package bootstrap

import (
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/services"
)

type AppServices struct {
	Outbox              *services.OutboxService
	IntegrationOutbox   *services.IntegrationOutboxService
	IntegrationProducer *services.IntegrationProducerService
}

func NewAppServices(adapters *Adapters) *AppServices {
	return &AppServices{
		Outbox:              services.NewOutboxService(adapters.Outbox, adapters.Events, adapters.TxMgr),
		IntegrationOutbox:   services.NewIntegrationOutboxService(adapters.IntegrationOutbox, adapters.IntegrationEvents, adapters.TxMgr),
		IntegrationProducer: services.NewIntegrationProducerService(adapters.IntegrationOutbox, adapters.TxMgr),
	}
}
