package domain

import (
	"testing"
)

func TestBuilding_MoveAndSpeedUpProduction_EmitsEvents(t *testing.T) {
	SetTestNow(t, 1_000)
	base := newBaseWithDefaults(1)

	proto := &BuildItemPrototype{
		ID:             1,
		Name:           "Mine",
		Category:       BuildCategoryResources,
		Price:          PriceModel{Credits: 100, Iron: 50, Titanium: 30, Antimatter: 20},
		ProductionTime: 120,
		Space:          2,
	}

	// queue -> should start production immediately via MoveBuildQueue in AddToBuildQueue
	if err := base.AddToBuildQueue(proto); err != nil {
		t.Fatalf("AddToBuildQueue error: %v", err)
	}
	if got := len(base.BuildingsInProduction); got != 1 {
		t.Fatalf("expected 1 building in production, got %d", got)
	}
	// resources should be debited exactly by prototype price
	if base.Stats.Credits != 10_000-100 || base.Stats.Iron != 10_000-50 || base.Stats.Titanium != 10_000-30 || base.Stats.Antimatter != 10_000-20 {
		t.Fatalf("unexpected stats after AddToBuildQueue: %+v", base.Stats)
	}
	// space usage should reflect the building's space
	if base.Stats.Space != proto.Space {
		t.Fatalf("expected space=%d after queuing building, got %d", proto.Space, base.Stats.Space)
	}

	// started event should reference the in-production item
	prodID := base.BuildingsInProduction[0].ID
	events := base.PullEvents()
	var started BuildingProductionStartedEvent
	foundStarted := false
	for _, e := range events {
		if ev, ok := e.(BuildingProductionStartedEvent); ok {
			started = ev
			foundStarted = true
			break
		}
	}
	if !foundStarted {
		t.Fatalf("expected BuildingProductionStartedEvent after AddToBuildQueue")
	}
	if started.BaseID != base.ID || started.ItemID != prodID {
		t.Fatalf("unexpected BuildingProductionStartedEvent payload: %+v", started)
	}

	// advance time to complete and move queue
	SetTestNow(t, 1_000+200)
	base.MoveBuildQueue()
	if len(base.BuildingsInProduction) != 0 {
		t.Fatalf("expected no buildings in production after completion")
	}
	if len(base.BuildingsPresent) != 1 {
		t.Fatalf("expected 1 present building after completion")
	}
	// completed building should match the prototype
	if base.BuildingsPresent[0].Prototype.ID != proto.ID {
		t.Fatalf("unexpected present building prototype after completion: %+v", base.BuildingsPresent[0].Prototype)
	}
	events = base.PullEvents()
	if len(events) < 1 {
		t.Fatalf("expected at least 1 event, got %d", len(events))
	}
	// last event should be BuildingProductionFinished for the completed item
	foundFinished := false
	for _, e := range events {
		if _, ok := e.(BuildingProductionFinishedEvent); ok {
			foundFinished = true
			break
		}
	}
	if !foundFinished {
		t.Fatalf("expected BuildingProductionFinishedEvent to be emitted")
	}

	// Now test speedup path: queue another building and speed it up
	SetTestNow(t, 2_000)
	if err := base.AddToBuildQueue(proto); err != nil {
		t.Fatalf("AddToBuildQueue (2) error: %v", err)
	}
	if len(base.BuildingsInProduction) != 1 {
		t.Fatalf("expected 1 building in production before speedup")
	}
	inProdID := base.BuildingsInProduction[0].ID

	// speed up should not panic and should emit finished + speedup events
	base.PullEvents() // clear
	if err := base.SpeedUpBuildingProduction(inProdID); err != nil {
		t.Fatalf("SpeedUpBuildingProduction error: %v", err)
	}
	if len(base.BuildingsInProduction) != 0 {
		t.Fatalf("expected no buildings in production after speedup")
	}
	events = base.PullEvents()
	var gotFinished, gotSpeedup bool
	for _, e := range events {
		switch e.(type) {
		case BuildingProductionFinishedEvent:
			gotFinished = true
		case BuildingProductionSpeedupEvent:
			gotSpeedup = true
		}
	}
	if !gotFinished || !gotSpeedup {
		t.Fatalf("expected finished and speedup events, got finished=%v speedup=%v (events=%T)", gotFinished, gotSpeedup, events)
	}
}

func TestBuilding_AddToBuildQueue_NotEnoughSpace(t *testing.T) {
	SetTestNow(t, 2_000)
	base := newBaseWithDefaults(10)
	// artificially restrict space capacity to simulate a nearly full base
	base.Stats.MaxSpace = 1

	proto := &BuildItemPrototype{
		ID:             2,
		Name:           "Big Tower",
		Category:       BuildCategoryResources,
		Price:          PriceModel{Credits: 100, Iron: 50, Titanium: 30, Antimatter: 20},
		ProductionTime: 60,
		Space:          2, // exceeds capacity
	}

	if err := base.AddToBuildQueue(proto); err == nil {
		t.Fatalf("expected error when queuing building without enough space")
	}
	// no items should be queued or in production and resources must be unchanged
	if len(base.BuildingsPending) != 0 || len(base.BuildingsInProduction) != 0 {
		t.Fatalf("expected no buildings queued or in production after space error")
	}
	if base.Stats.Credits != 10_000 || base.Stats.Iron != 10_000 || base.Stats.Titanium != 10_000 || base.Stats.Antimatter != 10_000 {
		t.Fatalf("expected resources to remain unchanged after space error, got %+v", base.Stats)
	}
}

func TestApplyRemainingDefensiveStructures_KeepsNonDefensiveAndAppliesCounts(t *testing.T) {
	base := newBaseWithDefaults(41)

	// One non-defensive building, and three defensive turrets (same prototype ID)
	defProto := BuildItemPrototype{ID: 400, Category: BuildCategoryDefense, DefenseData: &DefenseBuildingData{DefenceBonus: 5}}
	base.BuildingsPresent = []BuildItemPresent{
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: BuildItemPrototype{ID: 399, Category: BuildCategoryResources}},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: defProto},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: defProto},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: defProto},
	}

	// Remaining structures say keep only 2 turrets of this prototype
	remaining := []DefenseStructureSnap{{PrototypeID: 400, Count: 2}}
	base.ApplyRemainingDefensiveStructures(remaining)

	var nonDefCount, defCount int
	for _, b := range base.BuildingsPresent {
		if b.Prototype.DefenseData == nil {
			nonDefCount++
		} else if b.Prototype.ID == 400 {
			defCount++
		}
	}
	if nonDefCount != 1 || defCount != 2 {
		t.Fatalf("unexpected defensive structure filtering: nonDef=%d def=%d", nonDefCount, defCount)
	}
}

func TestUserBase_TotalRadarStealthStrength(t *testing.T) {
	base := newBaseWithDefaults(1)
	if base.TotalRadarStealthStrength() != 0 {
		t.Errorf("expected 0 initial radar strength, got %d", base.TotalRadarStealthStrength())
	}

	radarProto := &BuildItemPrototype{
		ID:       10,
		Category: BuildCategoryIntelligence,
		IntelligenceData: &IntelligenceBuildingData{
			Subtype:         IntelligenceSubtypeRadar,
			StealthStrength: 50,
			ScanRange:       100,
		},
	}

	base.BuildingsPresent = append(base.BuildingsPresent, BuildItemPresent{
		Prototype: *radarProto,
	})

	base.recalculateStats()

	if base.TotalRadarStealthStrength() != 50 {
		t.Errorf("expected 50 radar strength, got %d", base.TotalRadarStealthStrength())
	}

	// Add another one
	base.BuildingsPresent = append(base.BuildingsPresent, BuildItemPresent{
		Prototype: *radarProto,
	})

	base.recalculateStats()

	if base.TotalRadarStealthStrength() != 100 {
		t.Errorf("expected 100 radar strength, got %d", base.TotalRadarStealthStrength())
	}
}
