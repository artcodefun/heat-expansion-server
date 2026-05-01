package domain

import (
	"math/rand"
)

// DangerousLocationModel represents a dangerous site in a sector.
type DangerousLocationModel struct {
	EventProducer
	ID              int
	Coordinates     Vector2i
	DefenderFaction Faction
	LocationDetails
	TotalWorth int // Rough worth in credits
	// ...other dangerous location-specific fields

	// Available resources at this location (lootable pool)
	Resources LocationResourceStats

	// Defenders as stacks (composition). Snapshots are materialized for battles.
	DefendingArmies     []ArmyStack
	DefendingStructures []DefenseStack

	// Trophies granted to the winner when this location is cleared.
	Trophies []TrophyStorageItem
}

// NewDangerousLocation creates and initializes a DangerousLocationModel.
// totalWorth drives the lootable trophy/resource pool; defenseTarget is the desired
// defender combat power and is independent of worth.
func NewDangerousLocation(
	coords Vector2i,
	faction Faction,
	totalWorth int,
	defenseTarget float64,
	storageProtos []*StorageItemPrototype,
	armyProtos []*ArmyItemPrototype,
	buildProtos []*BuildItemPrototype,
) *DangerousLocationModel {
	loc := &DangerousLocationModel{
		Coordinates:     coords,
		DefenderFaction: faction,
		TotalWorth:      totalWorth,
	}
	loc.FillTrophiesAndResources(storageProtos)
	loc.FillDefenders(defenseTarget, armyProtos, buildProtos)
	return loc
}

// CheckIfDrained checks if resources, trophies, and defenders are all empty and emits an event if so.
func (dl *DangerousLocationModel) CheckIfDrained() {
	if !dl.Resources.IsEmpty() || len(dl.Trophies) > 0 {
		return
	}
	for _, stack := range dl.DefendingArmies {
		if stack.Count > 0 {
			return
		}
	}
	for _, stack := range dl.DefendingStructures {
		if stack.Count > 0 {
			return
		}
	}
	dl.AddEvent(NewLocationDrainedEvent(dl.Coordinates.X, dl.Coordinates.Y, LocationTypeDangerous))
}

// TrophyStorageItem represents a special item granted when a dangerous location is defeated.
type TrophyStorageItem struct {
	PrototypeID int
}

// FillTrophiesAndResources populates the location with trophies from a pool and fills the rest with resources.
func (dl *DangerousLocationModel) FillTrophiesAndResources(availableTrophies []*StorageItemPrototype) {
	totalBudget := float64(dl.TotalWorth)
	trophyBudget := totalBudget * 0.8
	resourceBudget := totalBudget * 0.2

	spentOnTrophies := 0.0

	// 1. Shuffle to ensure variety and prevent cheap buffs from eating the whole budget
	shuffled := make([]*StorageItemPrototype, len(availableTrophies))
	copy(shuffled, availableTrophies)
	r := rand.New(rand.NewSource(NowUnixNano()))
	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	for _, p := range shuffled {
		worth := float64(p.EstimatedWorth)
		if worth > 0 && worth <= (trophyBudget-spentOnTrophies) {
			dl.Trophies = append(dl.Trophies, TrophyStorageItem{
				PrototypeID: p.ID,
			})
			spentOnTrophies += worth
		}
	}

	// 2. Any remainder from the trophy budget is added to the resource budget (20%)
	remainingTrophyBudget := trophyBudget - spentOnTrophies
	totalResourceBudget := resourceBudget + remainingTrophyBudget

	// 3. Fill the resources using the helper
	dl.Resources.FillFromBudget(totalResourceBudget, "", 0)
}

// FillDefenders populates the location with defenders targeting defenseTarget combat power.
func (dl *DangerousLocationModel) FillDefenders(defenseTarget float64, armyProtos []*ArmyItemPrototype, buildProtos []*BuildItemPrototype) {
	FillDefenders(&dl.DefendingArmies, &dl.DefendingStructures, dl.DefenderFaction, defenseTarget, armyProtos, buildProtos)
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
		out = append(out, MilitaryUnitFromStack(s))
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
		out = append(out, DefenseStructureFromStack(s))
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

	dl.CheckIfDrained()
}
