package ports

import (
	"context"

	authevents "github.com/artcodefun/heat-expansion-server/contracts/auth/events"
)

type IntegrationEventPublisher interface {
	Publish(ctx context.Context, event authevents.IntegrationEvent) error
}
