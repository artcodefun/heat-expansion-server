package v1

import "github.com/google/uuid"

const EventAccountRegisteredV1 = "auth.account.registered.v1"

// AccountRegisteredV1 is an integration event payload emitted when a new account registers.
type AccountRegisteredV1 struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Email  string    `json:"email"`
}
