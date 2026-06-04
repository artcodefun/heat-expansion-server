package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"

	authv1 "github.com/artcodefun/heat-expansion-server/contracts/auth/events/v1"
	"github.com/artcodefun/heat-expansion-server/contracts/events"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	platformevents "github.com/artcodefun/heat-expansion-server/internal/platform/events"
	"github.com/artcodefun/heat-expansion-server/internal/platform/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

var billingTracer = otel.Tracer("heat-expansion-billing")

// WireConsumerIntegrationEvents binds inbound auth integration events to command
// handlers. Billing maintains a local users projection from these events so it
// can attach a customer email to payment receipts.
func WireConsumerIntegrationEvents(c *Commands, consumer *rabbitmq.RabbitMQConsumer, authExchange, authQueue string) {
	consumer.Subscribe(authExchange, authQueue, "auth.#", func(ctx context.Context, d amqp.Delivery) error {
		ctx, span := billingTracer.Start(ctx, "billing.integration."+d.RoutingKey)
		defer span.End()

		err := func() error {
			var envelope events.IntegrationEvent
			if err := json.Unmarshal(d.Body, &envelope); err != nil {
				return err
			}

			switch envelope.Type {
			case authv1.EventAccountRegisteredV1:
				var p authv1.AccountRegisteredV1
				if err := json.Unmarshal(envelope.Payload, &p); err != nil {
					return err
				}
				return c.User.HandleAccountRegisteredV1Event(ctx, p)
			default:
				slog.WarnContext(ctx, "received unknown auth integration event type", "type", envelope.Type)
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

// WireIntegrationEvents subscribes integration producer handlers to domain events.
func WireIntegrationEvents(s *AppServices, pub ports.EventPublisher) {
	p, ok := pub.(*platformevents.SimplePublisher[domain.DomainEvent])
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
