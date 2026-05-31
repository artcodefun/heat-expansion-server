package bootstrap

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"go.opentelemetry.io/otel/codes"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/events"
)

type Workers struct {
	DomainOutboxLoop      func(ctx context.Context)
	IntegrationOutboxLoop func(ctx context.Context)
	IntegrationEvtLoop    func(ctx context.Context)
	wg                    sync.WaitGroup
}

func NewWorkers(
	dbURL string,
	outbox *services.OutboxService,
	intOutbox *services.IntegrationOutboxService,
	consumer *events.RabbitMQConsumer,
) *Workers {
	return &Workers{
		DomainOutboxLoop: func(ctx context.Context) {
			slog.InfoContext(ctx, "billing outbox worker started")
			defer slog.InfoContext(ctx, "billing outbox worker stopped")

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "billing_domain_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					processBatch(ctx, "billing.outbox.process_batch", "billing outbox dispatch failed", func(batchCtx context.Context) error {
						return outbox.ProcessBatch(batchCtx, 100)
					})
				case <-signalChan:
					processBatch(ctx, "billing.outbox.process_batch", "billing outbox dispatch failed", func(batchCtx context.Context) error {
						return outbox.ProcessBatch(batchCtx, 100)
					})
				}
			}
		},
		IntegrationOutboxLoop: func(ctx context.Context) {
			slog.InfoContext(ctx, "billing integration outbox worker started")
			defer slog.InfoContext(ctx, "billing integration outbox worker stopped")

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "billing_integration_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					processBatch(ctx, "billing.integration_outbox.process_batch", "billing integration outbox dispatch failed", func(batchCtx context.Context) error {
						return intOutbox.ProcessBatch(batchCtx, 100)
					})
				case <-signalChan:
					processBatch(ctx, "billing.integration_outbox.process_batch", "billing integration outbox dispatch failed", func(batchCtx context.Context) error {
						return intOutbox.ProcessBatch(batchCtx, 100)
					})
				}
			}
		},
		IntegrationEvtLoop: func(ctx context.Context) {
			slog.InfoContext(ctx, "billing integration consumer started")
			defer slog.InfoContext(ctx, "billing integration consumer stopped")

			if err := consumer.Start(ctx); err != nil {
				slog.WarnContext(ctx, "failed to start billing integration consumer", "error", err.Error())
				return
			}

			<-ctx.Done()
		},
	}
}

// Start launches the billing background worker loops.
func (w *Workers) Start(ctx context.Context) {
	w.wg.Add(3)
	go func() { defer w.wg.Done(); w.DomainOutboxLoop(ctx) }()
	go func() { defer w.wg.Done(); w.IntegrationOutboxLoop(ctx) }()
	go func() { defer w.wg.Done(); w.IntegrationEvtLoop(ctx) }()
}

// Wait blocks until all billing background worker loops have exited.
func (w *Workers) Wait() {
	w.wg.Wait()
}

func processBatch(ctx context.Context, spanName, errMsg string, fn func(context.Context) error) {
	batchCtx, span := billingTracer.Start(ctx, spanName)
	defer span.End()

	if err := fn(batchCtx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		slog.ErrorContext(batchCtx, errMsg, "error", err.Error())
	}
}
