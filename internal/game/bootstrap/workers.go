package bootstrap

import (
	"context"
	"log/slog"
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
}

func NewWorkers(
	dbURL string,
	outbox *services.OutboxService,
	scheduler ports.Scheduler,
	consumer *events.RabbitMQConsumer,
) *Workers {
	return &Workers{
		OutboxLoop: func(ctx context.Context) {
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
			runner, ok := scheduler.(*jobs.DBScheduler)
			if !ok {
				return
			}
			runner.Run(ctx)
		},
		IntegrationEvtLoop: func(ctx context.Context) {
			if err := consumer.Start(ctx); err != nil {
				slog.WarnContext(ctx, "failed to start game integration consumer", "error", err.Error())
			}
		},
	}
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
