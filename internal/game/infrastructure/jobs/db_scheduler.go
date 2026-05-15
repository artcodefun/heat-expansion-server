package jobs

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/repo"
)

// DBScheduler is a durable scheduler implementation backed by the scheduled_jobs table.
// It implements ports.Scheduler and forwards each claimed job to a single listener
// registered at startup. Not safe for multi-process usage by itself; safety comes
// from the database-level locking used in the repository.
type DBScheduler struct {
	txMgr ports.TransactionManager
	repo  *repo.ScheduledJobRepo

	listener func(context.Context, ports.SchadulableJob) error
}

var _ ports.Scheduler = (*DBScheduler)(nil)

// NewDBScheduler constructs a new DB-backed scheduler.
func NewDBScheduler(txMgr ports.TransactionManager, r *repo.ScheduledJobRepo) *DBScheduler {
	return &DBScheduler{
		txMgr: txMgr,
		repo:  r,
	}
}

// Schedule persists a job into the scheduled_jobs table.
func (s *DBScheduler) Schedule(ctx context.Context, job ports.SchadulableJob, executeAt int64) error {
	if job == nil {
		return fmt.Errorf("nil job")
	}
	now := time.Now().Unix()
	if executeAt <= 0 {
		executeAt = now
	}
	_, err := s.repo.Insert(ctx, job, executeAt, now)
	return err
}

// Listen registers a single callback to receive job payloads as they are dispatched.
// Returns an unsubscribe function.
func (s *DBScheduler) Listen(cb func(context.Context, ports.SchadulableJob) error) (unsubscribe func()) {
	s.listener = cb
	return func() {
		s.listener = nil
	}
}

// Run starts the dispatch loop, claiming due jobs and invoking subscribers until ctx is done.
func (s *DBScheduler) Run(ctx context.Context) {
	const batchLimit int32 = 100
	const idleSleep = 5 * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		now := time.Now().Unix()
		pollCtx, span := otel.Tracer("heat-expansion-game").Start(ctx, "game.scheduler.process_batch")
		if err := s.processDueJobs(pollCtx, now, batchLimit); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			slog.ErrorContext(pollCtx, "game scheduler batch processing failed", "error", err.Error())
		}
		next, err := s.repo.GetNext(pollCtx)
		span.End()
		sleepFor := idleSleep

		if err == nil && next != nil {
			now = time.Now().Unix()
			wait := time.Duration(next.ExecuteAt-now) * time.Second
			sleepFor = min(idleSleep, wait)
		}

		select {
		case <-time.After(sleepFor):
		case <-ctx.Done():
			return
		}
		continue

	}
}

func (s *DBScheduler) processDueJobs(ctx context.Context, now int64, limit int32) error {
	return s.txMgr.WithTx(ctx, func(tx ports.Transaction) error {
		repoTx := s.repo.Tx(tx)

		records, err := repoTx.ClaimDue(ctx, now, limit)
		if err != nil {
			return err
		}
		if len(records) == 0 {
			return nil
		}

		nowTs := time.Now().Unix()
		for _, row := range records {
			// Deliver the job to listener before marking it as dispatched.
			// If a handler fails, skip marking dispatched so it can be retried.
			if err := s.notifyListener(ctx, row.Job); err != nil {
				slog.WarnContext(ctx, "game scheduled job handler failed; leaving job undispatched for retry",
					"scheduled_job_id", row.ID,
					"job_type", fmt.Sprintf("%T", row.Job),
					"execute_at", row.ExecuteAt,
					"error", err.Error(),
				)
				continue
			}

			if err := repoTx.MarkDispatched(ctx, row.ID, nowTs); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *DBScheduler) notifyListener(ctx context.Context, job ports.SchadulableJob) error {
	listener := s.listener
	if listener == nil {
		return nil
	}
	return listener(ctx, job)
}
