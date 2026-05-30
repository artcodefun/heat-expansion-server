package ports

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
)

type EventPublisher interface {
	Publish(ctx context.Context, event domain.DomainEvent) error
}
