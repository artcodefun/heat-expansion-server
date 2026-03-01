package cqrs

import "github.com/google/uuid"

// Actor carries caller identity and authorization scope for app-layer checks.
type Actor struct {
	UserID uuid.UUID
	Roles  []string
}
