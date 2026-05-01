package domain

import (
	"testing"
)

func TestResourceLocation_DrainedLifecycle(t *testing.T) {
	coords := Vector2i{X: 10, Y: 20}

	// Create sample prototypes
	armyProto := &ArmyItemPrototype{
		ID:      100,
		Faction: FactionMarauders,
		Defence: 10,
	}
	buildProto := &BuildItemPrototype{
		ID:          200,
		Faction:     FactionMarauders,
		DefenseData: &DefenseBuildingData{DefenceBonus: 20},
	}

	loc := NewResourceLocation(
		coords,
		ResourceTypeIron,
		FactionMarauders,
		1000,
		20.0,
		[]*ArmyItemPrototype{armyProto},
		[]*BuildItemPrototype{buildProto},
	)

	// Verify initial state
	if loc.Resources.IsEmpty() {
		t.Errorf("expected resources to be populated")
	}
	if len(loc.DefendingArmies) == 0 && len(loc.DefendingStructures) == 0 {
		t.Errorf("expected defenders to be populated")
	}

	// 1. Partial drain (loot only)
	loot := PriceModel{Iron: 1}
	loc.DeductLoot(loot)
	if len(loc.PullEvents()) != 0 {
		t.Errorf("expected no events after partial loot")
	}

	// 2. Clear all resources
	loc.Resources = LocationResourceStats{}
	loc.CheckIfDrained()
	if len(loc.PullEvents()) != 0 {
		t.Errorf("expected no events while defenders remain")
	}

	// 3. Clear all defenders
	loc.DefendingArmies = nil
	loc.DefendingStructures = nil
	loc.CheckIfDrained()

	events := loc.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 drained event, got %d", len(events))
	}
	ev, ok := events[0].(LocationDrainedEvent)
	if !ok {
		t.Fatalf("expected LocationDrainedEvent, got %T", events[0])
	}
	if ev.X != coords.X || ev.Y != coords.Y || ev.Type != LocationTypeResourceful {
		t.Errorf("unexpected event payload: %+v", ev)
	}
}

func TestResourceLocation_ApplyDefendersRemaining(t *testing.T) {
	loc := &ResourceLocationModel{
		DefendingArmies: []ArmyStack{
			{Prototype: ArmyItemPrototype{ID: 1}, Count: 10},
			{Prototype: ArmyItemPrototype{ID: 2}, Count: 5},
		},
	}

	// Only 2 units of ID 1 survived, ID 2 wiped out
	remaining := []MilitaryUnitSnap{
		{PrototypeID: 1, Count: 2},
	}

	loc.ApplyDefenderArmyRemaining(remaining)

	if len(loc.DefendingArmies) != 1 {
		t.Fatalf("expected 1 army stack remaining, got %d", len(loc.DefendingArmies))
	}
	if loc.DefendingArmies[0].Prototype.ID != 1 || loc.DefendingArmies[0].Count != 2 {
		t.Errorf("unexpected remaining army stack: %+v", loc.DefendingArmies[0])
	}
}
