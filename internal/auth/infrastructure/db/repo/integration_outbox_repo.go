package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/contracts/auth"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type IntegrationOutboxRepo struct {
	db *gen.Queries
}

func NewIntegrationOutboxRepo(db *sql.DB) *IntegrationOutboxRepo {
	return &IntegrationOutboxRepo{
		db: gen.New(db),
	}
}

func (r *IntegrationOutboxRepo) Tx(tx ports.Transaction) ports.IntegrationOutboxRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &IntegrationOutboxRepo{
			db: r.db.WithTx(sqlTx),
		}
	}
	return r
}

func (r *IntegrationOutboxRepo) Save(ctx context.Context, event auth.IntegrationEvent) error {
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
	return err
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

func (r *IntegrationOutboxRepo) ClaimUnpublished(limit int) ([]auth.IntegrationEvent, error) {
	rows, err := r.db.ClaimUnpublishedIntegrationEvents(context.Background(), int32(limit))
	if err != nil {
		return nil, err
	}

	events := make([]auth.IntegrationEvent, 0, len(rows))
	for _, row := range rows {
		evt, err := mappers.DecodeIntegrationEvent(row.Kind, row.Payload)
		if err != nil {
			return nil, err
		}

		events = append(events, evt)
	}

	return events, nil
}

func (r *IntegrationOutboxRepo) MarkPublished(id uuid.UUID, publishedAt int64) error {
	return r.db.MarkIntegrationEventPublished(context.Background(), gen.MarkIntegrationEventPublishedParams{
		ID: id,
		PublishedAt: sql.NullInt64{
			Int64: publishedAt,
			Valid: true,
		},
	})
}
