package bootstrap

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	infraevents "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/events"
)

var billingTracer = otel.Tracer("heat-expansion-billing")

// WireIntegrationEvents subscribes integration producer handlers to domain events.
func WireIntegrationEvents(s *AppServices, pub ports.EventPublisher) {
	p, ok := pub.(*infraevents.SimplePublisher)
	if !ok {
		return
	}

	p.Listen(func(ctx context.Context, e domain.DomainEvent) error {
		ctx, span := billingTracer.Start(ctx, fmt.Sprintf("%T", e))
		defer span.End()

		err := func() error {
			switch ev := e.(type) {
			case domain.OrderPaidEvent:
				return s.IntegrationProducer.HandleOrderPaid(ctx, ev)
			}
			return nil
		}()

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		return err
	})
}
