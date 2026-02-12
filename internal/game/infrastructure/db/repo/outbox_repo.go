package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
)

// OutboxEventRepo is a DB-backed implementation of ports.OutboxEventRepository
// using the domain_events outbox table via sqlc-generated queries.
type OutboxEventRepo struct {
	q *gen.Queries
}

func NewOutboxEventRepo(q *gen.Queries) *OutboxEventRepo { return &OutboxEventRepo{q: q} }

func (r *OutboxEventRepo) Tx(tx ports.Transaction) ports.OutboxEventRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &OutboxEventRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *OutboxEventRepo) Save(events []domain.DomainEvent) error {
	if len(events) == 0 {
		return nil
	}

	now := time.Now().Unix()
	for _, ev := range events {
		if ev == nil {
			continue
		}

		kind, payload, err := mappers.EncodeEvent(ev)
		if err != nil {
			return err
		}

		if _, err := r.q.InsertOutboxEvent(context.Background(), gen.InsertOutboxEventParams{
			Kind:      kind,
			Payload:   payload,
			CreatedAt: now,
		}); err != nil {
			return err
		}
	}

	return nil
}

func (r *OutboxEventRepo) ClaimUnpublished(limit int) ([]*ports.OutboxEventRecord, error) {
	if limit <= 0 {
		limit = 100
	}

	rows, err := r.q.ClaimUnpublishedOutboxEvents(context.Background(), int32(limit))
	if err != nil {
		return nil, err
	}

	out := make([]*ports.OutboxEventRecord, 0, len(rows))
	for _, row := range rows {
		decoded, err := mappers.DecodeEvent(row.Kind, row.Payload)
		if err != nil {
			// Skip malformed rows but continue processing others.
			continue
		}

		publishedAt := int64(0)
		if row.PublishedAt.Valid {
			publishedAt = row.PublishedAt.Int64
		}

		out = append(out, &ports.OutboxEventRecord{
			ID:          row.ID,
			Event:       decoded,
			CreatedAt:   row.CreatedAt,
			Published:   row.Published,
			PublishedAt: publishedAt,
		})
	}

	return out, nil
}

func (r *OutboxEventRepo) MarkPublished(id int64, publishedAt int64) error {
	if publishedAt == 0 {
		publishedAt = time.Now().Unix()
	}

	return r.q.MarkOutboxEventPublished(context.Background(), gen.MarkOutboxEventPublishedParams{
		PublishedAt: sql.NullInt64{Int64: publishedAt, Valid: true},
		ID:          id,
	})
}
