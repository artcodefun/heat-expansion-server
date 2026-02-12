package readmodels

// MilitaryUnitSnap captures a snapshot of an army unit participating in an operation.
// Using a snapshot decouples battle resolution from prototype changes over time.
type MilitaryUnitSnap struct {
	PrototypeID      int
	CurrentPrototype ArmyItemPrototype
	Name             string
	ImageURL         string
	Category         ArmyCategory
	Attack           int
	Defence          int
	Capacity         int
	Stealth          int
	Speed            int
	Space            int
	Count            int
}

// DefenseStructureSnap is a simplified snapshot of a defending structure (e.g., turrets, shields).
// This provides a lightweight input for attack resolution without coupling to build prototypes.
type DefenseStructureSnap struct {
	PrototypeID      int
	CurrentPrototype BuildItemPrototype
	Name             string
	ImageURL         string
	Defence          int
	Space            int
	Count            int
}
