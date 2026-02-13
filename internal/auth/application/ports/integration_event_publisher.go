package ports

import "github.com/artcodefun/heat-expansion-server/contracts/auth"

type IntegrationEventPublisher interface {
	Publish(event auth.IntegrationEvent) error
}
