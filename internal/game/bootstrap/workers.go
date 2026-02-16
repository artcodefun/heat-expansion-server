package bootstrap

import (
	"context"
	"log/slog"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/repo"
)

type Workers struct {
	OutboxLoop func(ctx context.Context)
}

func NewWorkers(
	dbURL string,
	outbox *services.OutboxService,
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
					if err := outbox.ProcessBatch(100); err != nil {
						slog.Error("game outbox dispatch failed", "error", err.Error())
					}
				case <-signalChan:
					if err := outbox.ProcessBatch(100); err != nil {
						slog.Error("game outbox dispatch failed", "error", err.Error())
					}
				}
			}
		},
	}
}
