package domain

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
	DefaultCrystalsBalance = 5
)

// Initialize sets default values and emits user created event.
func (u *User) Initialize() {
	u.Crystals = DefaultCrystalsBalance
	u.AddEvent(NewUserAccountCreatedEvent(u.ID))
}
