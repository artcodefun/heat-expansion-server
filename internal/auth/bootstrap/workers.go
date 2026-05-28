package bootstrap

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/repo"
)

var authTracer = otel.Tracer("heat-expansion-auth")

type Workers struct {
	DomainOutboxLoop      func(ctx context.Context)
	IntegrationOutboxLoop func(ctx context.Context)
	wg                    sync.WaitGroup
}

func NewWorkers(
	dbURL string,
	outbox *services.OutboxService,
	intOutbox *services.IntegrationOutboxService,
) *Workers {
	return &Workers{
		DomainOutboxLoop: func(ctx context.Context) {
			slog.InfoContext(ctx, "auth outbox worker started")
			defer slog.InfoContext(ctx, "auth outbox worker stopped")

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "auth_domain_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					processBatch(ctx, "auth.outbox.process_batch", "auth outbox dispatch failed", func(batchCtx context.Context) error {
						return outbox.ProcessBatch(batchCtx, 100)
					})
				case <-signalChan:
					processBatch(ctx, "auth.outbox.process_batch", "auth outbox dispatch failed", func(batchCtx context.Context) error {
						return outbox.ProcessBatch(batchCtx, 100)
					})
				}
			}
		},
		IntegrationOutboxLoop: func(ctx context.Context) {
			slog.InfoContext(ctx, "auth integration outbox worker started")
			defer slog.InfoContext(ctx, "auth integration outbox worker stopped")

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "auth_integration_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					processBatch(ctx, "auth.integration_outbox.process_batch", "auth integration outbox dispatch failed", func(batchCtx context.Context) error {
						return intOutbox.ProcessBatch(batchCtx, 100)
					})
				case <-signalChan:
					processBatch(ctx, "auth.integration_outbox.process_batch", "auth integration outbox dispatch failed", func(batchCtx context.Context) error {
						return intOutbox.ProcessBatch(batchCtx, 100)
					})
				}
			}
		},
	}
}

// Start launches the auth background worker loops.
func (w *Workers) Start(ctx context.Context) {
	w.wg.Add(2)
	go func() { defer w.wg.Done(); w.DomainOutboxLoop(ctx) }()
	go func() { defer w.wg.Done(); w.IntegrationOutboxLoop(ctx) }()
}

// Wait blocks until all auth background worker loops have exited.
func (w *Workers) Wait() {
	w.wg.Wait()
}

func processBatch(ctx context.Context, spanName, errMsg string, fn func(context.Context) error) {
	batchCtx, span := authTracer.Start(ctx, spanName)
	defer span.End()

	if err := fn(batchCtx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(batchCtx, errMsg, "error", err.Error())
	}
}
