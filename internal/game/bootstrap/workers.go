package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"go.opentelemetry.io/otel/codes"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/events"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/jobs"
)

type Workers struct {
	OutboxLoop         func(ctx context.Context)
	SchedulerLoop      func(ctx context.Context)
	IntegrationEvtLoop func(ctx context.Context)
	wg                 sync.WaitGroup
}

func NewWorkers(
	dbURL string,
	outbox *services.OutboxService,
	scheduler ports.Scheduler,
	consumer *events.RabbitMQConsumer,
) *Workers {
	runner, ok := scheduler.(*jobs.DBScheduler)
	if !ok {
		panic(fmt.Sprintf("game workers require *jobs.DBScheduler, got %T", scheduler))
	}

	return &Workers{
		OutboxLoop: func(ctx context.Context) {
			slog.InfoContext(ctx, "game outbox worker started")
			defer slog.InfoContext(ctx, "game outbox worker stopped")

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "game_domain_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return
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
		SchedulerLoop: func(ctx context.Context) {
			slog.InfoContext(ctx, "game scheduler worker started")
			defer slog.InfoContext(ctx, "game scheduler worker stopped")
			runner.Run(ctx)
		},
		IntegrationEvtLoop: func(ctx context.Context) {
			slog.InfoContext(ctx, "game integration consumer started")
			defer slog.InfoContext(ctx, "game integration consumer stopped")

			if err := consumer.Start(ctx); err != nil {
				slog.WarnContext(ctx, "failed to start game integration consumer", "error", err.Error())
				return
			}

			<-ctx.Done()
		},
	}
}

// Start launches all background worker loops.
func (w *Workers) Start(ctx context.Context) {
	w.wg.Add(3)
	go func() { defer w.wg.Done(); w.SchedulerLoop(ctx) }()
	go func() { defer w.wg.Done(); w.IntegrationEvtLoop(ctx) }()
	go func() { defer w.wg.Done(); w.OutboxLoop(ctx) }()
}

// Wait blocks until all background worker loops have exited.
func (w *Workers) Wait() {
	w.wg.Wait()
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
