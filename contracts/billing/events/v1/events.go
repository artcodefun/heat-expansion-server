package v1

import (
	"github.com/artcodefun/heat-expansion-server/contracts/billing/events"
	"github.com/google/uuid"
)

func init() {
	events.RegisterPayload(EventCrystalsPurchasedV1, func() events.IntegrationEventPayload {
		return &CrystalsPurchasedV1{}
	})
}

const (
	EventCrystalsPurchasedV1 = "billing.crystals.purchased.v1"
)

// CrystalsPurchasedV1 is an integration event payload emitted when a player successfully purchases a crystal package.
type CrystalsPurchasedV1 struct {
	UserID    uuid.UUID `json:"user_id"`
	OrderID   uuid.UUID `json:"order_id"`
	PackageID uuid.UUID `json:"package_id"`
	Crystals  int       `json:"crystals"`
}

func (e CrystalsPurchasedV1) IntegrationEventType() string {
	return EventCrystalsPurchasedV1
}
