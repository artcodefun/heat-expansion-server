package domain

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
	DefendingUnits []MilitaryUnit     // units guarding this location
	Structures     []DefenseStructure // static defenses like turrets/shields
}

// ApplyDefenderArmyRemaining sets the defending units to the provided remaining snapshot.
func (dl *DangerousLocationModel) ApplyDefenderArmyRemaining(remaining []MilitaryUnit) {
	dl.DefendingUnits = CloneOperationUnits(remaining)
}

// ApplyRemainingDefensiveStructures sets the defensive structures to the provided remaining snapshot.
func (dl *DangerousLocationModel) ApplyRemainingDefensiveStructures(remaining []DefenseStructure) {
	dl.Structures = CloneDefenseStructures(remaining)
}

// DeductLoot subtracts the provided loot from the location's resource pool, clamped at zero.
func (dl *DangerousLocationModel) DeductLoot(loot PriceModel) {
	dl.Resources.Credits = maxInt(dl.Resources.Credits-loot.Credits, 0)
	dl.Resources.Iron = maxInt(dl.Resources.Iron-loot.Iron, 0)
	dl.Resources.Titanium = maxInt(dl.Resources.Titanium-loot.Titanium, 0)
	dl.Resources.Antimatter = maxInt(dl.Resources.Antimatter-loot.Antimatter, 0)
}
