package ports

import (
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
)

type EventPublisher interface {
	Publish(event domain.DomainEvent) error
}
