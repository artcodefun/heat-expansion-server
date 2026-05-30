package ports

import (
	"context"

	billingevents "github.com/artcodefun/heat-expansion-server/contracts/billing/events"
)

type IntegrationEventPublisher interface {
	Publish(ctx context.Context, event billingevents.IntegrationEvent) error
}
