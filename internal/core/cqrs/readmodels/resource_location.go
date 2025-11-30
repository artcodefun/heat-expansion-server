package readmodels

// ResourceLocationModel represents a resource site in a sector.
type ResourceLocationModel struct {
	ID          int
	Coordinates Vector2i
	LocationDetails
	Type   string // e.g., "IRON_MINE", "TITANIUM_FIELD"
	Amount int
	// ...other resource-specific fields

	// Available resources at this location (lootable pool)
	Resources LocationResourceStats

	// Defending forces snapshot
	DefendingUnits []MilitaryUnit     // units guarding this resource location
	Structures     []DefenseStructure // static defenses like turrets/shields
}
