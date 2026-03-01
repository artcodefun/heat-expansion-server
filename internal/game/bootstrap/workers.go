package bootstrap

import (
	"context"
	"log/slog"
	"time"

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
					if err := outbox.ProcessBatch(ctx, 100); err != nil {
						slog.Error("game outbox dispatch failed", "error", err.Error())
					}
				case <-signalChan:
					if err := outbox.ProcessBatch(ctx, 100); err != nil {
						slog.Error("game outbox dispatch failed", "error", err.Error())
					}
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
				slog.Warn("failed to start game integration consumer", "error", err.Error())
			}
		},
	}
}
