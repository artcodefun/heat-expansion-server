package bootstrap

import (
	"context"
	"log/slog"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/repo"
)

type Workers struct {
	DomainOutboxLoop      func(ctx context.Context)
	IntegrationOutboxLoop func(ctx context.Context)
}

func NewWorkers(
	dbURL string,
	outbox *services.OutboxService,
	intOutbox *services.IntegrationOutboxService,
) *Workers {
	return &Workers{
		DomainOutboxLoop: func(ctx context.Context) {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "auth_domain_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := outbox.ProcessBatch(ctx, 100); err != nil {
						slog.Error("auth outbox dispatch failed", "error", err.Error())
					}
				case <-signalChan:
					if err := outbox.ProcessBatch(ctx, 100); err != nil {
						slog.Error("auth outbox dispatch failed", "error", err.Error())
					}
				}
			}
		},
		IntegrationOutboxLoop: func(ctx context.Context) {
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "auth_integration_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					if err := intOutbox.ProcessBatch(ctx, 100); err != nil {
						slog.Error("auth integration outbox dispatch failed", "error", err.Error())
					}
				case <-signalChan:
					if err := intOutbox.ProcessBatch(ctx, 100); err != nil {
						slog.Error("auth integration outbox dispatch failed", "error", err.Error())
					}
				}
			}
		},
	}
}
