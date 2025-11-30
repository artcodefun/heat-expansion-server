package readmodels

// MilitaryUnit captures a snapshot of an army unit participating in an operation.
// Using a snapshot decouples battle resolution from prototype changes over time.
type MilitaryUnit struct {
	PrototypeID int
	Category    ArmyCategory
	Attack      int
	Defence     int
	Capacity    int
	Stealth     int
	Speed       int
	Count       int
}

// DefenseStructure is a simplified snapshot of a defending structure (e.g., turrets, shields).
// This provides a lightweight input for attack resolution without coupling to build prototypes.
type DefenseStructure struct {
	PrototypeID int
	Defence     int
	Count       int
}
