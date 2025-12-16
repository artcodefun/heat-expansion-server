package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/mappers"
)

// ScheduledJobRecord is an infrastructure-level representation of a scheduled job
// row used by the durable scheduler implementation. It carries a typed
// SchadulableJob; infrastructure is responsible for encoding/decoding via
// mappers when talking to the DB.
type ScheduledJobRecord struct {
	ID           int64
	Job          ports.SchadulableJob
	ExecuteAt    int64
	CreatedAt    int64
	Dispatched   bool
	DispatchedAt int64
}

// ScheduledJobRepo wraps sqlc-generated queries for scheduled_jobs.
type ScheduledJobRepo struct {
	q *gen.Queries
}

func NewScheduledJobRepo(q *gen.Queries) *ScheduledJobRepo { return &ScheduledJobRepo{q: q} }

// Tx returns a copy of the repository bound to the given transaction.
func (r *ScheduledJobRepo) Tx(tx ports.Transaction) *ScheduledJobRepo {
	if tx == nil {
		return r
	}
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &ScheduledJobRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

// Insert stores a new scheduled job row for the given typed job.
func (r *ScheduledJobRepo) Insert(ctx context.Context, job ports.SchadulableJob, executeAt, createdAt int64) (int64, error) {
	if createdAt == 0 {
		createdAt = time.Now().Unix()
	}

	kind, payload, err := mappers.EncodeJob(job)
	if err != nil {
		return 0, err
	}

	return r.q.InsertScheduledJob(ctx, gen.InsertScheduledJobParams{
		Kind:      kind,
		Payload:   payload,
		ExecuteAt: executeAt,
		CreatedAt: createdAt,
	})
}

// ClaimDue returns up to limit jobs that are due for execution (execute_at <= now)
// using FOR UPDATE SKIP LOCKED semantics. It is expected to be called within a
// transaction to guarantee exclusive claims in multi-node setups.
func (r *ScheduledJobRepo) ClaimDue(ctx context.Context, now int64, limit int32) ([]ScheduledJobRecord, error) {
	if limit <= 0 {
		limit = 100
	}

	rows, err := r.q.ClaimDueScheduledJobs(ctx, gen.ClaimDueScheduledJobsParams{
		ExecuteAt: now,
		Limit:     limit,
	})
	if err != nil {
		return nil, err
	}

	out := make([]ScheduledJobRecord, 0, len(rows))
	for _, row := range rows {
		job, err := mappers.DecodeJob(row.Kind, row.Payload)
		if err != nil {
			// Skip malformed rows but continue processing others.
			continue
		}

		rec := ScheduledJobRecord{
			ID:         row.ID,
			Job:        job,
			ExecuteAt:  row.ExecuteAt,
			CreatedAt:  row.CreatedAt,
			Dispatched: row.Dispatched,
		}
		if row.DispatchedAt.Valid {
			rec.DispatchedAt = row.DispatchedAt.Int64
		}
		out = append(out, rec)
	}

	return out, nil
}

// MarkDispatched marks a job as dispatched at the given timestamp.
func (r *ScheduledJobRepo) MarkDispatched(ctx context.Context, id int64, dispatchedAt int64) error {
	if dispatchedAt == 0 {
		dispatchedAt = time.Now().Unix()
	}

	return r.q.MarkScheduledJobDispatched(ctx, gen.MarkScheduledJobDispatchedParams{
		DispatchedAt: sql.NullInt64{Int64: dispatchedAt, Valid: true},
		ID:           id,
	})
}

// GetNext returns the next (earliest execute_at) undispatched job, or nil if none.
func (r *ScheduledJobRepo) GetNext(ctx context.Context) (*ScheduledJobRecord, error) {
	row, err := r.q.GetNextScheduledJob(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	job, err := mappers.DecodeJob(row.Kind, row.Payload)
	if err != nil {
		return nil, err
	}

	rec := &ScheduledJobRecord{
		ID:         row.ID,
		Job:        job,
		ExecuteAt:  row.ExecuteAt,
		CreatedAt:  row.CreatedAt,
		Dispatched: row.Dispatched,
	}
	if row.DispatchedAt.Valid {
		rec.DispatchedAt = row.DispatchedAt.Int64
	}
	return rec, nil
}
