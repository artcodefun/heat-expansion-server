package readmodels

// User represents a player in the game.
type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
	Crystals     int // Global in-game currency for the user
}
