package domain

import "fmt"

// User represents a player in the game.
type User struct {
	EventProducer
	ID           int
	Name         string
	Email        string
	PasswordHash string
	Crystals     int // Global in-game currency for the user
}

// Default values for new users
const (
	DefaultCrystalsBalance = 50
)

// Initialize sets default values and emits user created event.
func (u *User) Initialize() {
	u.Crystals = DefaultCrystalsBalance
	u.AddEvent(NewUserAccountCreatedEvent(u.ID))
}

// SpendCrystals deducts the given amount from the user's crystal balance.
// It returns an error if the amount is non-positive or if the user does not
// have enough crystals available.
func (u *User) SpendCrystals(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("invalid crystal spend amount: %d", amount)
	}
	if u.Crystals < amount {
		return fmt.Errorf("not enough crystals")
	}
	u.Crystals -= amount
	return nil
}
