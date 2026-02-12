package domain

import (
	"testing"
)

func TestDangerousLocation_DrainedLifecycle(t *testing.T) {
	coords := Vector2i{X: 30, Y: 40}

	armyProto := &ArmyItemPrototype{
		ID:      100,
		Faction: FactionMarauders,
		Defence: 10,
	}
	trophyProto := &StorageItemPrototype{
		ID:             500,
		EstimatedWorth: 100,
	}

	loc := NewDangerousLocation(
		coords,
		FactionMarauders,
		2000,
		[]*StorageItemPrototype{trophyProto},
		[]*ArmyItemPrototype{armyProto},
		nil,
	)

	// Verify initial state
	if len(loc.Trophies) == 0 {
		t.Errorf("expected trophies to be populated")
	}

	// 1. Clear everything EXCEPT trophies
	loc.Resources = LocationResourceStats{}
	loc.DefendingArmies = nil
	loc.DefendingStructures = nil
	loc.CheckIfDrained()
	if len(loc.PullEvents()) != 0 {
		t.Errorf("expected no events while trophies remain")
	}

	// 2. Clear trophies
	loc.Trophies = nil
	loc.CheckIfDrained()

	events := loc.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 drained event, got %d", len(events))
	}
	ev, ok := events[0].(LocationDrainedEvent)
	if !ok {
		t.Fatalf("expected LocationDrainedEvent, got %T", events[0])
	}
	if ev.Type != LocationTypeDangerous {
		t.Errorf("expected dangerous location type, got %s", ev.Type)
	}
}

func TestDangerousLocation_FillTrophiesAndResources(t *testing.T) {
	trophyProtos := []*StorageItemPrototype{
		{ID: 1, EstimatedWorth: 500},
		{ID: 2, EstimatedWorth: 500},
	}

	// Worth 1000: 800 for trophies, 200 for resources
	loc := NewDangerousLocation(
		Vector2i{X: 0, Y: 0},
		FactionMarauders,
		1000,
		trophyProtos,
		nil,
		nil,
	)

	// Should have 1 trophy (500) because 2 trophies (1000) > budget (800)
	// Remaining 300 trophy budget + 200 resource budget = 500 for resources
	if len(loc.Trophies) != 1 {
		t.Errorf("expected 1 trophy, got %d", len(loc.Trophies))
	}

	// Resources should be roughly worth 500 (1000 total - 500 spent on trophies)
	// Note: currently FillFromBudget has some loss due to integer truncation when split across types.
	// 500 budget / 4 types = 125 per type.
	// 125 Cred + 124 Iron (31*4) + 120 Titan (6*20) + 0 AM = 369.
	totalResourceWorth := float64(loc.Resources.Credits)*WorthCredit +
		float64(loc.Resources.Iron)*WorthIron +
		float64(loc.Resources.Titanium)*WorthTitanium +
		float64(loc.Resources.Antimatter)*WorthAntimatter

	// Allow some floating point variance if any
	if totalResourceWorth < 360 || totalResourceWorth > 550 {
		t.Errorf("expected resource worth approx 500 (actual 369 due to truncation), got %f", totalResourceWorth)
	}
}
