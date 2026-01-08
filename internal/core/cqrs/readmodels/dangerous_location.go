package readmodels

// DangerousLocationModel represents a dangerous site in a sector.
type DangerousLocationModel struct {
	ID          int
	Coordinates Vector2i
	LocationDetails
	DangerLevel int
	// ...other dangerous location-specific fields

	// Available resources at this location (lootable pool)
	Resources LocationResourceStats

	// Defending forces snapshot
	DefendingArmies     []MilitaryUnitSnap     // units guarding this location
	DefendingStructures []DefenseStructureSnap // static defenses like turrets/shields
}
