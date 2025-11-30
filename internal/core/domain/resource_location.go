package domain

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

// ApplyDefenderArmyRemaining sets the defending units to the provided remaining snapshot.
func (rl *ResourceLocationModel) ApplyDefenderArmyRemaining(remaining []MilitaryUnit) {
	rl.DefendingUnits = CloneOperationUnits(remaining)
}

// ApplyRemainingDefensiveStructures sets the defensive structures to the provided remaining snapshot.
func (rl *ResourceLocationModel) ApplyRemainingDefensiveStructures(remaining []DefenseStructure) {
	rl.Structures = CloneDefenseStructures(remaining)
}

// DeductLoot subtracts the provided loot from the location's resource pool, clamped at zero.
func (rl *ResourceLocationModel) DeductLoot(loot PriceModel) {
	rl.Resources.Credits = maxInt(rl.Resources.Credits-loot.Credits, 0)
	rl.Resources.Iron = maxInt(rl.Resources.Iron-loot.Iron, 0)
	rl.Resources.Titanium = maxInt(rl.Resources.Titanium-loot.Titanium, 0)
	rl.Resources.Antimatter = maxInt(rl.Resources.Antimatter-loot.Antimatter, 0)
}
