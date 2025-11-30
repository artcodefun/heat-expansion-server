package commands

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
)

func publishEvents(events []domain.DomainEvent, publisher ports.EventPublisher) {
	for _, e := range events {
		publisher.Publish(e)
	}
}
