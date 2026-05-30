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
	q *gen.Queries
}

func NewOutboxEventRepo(q *gen.Queries) *OutboxEventRepo {
	return &OutboxEventRepo{q: q}
}

func (r *OutboxEventRepo) Tx(tx ports.Transaction) ports.OutboxEventRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &OutboxEventRepo{
			q: r.q.WithTx(sqlTx),
		}
	}
	return r
}

func (r *OutboxEventRepo) Save(ctx context.Context, events []domain.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	for _, evt := range events {
		if evt == nil {
			continue
		}

		kind, payload, err := mappers.EncodeDomainEvent(evt)
		if err != nil {
			return err
		}

		err = r.q.SaveOutboxEvent(ctx, gen.SaveOutboxEventParams{
			ID:        evt.ID(),
			Kind:      kind,
			Payload:   payload,
			CreatedAt: evt.OccurredAt(),
		})
		if err != nil {
			return err
		}
	}
	return r.q.NotifyOutboxEvent(ctx)
}

func (r *OutboxEventRepo) ClaimUnpublished(ctx context.Context, limit int) ([]domain.DomainEvent, error) {
	if limit <= 0 {
		limit = 100
	}

	rows, err := r.q.ClaimUnpublishedEvents(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	events := make([]domain.DomainEvent, 0, len(rows))
	for _, row := range rows {
		evt, err := mappers.DecodeDomainEvent(row.Kind, row.Payload)
		if err != nil {
			// Skip malformed rows but continue processing others so a single
			// bad row cannot stall the publisher.
			continue
		}

		events = append(events, evt)
	}

	return events, nil
}

func (r *OutboxEventRepo) MarkPublished(ctx context.Context, id uuid.UUID, publishedAt int64) error {
	return r.q.MarkEventPublished(ctx, gen.MarkEventPublishedParams{
		ID: id,
		PublishedAt: sql.NullInt64{
			Int64: publishedAt,
			Valid: true,
		},
	})
}
