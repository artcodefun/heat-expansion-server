package ports

import "github.com/artcodefun/heat-expansion-api/internal/game/core/domain"

// EventPublisher defines the interface for publishing domain events.
type EventPublisher interface {
	Publish(event domain.DomainEvent) error
}
