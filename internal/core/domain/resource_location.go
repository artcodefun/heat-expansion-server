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

	// Defenders as stacks (composition). Snapshots are materialized for battles.
	DefendingArmies     []ArmyStack
	DefendingStructures []DefenseStack
}

// MaterializeDefenderArmySnapshot builds battle-ready snapshots using current prototype values.
func (rl *ResourceLocationModel) MaterializeDefenderArmySnapshot() []MilitaryUnitSnap {
	if len(rl.DefendingArmies) == 0 {
		return nil
	}
	out := make([]MilitaryUnitSnap, 0, len(rl.DefendingArmies))
	for _, s := range rl.DefendingArmies {
		if s.Count <= 0 {
			continue
		}
		out = append(out, s.ToSnap())
	}
	return out
}

// MaterializeDefenderStructureSnapshot builds battle-ready snapshots for static defenses.
func (rl *ResourceLocationModel) MaterializeDefenderStructureSnapshot() []DefenseStructureSnap {
	if len(rl.DefendingStructures) == 0 {
		return nil
	}
	out := make([]DefenseStructureSnap, 0, len(rl.DefendingStructures))
	for _, s := range rl.DefendingStructures {
		if s.Count <= 0 {
			continue
		}
		out = append(out, s.ToSnap())
	}
	return out
}

// ApplyDefenderArmyRemaining sets the defending units to the provided remaining snapshot.
func (rl *ResourceLocationModel) ApplyDefenderArmyRemaining(remaining []MilitaryUnitSnap) {
	remByProto := make(map[int]int, len(remaining))
	for _, u := range remaining {
		remByProto[u.PrototypeID] += u.Count
	}

	newArmies := make([]ArmyStack, 0, len(rl.DefendingArmies))
	for _, s := range rl.DefendingArmies {
		if count, ok := remByProto[s.Prototype.ID]; ok && count > 0 {
			s.Count = count
			newArmies = append(newArmies, s)
		}
	}
	rl.DefendingArmies = newArmies
}

// ApplyRemainingDefensiveStructures sets the defensive structures to the provided remaining snapshot.
func (rl *ResourceLocationModel) ApplyRemainingDefensiveStructures(remaining []DefenseStructureSnap) {
	remByProto := make(map[int]int, len(remaining))
	for _, s := range remaining {
		remByProto[s.PrototypeID] += s.Count
	}

	newStructures := make([]DefenseStack, 0, len(rl.DefendingStructures))
	for _, s := range rl.DefendingStructures {
		if count, ok := remByProto[s.Prototype.ID]; ok && count > 0 {
			s.Count = count
			newStructures = append(newStructures, s)
		}
	}
	rl.DefendingStructures = newStructures
}

// DeductLoot subtracts the provided loot from the location's resource pool, clamped at zero.
func (rl *ResourceLocationModel) DeductLoot(loot PriceModel) {
	rl.Resources.Credits = maxInt(rl.Resources.Credits-loot.Credits, 0)
	rl.Resources.Iron = maxInt(rl.Resources.Iron-loot.Iron, 0)
	rl.Resources.Titanium = maxInt(rl.Resources.Titanium-loot.Titanium, 0)
	rl.Resources.Antimatter = maxInt(rl.Resources.Antimatter-loot.Antimatter, 0)
}
