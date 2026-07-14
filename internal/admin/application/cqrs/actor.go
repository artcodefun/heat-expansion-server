package cqrs

import "github.com/google/uuid"

// Actor carries the caller's identity for app-layer authorization.
// Unauthenticated calls (e.g. Register, Login) receive a zero-value Actor.
type Actor struct {
	AdminID uuid.UUID
}
