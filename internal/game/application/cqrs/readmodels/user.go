package readmodels

import "github.com/google/uuid"

// User represents a player in the game.
type User struct {
	ID       uuid.UUID
	Name     string
	Crystals int // Global in-game currency for the user
}
