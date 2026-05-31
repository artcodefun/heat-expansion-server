package domain

import "github.com/google/uuid"

// User is a local projection of an auth account, populated by consuming
// auth.account.registered.v1 integration events. Billing keeps it solely so it
// can attach a customer email to payment receipts. It is not an aggregate and
// emits no domain events.
type User struct {
	ID    uuid.UUID
	Email string
}
