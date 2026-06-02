package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"golang.org/x/sync/errgroup"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/db/repo"
	"github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/events"
)

var authTracer = otel.Tracer("heat-expansion-auth")

type Workers struct {
	DomainOutboxLoop      func(ctx context.Context) error
	IntegrationOutboxLoop func(ctx context.Context) error
	PublisherLoop         func(ctx context.Context) error
	g                     *errgroup.Group
}

func NewWorkers(
	dbURL string,
	outbox *services.OutboxService,
	intOutbox *services.IntegrationOutboxService,
	publisher *events.RabbitMQPublisher,
) *Workers {
	return &Workers{
		DomainOutboxLoop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "auth outbox worker started")
			defer slog.InfoContext(ctx, "auth outbox worker stopped")

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "auth_domain_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return nil
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
		IntegrationOutboxLoop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "auth integration outbox worker started")
			defer slog.InfoContext(ctx, "auth integration outbox worker stopped")

			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()

			listener := repo.NewPostgresListener(dbURL, "auth_integration_events")
			signalChan := listener.Events

			for {
				select {
				case <-ctx.Done():
					return nil
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
		PublisherLoop: func(ctx context.Context) error {
			slog.InfoContext(ctx, "auth rabbitmq publisher started")
			defer slog.InfoContext(ctx, "auth rabbitmq publisher stopped")

			if err := publisher.Start(ctx); err != nil {
				return fmt.Errorf("auth rabbitmq publisher: %w", err)
			}
			return nil
		},
	}
}

// Start launches the auth background worker loops. A loop failure cancels the
// group context so the sibling loops stop; Wait reports the first error.
func (w *Workers) Start(ctx context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	w.g = g
	g.Go(func() error { return w.DomainOutboxLoop(ctx) })
	g.Go(func() error { return w.IntegrationOutboxLoop(ctx) })
	g.Go(func() error { return w.PublisherLoop(ctx) })
}

// Wait blocks until all auth background worker loops have exited and returns
// the first loop failure, if any.
func (w *Workers) Wait() error {
	return w.g.Wait()
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
