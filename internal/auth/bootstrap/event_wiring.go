package bootstrap

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/codes"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
	platformevents "github.com/artcodefun/heat-expansion-server/internal/platform/events"
)

// WireIntegrationEvents subscribes integration producer handlers to domain events.
func WireIntegrationEvents(s *AppServices, pub ports.EventPublisher) {
	p, ok := pub.(*platformevents.SimplePublisher[domain.DomainEvent])
	if !ok {
		return
	}

	p.Listen(func(ctx context.Context, e domain.DomainEvent) error {
		ctx, span := authTracer.Start(ctx, fmt.Sprintf("%T", e))
		defer span.End()

		err := func() error {
			switch ev := e.(type) {
			case domain.AccountRegisteredEvent:
				return s.IntegrationProducer.HandleAccountRegistered(ctx, ev)
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
