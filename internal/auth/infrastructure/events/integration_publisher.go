package events

import (
	"fmt"

	"github.com/artcodefun/heat-expansion-server/contracts/auth"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
)

type ConsoleIntegrationEventPublisher struct{}

func NewConsoleIntegrationEventPublisher() *ConsoleIntegrationEventPublisher {
	return &ConsoleIntegrationEventPublisher{}
}

func (p *ConsoleIntegrationEventPublisher) Publish(event auth.IntegrationEvent) error {
	fmt.Printf("[Integration Event] Publishing Type: %s, ID: %s, OriginID: %s, Payload: %+v\n",
		event.Type, event.ID, event.OriginID, event.Payload)
	return nil
}

var _ ports.IntegrationEventPublisher = (*ConsoleIntegrationEventPublisher)(nil)
