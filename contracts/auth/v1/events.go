package v1

import (
	"github.com/artcodefun/heat-expansion-server/contracts/auth"
	"github.com/google/uuid"
)

func init() {
	auth.RegisterPayload(EventAccountRegisteredV1, func() auth.IntegrationEventPayload {
		return &AccountRegisteredV1{}
	})
}

const (
	EventAccountRegisteredV1 = "auth.account.registered.v1"
)

// AccountRegisteredV1 is an integration event payload emitted when a new account registers.
type AccountRegisteredV1 struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
}

func (e AccountRegisteredV1) IntegrationEventType() string {
	return EventAccountRegisteredV1
}
