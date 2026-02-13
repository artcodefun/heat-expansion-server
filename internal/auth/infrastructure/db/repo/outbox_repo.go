package repo

import (
	"context"
	"database/sql"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type OutboxEventRepo struct {
	db *gen.Queries
}

func NewOutboxEventRepo(db *sql.DB) *OutboxEventRepo {
	return &OutboxEventRepo{
		db: gen.New(db),
	}
}

func (r *OutboxEventRepo) Tx(tx ports.Transaction) ports.OutboxEventRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &OutboxEventRepo{
			db: r.db.WithTx(sqlTx),
		}
	}
	return r
}

func (r *OutboxEventRepo) Save(ctx context.Context, events []domain.DomainEvent) error {
	for _, evt := range events {
		kind, payload, err := mappers.EncodeDomainEvent(evt)
		if err != nil {
			return err
		}

		err = r.db.SaveOutboxEvent(ctx, gen.SaveOutboxEventParams{
			ID:        evt.ID(),
			Kind:      kind,
			Payload:   payload,
			CreatedAt: evt.OccurredAt(),
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *OutboxEventRepo) ClaimUnpublished(limit int) ([]domain.DomainEvent, error) {
	rows, err := r.db.ClaimUnpublishedEvents(context.Background(), int32(limit))
	if err != nil {
		return nil, err
	}

	events := make([]domain.DomainEvent, 0, len(rows))
	for _, row := range rows {
		evt, err := mappers.DecodeDomainEvent(row.Kind, row.Payload)
		if err != nil {
			return nil, err
		}

		events = append(events, evt)
	}

	return events, nil
}

func (r *OutboxEventRepo) MarkPublished(id uuid.UUID, publishedAt int64) error {
	return r.db.MarkEventPublished(context.Background(), gen.MarkEventPublishedParams{
		ID: id,
		PublishedAt: sql.NullInt64{
			Int64: publishedAt,
			Valid: true,
		},
	})
}
