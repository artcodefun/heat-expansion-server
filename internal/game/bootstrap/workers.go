package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/jobs"
	"github.com/artcodefun/heat-expansion-server/internal/platform/rabbitmq"
)

type Workers struct {
	OutboxLoop         func(ctx context.Context) error
	SchedulerLoop      func(ctx context.Context) error
	IntegrationEvtLoop func(ctx context.Context) error
	g                  *errgroup.Group
}

func NewWorkers(
	dbURL string,
	outbox *services.OutboxService,
	scheduler ports.Scheduler,
	consumer *rabbitmq.RabbitMQConsumer,
) *Workers {
	runner, ok := scheduler.(*jobs.DBScheduler)
	if !ok {
		panic(fmt.Sprintf("game workers require *jobs.DBScheduler, got %T", scheduler))
	}

	return &Workers{
		OutboxLoop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "game outbox worker started")
			defer slog.InfoContext(ctx, "game outbox worker stopped")

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "game_domain_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return nil
				case <-ticker.C:
					processBatch(ctx, "game.outbox.process_batch", "game outbox dispatch failed", func(batchCtx context.Context) error {
						return outbox.ProcessBatch(batchCtx, 100)
					})
				case <-signalChan:
					processBatch(ctx, "game.outbox.process_batch", "game outbox dispatch failed", func(batchCtx context.Context) error {
						return outbox.ProcessBatch(batchCtx, 100)
					})
				}
			}
		},
		SchedulerLoop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "game scheduler worker started")
			defer slog.InfoContext(ctx, "game scheduler worker stopped")
			runner.Run(ctx)
			return nil
		},
		IntegrationEvtLoop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "game integration consumer started")
			defer slog.InfoContext(ctx, "game integration consumer stopped")

			if err := consumer.Start(ctx); err != nil {
				return fmt.Errorf("game integration consumer: %w", err)
			}
			return nil
		},
	}
}

// Start launches all background worker loops. A loop failure cancels the
// group context so the sibling loops stop; Wait reports the first error.
func (w *Workers) Start(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	w.g = g
	g.Go(func() error { return w.SchedulerLoop(ctx) })
	g.Go(func() error { return w.IntegrationEvtLoop(ctx) })
	g.Go(func() error { return w.OutboxLoop(ctx) })
}

// Wait blocks until all background worker loops have exited and returns the
// first loop failure, if any.
func (w *Workers) Wait() error {
	return w.g.Wait()
}

func processBatch(ctx context.Context, spanName, errMsg string, fn func(context.Context) error) {
	batchCtx, span := gameTracer.Start(ctx, spanName)
	defer span.End()

	if err := fn(batchCtx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(batchCtx, errMsg, "error", err.Error())
	}
}
