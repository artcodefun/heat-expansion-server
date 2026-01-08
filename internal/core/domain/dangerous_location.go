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

	// Defenders as stacks (composition). Snapshots are materialized for battles.
	DefendingArmies     []ArmyStack
	DefendingStructures []DefenseStack
}

// MaterializeDefenderArmySnapshot builds battle-ready snapshots using current prototype values.
func (dl *DangerousLocationModel) MaterializeDefenderArmySnapshot() []MilitaryUnitSnap {
	if len(dl.DefendingArmies) == 0 {
		return nil
	}
	out := make([]MilitaryUnitSnap, 0, len(dl.DefendingArmies))
	for _, s := range dl.DefendingArmies {
		if s.Count <= 0 {
			continue
		}
		out = append(out, s.ToSnap())
	}
	return out
}

// MaterializeDefenderStructureSnapshot builds battle-ready snapshots for static defenses.
func (dl *DangerousLocationModel) MaterializeDefenderStructureSnapshot() []DefenseStructureSnap {
	if len(dl.DefendingStructures) == 0 {
		return nil
	}
	out := make([]DefenseStructureSnap, 0, len(dl.DefendingStructures))
	for _, s := range dl.DefendingStructures {
		if s.Count <= 0 {
			continue
		}
		out = append(out, s.ToSnap())
	}
	return out
}

// ApplyDefenderArmyRemaining sets the defending units to the provided remaining snapshot.
func (dl *DangerousLocationModel) ApplyDefenderArmyRemaining(remaining []MilitaryUnitSnap) {
	remByProto := make(map[int]int, len(remaining))
	for _, u := range remaining {
		remByProto[u.PrototypeID] += u.Count
	}

	newArmies := make([]ArmyStack, 0, len(dl.DefendingArmies))
	for _, s := range dl.DefendingArmies {
		if count, ok := remByProto[s.Prototype.ID]; ok && count > 0 {
			s.Count = count
			newArmies = append(newArmies, s)
		}
	}
	dl.DefendingArmies = newArmies
}

// ApplyRemainingDefensiveStructures sets the defensive structures to the provided remaining snapshot.
func (dl *DangerousLocationModel) ApplyRemainingDefensiveStructures(remaining []DefenseStructureSnap) {
	remByProto := make(map[int]int, len(remaining))
	for _, s := range remaining {
		remByProto[s.PrototypeID] += s.Count
	}

	newStructures := make([]DefenseStack, 0, len(dl.DefendingStructures))
	for _, s := range dl.DefendingStructures {
		if count, ok := remByProto[s.Prototype.ID]; ok && count > 0 {
			s.Count = count
			newStructures = append(newStructures, s)
		}
	}
	dl.DefendingStructures = newStructures
}

// DeductLoot subtracts the provided loot from the location's resource pool, clamped at zero.
func (dl *DangerousLocationModel) DeductLoot(loot PriceModel) {
	dl.Resources.Credits = maxInt(dl.Resources.Credits-loot.Credits, 0)
	dl.Resources.Iron = maxInt(dl.Resources.Iron-loot.Iron, 0)
	dl.Resources.Titanium = maxInt(dl.Resources.Titanium-loot.Titanium, 0)
	dl.Resources.Antimatter = maxInt(dl.Resources.Antimatter-loot.Antimatter, 0)
}
