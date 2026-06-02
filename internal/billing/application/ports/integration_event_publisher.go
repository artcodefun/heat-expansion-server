package ports

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/contracts/events"
)

type IntegrationEventPublisher interface {
	Publish(ctx context.Context, event events.IntegrationEvent) error
}
