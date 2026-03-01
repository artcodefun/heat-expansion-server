package bootstrap

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
	infraevents "github.com/artcodefun/heat-expansion-server/internal/auth/infrastructure/events"
)

// WireIntegrationEvents subscribes integration producer handlers to domain events.
func WireIntegrationEvents(s *AppServices, pub ports.EventPublisher) {
	p, ok := pub.(*infraevents.SimplePublisher)
	if !ok {
		return
	}

	p.Listen(func(ctx context.Context, e domain.DomainEvent) error {
		switch ev := e.(type) {
		case domain.AccountRegisteredEvent:
			return s.IntegrationProducer.HandleAccountRegistered(ctx, ev)
		}
		return nil
	})
}
