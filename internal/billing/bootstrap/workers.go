package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/platform/rabbitmq"
)

type Workers struct {
	DomainOutboxLoop      func(ctx context.Context) error
	IntegrationOutboxLoop func(ctx context.Context) error
	IntegrationEvtLoop    func(ctx context.Context) error
	PublisherLoop         func(ctx context.Context) error
	g                     *errgroup.Group
}

func NewWorkers(
	dbURL string,
	outbox *services.OutboxService,
	intOutbox *services.IntegrationOutboxService,
	consumer *rabbitmq.RabbitMQConsumer,
	publisher *rabbitmq.RabbitMQPublisher,
) *Workers {
	return &Workers{
		DomainOutboxLoop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "billing outbox worker started")
			defer slog.InfoContext(ctx, "billing outbox worker stopped")

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "billing_domain_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return nil
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
		IntegrationOutboxLoop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "billing integration outbox worker started")
			defer slog.InfoContext(ctx, "billing integration outbox worker stopped")

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "billing_integration_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return nil
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
		IntegrationEvtLoop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "billing integration consumer started")
			defer slog.InfoContext(ctx, "billing integration consumer stopped")

			if err := consumer.Start(ctx); err != nil {
				return fmt.Errorf("billing integration consumer: %w", err)
			}
			return nil
		},
		PublisherLoop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "billing rabbitmq publisher started")
			defer slog.InfoContext(ctx, "billing rabbitmq publisher stopped")

			if err := publisher.Start(ctx); err != nil {
				return fmt.Errorf("billing rabbitmq publisher: %w", err)
			}
			return nil
		},
	}
}

// Start launches the billing background worker loops. A loop failure cancels
// the group context so the sibling loops stop; Wait reports the first error.
func (w *Workers) Start(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	w.g = g
	g.Go(func() error { return w.DomainOutboxLoop(ctx) })
	g.Go(func() error { return w.IntegrationOutboxLoop(ctx) })
	g.Go(func() error { return w.IntegrationEvtLoop(ctx) })
	g.Go(func() error { return w.PublisherLoop(ctx) })
}

// Wait blocks until all billing background worker loops have exited and
// returns the first loop failure, if any.
func (w *Workers) Wait() error {
	return w.g.Wait()
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
