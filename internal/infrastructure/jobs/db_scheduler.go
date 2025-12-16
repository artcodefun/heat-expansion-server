package jobs

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/repo"
)

// DBScheduler is a durable scheduler implementation backed by the scheduled_jobs table.
// It implements ports.Scheduler and exposes a Subscribe method compatible with
// the existing wiring in bootstrap/command_wiring.go.
type DBScheduler struct {
	txMgr ports.TransactionManager
	repo  *repo.ScheduledJobRepo

	mu          sync.Mutex
	subscribers map[int]func(ports.SchadulableJob)
	nextSubID   int
}

var _ ports.Scheduler = (*DBScheduler)(nil)

// NewDBScheduler constructs a new DB-backed scheduler.
func NewDBScheduler(txMgr ports.TransactionManager, r *repo.ScheduledJobRepo) *DBScheduler {
	return &DBScheduler{
		txMgr:       txMgr,
		repo:        r,
		subscribers: make(map[int]func(ports.SchadulableJob)),
	}
}

// Schedule persists a job into the scheduled_jobs table.
func (s *DBScheduler) Schedule(job ports.SchadulableJob, executeAt int64) error {
	if job == nil {
		return fmt.Errorf("nil job")
	}
	now := time.Now().Unix()
	if executeAt <= 0 {
		executeAt = now
	}
	_, err := s.repo.Insert(context.Background(), job, executeAt, now)
	return err
}

// Subscribe registers a callback to receive job payloads as they are dispatched.
// Returns an unsubscribe function.
func (s *DBScheduler) Subscribe(cb func(ports.SchadulableJob)) (unsubscribe func()) {
	s.mu.Lock()
	id := s.nextSubID
	s.nextSubID++
	s.subscribers[id] = cb
	s.mu.Unlock()
	return func() {
		s.mu.Lock()
		delete(s.subscribers, id)
		s.mu.Unlock()
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
		_ = s.processDueJobs(ctx, now, batchLimit)

		next, err := s.repo.GetNext(ctx)
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
	return s.txMgr.WithTx(func(tx ports.Transaction) error {
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
			// Deliver the job to subscribers before marking it as dispatched,
			// mirroring OutboxService semantics.
			s.notifySubscribers(row.Job)

			if err := repoTx.MarkDispatched(ctx, row.ID, nowTs); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *DBScheduler) notifySubscribers(job ports.SchadulableJob) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, cb := range s.subscribers {
		cb(job)
	}
}
