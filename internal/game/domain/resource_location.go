package domain

// ResourceType represents the primary resource found at a location.
type ResourceType string

const (
	ResourceTypeCredits    ResourceType = "CREDITS"
	ResourceTypeIron       ResourceType = "IRON"
	ResourceTypeTitanium   ResourceType = "TITANIUM"
	ResourceTypeAntimatter ResourceType = "ANTIMATTER"
)

// ResourceLocationModel represents a resource site in a sector.
type ResourceLocationModel struct {
	EventProducer
	ID              int
	Coordinates     Vector2i
	DefenderFaction Faction
	LocationDetails
	Type       ResourceType
	TotalWorth int // Rough worth in credits

	// Available resources at this location (lootable pool)
	Resources LocationResourceStats

	// Defenders as stacks (composition). Snapshots are materialized for battles.
	DefendingArmies     []ArmyStack
	DefendingStructures []DefenseStack
}

// NewResourceLocation creates and initializes a ResourceLocationModel.
func NewResourceLocation(
	coords Vector2i,
	resType ResourceType,
	faction Faction,
	totalWorth int,
	armyProtos []*ArmyItemPrototype,
	buildProtos []*BuildItemPrototype,
) *ResourceLocationModel {
	loc := &ResourceLocationModel{
		Coordinates:     coords,
		Type:            resType,
		DefenderFaction: faction,
		TotalWorth:      totalWorth,
	}
	loc.FillResources()
	loc.FillDefenders(armyProtos, buildProtos)
	return loc
}

// CheckIfDrained checks if resources are empty and all defenders are defeated and emits an event if so.
func (rl *ResourceLocationModel) CheckIfDrained() {
	if !rl.Resources.IsEmpty() {
		return
	}
	for _, stack := range rl.DefendingArmies {
		if stack.Count > 0 {
			return
		}
	}
	for _, stack := range rl.DefendingStructures {
		if stack.Count > 0 {
			return
		}
	}
	rl.AddEvent(NewLocationDrainedEvent(rl.Coordinates.X, rl.Coordinates.Y, LocationTypeResourceful))
}

// FillResources populates the location's resource pool based on TotalWorth and ResourceType.
// FillResources populates the location with resources based on its TotalWorth.
func (dl *ResourceLocationModel) FillResources() {
	dl.Resources.FillFromBudget(float64(dl.TotalWorth), dl.Type, 0.7)
}

// FillDefenders populates the location with defenders based on its TotalWorth and Faction.
func (rl *ResourceLocationModel) FillDefenders(armyProtos []*ArmyItemPrototype, buildProtos []*BuildItemPrototype) {
	FillDefenders(&rl.DefendingArmies, &rl.DefendingStructures, rl.DefenderFaction, rl.TotalWorth, armyProtos, buildProtos)
}

// MaterializeDefenderArmySnapshot builds battle-ready snapshots using current prototype values.
func (dl *ResourceLocationModel) MaterializeDefenderArmySnapshot() []MilitaryUnitSnap {
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
func (dl *ResourceLocationModel) MaterializeDefenderStructureSnapshot() []DefenseStructureSnap {
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
func (dl *ResourceLocationModel) ApplyDefenderArmyRemaining(remaining []MilitaryUnitSnap) {
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
func (dl *ResourceLocationModel) ApplyRemainingDefensiveStructures(remaining []DefenseStructureSnap) {
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
func (dl *ResourceLocationModel) DeductLoot(loot PriceModel) {
	dl.Resources.Credits = maxInt(dl.Resources.Credits-loot.Credits, 0)
	dl.Resources.Iron = maxInt(dl.Resources.Iron-loot.Iron, 0)
	dl.Resources.Titanium = maxInt(dl.Resources.Titanium-loot.Titanium, 0)
	dl.Resources.Antimatter = maxInt(dl.Resources.Antimatter-loot.Antimatter, 0)

	dl.CheckIfDrained()
}
