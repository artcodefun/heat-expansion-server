package ports

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/contracts/auth"
)

type IntegrationEventPublisher interface {
	Publish(ctx context.Context, event auth.IntegrationEvent) error
}
