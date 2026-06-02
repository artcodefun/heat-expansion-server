package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/contracts/events"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type IntegrationOutboxRepo struct {
	q *gen.Queries
}

func NewIntegrationOutboxRepo(q *gen.Queries) *IntegrationOutboxRepo {
	return &IntegrationOutboxRepo{q: q}
}

func (r *IntegrationOutboxRepo) Tx(tx ports.Transaction) ports.IntegrationOutboxRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &IntegrationOutboxRepo{
			q: r.q.WithTx(sqlTx),
		}
	}
	return r
}

func (r *IntegrationOutboxRepo) Save(ctx context.Context, event events.IntegrationEvent) error {
	kind, payload, err := mappers.EncodeIntegrationEvent(event)
	if err != nil {
		return err
	}

	err = r.q.SaveIntegrationEvent(ctx, gen.SaveIntegrationEventParams{
		ID:        event.ID,
		Kind:      kind,
		Payload:   payload,
		CreatedAt: event.OccurredAt,
		OriginID: uuid.NullUUID{
			UUID:  event.OriginID,
			Valid: event.OriginID != uuid.Nil,
		},
	})
	if err != nil {
		return err
	}
	return r.q.NotifyIntegrationOutboxEvent(ctx)
}

func (r *IntegrationOutboxRepo) Exists(ctx context.Context, originID uuid.UUID, eventType string) (bool, error) {
	return r.q.IntegrationEventExists(ctx, gen.IntegrationEventExistsParams{
		OriginID: uuid.NullUUID{
			UUID:  originID,
			Valid: originID != uuid.Nil,
		},
		Kind: eventType,
	})
}

func (r *IntegrationOutboxRepo) ClaimUnpublished(ctx context.Context, limit int) ([]events.IntegrationEvent, error) {
	rows, err := r.q.ClaimUnpublishedIntegrationEvents(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	evts := make([]events.IntegrationEvent, 0, len(rows))
	for _, row := range rows {
		evt, err := mappers.DecodeIntegrationEvent(row.Kind, row.Payload)
		if err != nil {
			// Skip malformed rows but continue processing others so a single
			// bad row cannot stall the publisher.
			continue
		}

		evts = append(evts, evt)
	}

	return evts, nil
}

func (r *IntegrationOutboxRepo) MarkPublished(ctx context.Context, id uuid.UUID, publishedAt int64) error {
	return r.q.MarkIntegrationEventPublished(ctx, gen.MarkIntegrationEventPublishedParams{
		ID: id,
		PublishedAt: sql.NullInt64{
			Int64: publishedAt,
			Valid: true,
		},
	})
}
