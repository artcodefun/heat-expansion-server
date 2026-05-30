package bootstrap

import (
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/services"
)

// AppServices aggregates application-level services shared across commands and bootstrap wiring.
type AppServices struct {
	Outbox              *services.OutboxService
	IntegrationOutbox   *services.IntegrationOutboxService
	IntegrationProducer *services.IntegrationProducerService
}

// NewAppServices constructs all application-level services using the provided secondary adapters.
func NewAppServices(a *Adapters) *AppServices {
	return &AppServices{
		Outbox:              services.NewOutboxService(a.Outbox, a.Events, a.TxMgr),
		IntegrationOutbox:   services.NewIntegrationOutboxService(a.IntegrationOutbox, a.IntegrationEvents, a.TxMgr),
		IntegrationProducer: services.NewIntegrationProducerService(a.IntegrationOutbox, a.TxMgr),
	}
}
