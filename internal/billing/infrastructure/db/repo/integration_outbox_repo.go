package repo

import (
	"context"
	"database/sql"

	billingevents "github.com/artcodefun/heat-expansion-server/contracts/billing/events"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type IntegrationOutboxRepo struct {
	db *gen.Queries
}

func NewIntegrationOutboxRepo(q *gen.Queries) *IntegrationOutboxRepo {
	return &IntegrationOutboxRepo{db: q}
}

func (r *IntegrationOutboxRepo) Tx(tx ports.Transaction) ports.IntegrationOutboxRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &IntegrationOutboxRepo{
			db: r.db.WithTx(sqlTx),
		}
	}
	return r
}

func (r *IntegrationOutboxRepo) Save(ctx context.Context, event billingevents.IntegrationEvent) error {
	kind, payload, err := mappers.EncodeIntegrationEvent(event)
	if err != nil {
		return err
	}

	err = r.db.SaveIntegrationEvent(ctx, gen.SaveIntegrationEventParams{
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
	return r.db.NotifyIntegrationOutboxEvent(ctx)
}

func (r *IntegrationOutboxRepo) Exists(ctx context.Context, originID uuid.UUID, eventType string) (bool, error) {
	return r.db.IntegrationEventExists(ctx, gen.IntegrationEventExistsParams{
		OriginID: uuid.NullUUID{
			UUID:  originID,
			Valid: originID != uuid.Nil,
		},
		Kind: eventType,
	})
}

func (r *IntegrationOutboxRepo) ClaimUnpublished(ctx context.Context, limit int) ([]billingevents.IntegrationEvent, error) {
	rows, err := r.db.ClaimUnpublishedIntegrationEvents(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	events := make([]billingevents.IntegrationEvent, 0, len(rows))
	for _, row := range rows {
		evt, err := mappers.DecodeIntegrationEvent(row.Kind, row.Payload)
		if err != nil {
			return nil, err
		}

		events = append(events, evt)
	}

	return events, nil
}

func (r *IntegrationOutboxRepo) MarkPublished(ctx context.Context, id uuid.UUID, publishedAt int64) error {
	return r.db.MarkIntegrationEventPublished(ctx, gen.MarkIntegrationEventPublishedParams{
		ID: id,
		PublishedAt: sql.NullInt64{
			Int64: publishedAt,
			Valid: true,
		},
	})
}
