package bootstrap

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/application/services"
)

// AppServices aggregates application-level services that are shared across
// commands, queries and other bootstrap wiring.
type AppServices struct {
	Provisioner *services.SectorProvisioningService
	Access      *services.AccessControlService
	Outbox      *services.OutboxService
}

// NewAppServices constructs all application-level services using the provided
// secondary adapters.
func NewAppServices(a *Adapters) *AppServices {
	provisioner := services.NewSectorProvisioningService(a.Content)
	access := services.NewAccessControlService(a.UserBases)
	outbox := services.NewOutboxService(a.OutboxEvents, a.Events, a.TxMgr)

	return &AppServices{
		Provisioner: provisioner,
		Access:      access,
		Outbox:      outbox,
	}
}
