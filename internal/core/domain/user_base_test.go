package domain

import (
	"testing"
)

func newBaseWithDefaults(id int) *UserBaseModel {
	b := &UserBaseModel{ID: id}
	// Zero production rates by default
	b.Stats = UserBaseStats{
		Credits:              10_000,
		Iron:                 10_000,
		Titanium:             10_000,
		Antimatter:           10_000,
		SpaceCapacity:        DefaultSpaceCapacity,
		CalculationTimestamp: NowUnix(),
	}
	return b
}

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

	// advance time to complete and move queue
	SetTestNow(t, 1_000+200)
	base.MoveBuildQueue()
	if len(base.BuildingsInProduction) != 0 {
		t.Fatalf("expected no buildings in production after completion")
	}
	if len(base.BuildingsPresent) != 1 {
		t.Fatalf("expected 1 present building after completion")
	}
	events := base.PullEvents()
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

func TestArmy_MoveAndSpeedUp_EmitsEvents(t *testing.T) {
	SetTestNow(t, 5_000)
	base := newBaseWithDefaults(2)
	// Unlock infantry via present military building
	base.BuildingsPresent = []BuildItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype: BuildItemPrototype{
			ID:           10,
			Name:         "Barracks",
			Category:     BuildCategoryMilitary,
			MilitaryData: &MilitaryBuildingData{UnlockArmyCategory: ArmyCategoryInfantry},
			Space:        0,
		},
	}}

	army := &ArmyItemPrototype{
		ID:             100,
		Name:           "Infantry",
		Category:       ArmyCategoryInfantry,
		Price:          PriceModel{},
		ProductionTime: 60,
		Space:          1,
		Attack:         10,
		Defence:        5,
		Capacity:       2,
		Speed:          100,
	}
	if err := base.QueueArmy(army, 1); err != nil {
		t.Fatalf("QueueArmy error: %v", err)
	}
	// start production
	base.MoveArmyQueue()
	if len(base.ArmiesInProduction) != 1 {
		t.Fatalf("expected 1 army in production after MoveArmyQueue start")
	}
	// started event expected
	events := base.PullEvents()
	hasStarted := false
	for _, e := range events {
		if _, ok := e.(ArmyProductionStartedEvent); ok {
			hasStarted = true
		}
	}
	if !hasStarted {
		t.Fatalf("expected ArmyProductionStartedEvent after starting queue")
	}

	// speed up and expect finished + speedup
	inProdID := base.ArmiesInProduction[0].ID
	base.PullEvents()
	if err := base.SpeedUpArmyProduction(inProdID); err != nil {
		t.Fatalf("SpeedUpArmyProduction error: %v", err)
	}
	events = base.PullEvents()
	var gotFinished, gotSpeedup bool
	for _, e := range events {
		switch e.(type) {
		case ArmyProductionFinishedEvent:
			gotFinished = true
		case ArmyProductionSpeedupEvent:
			gotSpeedup = true
		}
	}
	if !gotFinished || !gotSpeedup {
		t.Fatalf("expected army finished and speedup events, got finished=%v speedup=%v", gotFinished, gotSpeedup)
	}
}

func TestTech_StartAndSpeedUp_EmitsEvents(t *testing.T) {
	SetTestNow(t, 10_000)
	base := newBaseWithDefaults(3)
	tech := &TechItemPrototype{
		ID:           200,
		Name:         "Improved Mining",
		Category:     TechCategoryBuild,
		Price:        PriceModel{},
		ResearchTime: 90,
	}
	if err := base.StartTechResearch(tech); err != nil {
		t.Fatalf("StartTechResearch error: %v", err)
	}
	events := base.PullEvents()
	if len(events) == 0 {
		t.Fatalf("expected TechResearchStartedEvent on start")
	}
	_, ok := events[0].(TechResearchStartedEvent)
	if !ok {
		t.Fatalf("expected first event to be TechResearchStartedEvent, got %T", events[0])
	}
	if got := len(base.TechnologiesInProgress); got != 1 {
		t.Fatalf("expected 1 tech in progress, got %d", got)
	}

	// speed up
	inProgID := base.TechnologiesInProgress[0].BaseOwnedItem.ID
	base.PullEvents()
	if err := base.SpeedUpTechResearch(inProgID); err != nil {
		t.Fatalf("SpeedUpTechResearch error: %v", err)
	}
	events = base.PullEvents()
	var gotFinished, gotSpeedup bool
	for _, e := range events {
		switch e.(type) {
		case TechResearchFinishedEvent:
			gotFinished = true
		case TechResearchSpeedupEvent:
			gotSpeedup = true
		}
	}
	if !gotFinished || !gotSpeedup {
		t.Fatalf("expected tech finished and speedup events, got finished=%v speedup=%v", gotFinished, gotSpeedup)
	}
}

func TestStorage_BuffActivateAndExpire(t *testing.T) {
	SetTestNow(t, 20_000)
	base := newBaseWithDefaults(4)
	// add a buff storage item
	buff := StorageItemPresent{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype: StorageItemPrototype{
			ID:       300,
			Name:     "Space Booster",
			Category: StorageCategoryBuff,
			BuffData: &BuffStorageData{DurationSeconds: 100},
		},
	}
	base.StorageItemsPresent = []StorageItemPresent{buff}

	// Activate
	if err := base.ActivateBuffByID(buff.ID); err != nil {
		t.Fatalf("ActivateBuffByID error: %v", err)
	}
	events := base.PullEvents()
	if len(events) == 0 {
		t.Fatalf("expected BuffActivatedEvent")
	}
	if _, ok := events[0].(BuffActivatedEvent); !ok {
		t.Fatalf("expected BuffActivatedEvent, got %T", events[0])
	}
	if base.StorageItemsPresent[0].Prototype.BuffData.ActivatedAt == nil || *base.StorageItemsPresent[0].Prototype.BuffData.ActivatedAt != 20_000 {
		t.Fatalf("expected ActivatedAt=20000")
	}

	// advance time past expiration and delete
	SetTestNow(t, 20_200)
	deleted := base.DeleteExpiredBuffs()
	if deleted != 1 {
		t.Fatalf("expected 1 expired buff deleted, got %d", deleted)
	}
}
