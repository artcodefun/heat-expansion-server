package domain

import (
	"github.com/google/uuid"
)

// User represents a player in the game.
type User struct {
	EventProducer
	ID       uuid.UUID
	Name     string
	Crystals int // Global in-game currency for the user
}

// Default values for new users
const (
	DefaultCrystalsBalance = 50
)

// NewUser creates a new user with default settings and a creation event.
func NewUser(id uuid.UUID, name string) *User {
	u := &User{
		ID:       id,
		Name:     name,
		Crystals: DefaultCrystalsBalance,
	}
	u.AddEvent(NewUserAccountCreatedEvent(u.ID))
	return u
}

// SpendCrystals deducts the given amount from the user's crystal balance.
// It returns an error if the amount is non-positive or if the user does not
// have enough crystals available.
func (u *User) SpendCrystals(amount int) error {
	if amount <= 0 {
		return NewError("error.domain.user.invalid_crystal_spend_amount", H{"amount": amount})
	}
	if u.Crystals < amount {
		return NewError("error.domain.user.not_enough_crystals", nil)
	}
	u.Crystals -= amount
	return nil
}

// AddCrystals credits the given amount to the user's crystal balance.
// It returns an error if amount is non-positive.
func (u *User) AddCrystals(amount int) error {
	if amount <= 0 {
		return NewError("error.domain.user.invalid_crystal_add_amount", H{"amount": amount})
	}
	u.Crystals += amount
	return nil
}
