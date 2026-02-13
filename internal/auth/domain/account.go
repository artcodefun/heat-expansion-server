package domain

import (
	"github.com/google/uuid"
)

// Account represents a user identity in the authentication service.
type Account struct {
	EventProducer
	ID           uuid.UUID
	Name         string
	Email        string
	PasswordHash string
}

// RegisterAccount is a factory function that creates a new Account
// and emits an AccountRegisteredEvent.
func RegisterAccount(name, email, passwordHash string) *Account {
	id, _ := uuid.NewV7()
	a := &Account{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	}

	a.AddEvent(NewAccountRegisteredEvent(a.ID, a.Name, a.Email))

	return a
}
