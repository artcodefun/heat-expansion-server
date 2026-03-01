package cqrs

import "github.com/google/uuid"

// Actor carries caller identity & auth scope used by app-layer authorization.
type Actor struct {
	AccountID uuid.UUID
}
