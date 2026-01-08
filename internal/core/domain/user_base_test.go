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
	base.Stats.SpaceCapacity = 1

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
	// After queueing but before starting production, pending armies should reserve space.
	if base.Stats.Space != army.Space {
		t.Fatalf("expected space=%d after queueing army, got %d", army.Space, base.Stats.Space)
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

func TestArmy_QueueArmy_NotEnoughSpace(t *testing.T) {
	SetTestNow(t, 5_100)
	base := newBaseWithDefaults(3)
	// artificially restrict space capacity to simulate a nearly full base
	base.Stats.SpaceCapacity = 1

	// Unlock infantry via present military building
	base.BuildingsPresent = []BuildItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype: BuildItemPrototype{
			ID:       12,
			Name:     "Barracks",
			Category: BuildCategoryMilitary,
			MilitaryData: &MilitaryBuildingData{
				UnlockArmyCategory: ArmyCategoryInfantry,
			},
		},
	}}

	army := &ArmyItemPrototype{
		ID:             105,
		Name:           "Heavy Infantry",
		Category:       ArmyCategoryInfantry,
		Price:          PriceModel{Credits: 100},
		ProductionTime: 60,
		Space:          2, // exceeds capacity
	}

	if err := base.QueueArmy(army, 1); err == nil {
		t.Fatalf("expected error when queueing army without enough space")
	}
	// no items should be queued or in production and resources must be unchanged
	if len(base.ArmiesPending) != 0 || len(base.ArmiesInProduction) != 0 {
		t.Fatalf("expected no armies queued or in production after space error")
	}
	if base.Stats.Credits != 10_000 {
		t.Fatalf("expected credits to remain unchanged after army space error, got %+v", base.Stats)
	}
	if events := base.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no events after failed QueueArmy due to space, got %v", events)
	}
}

func TestArmy_QueueArmy_NotAvailableWithoutBuilding(t *testing.T) {
	SetTestNow(t, 5_500)
	base := newBaseWithDefaults(20)

	army := &ArmyItemPrototype{
		ID:             101,
		Name:           "Infantry",
		Category:       ArmyCategoryInfantry,
		Price:          PriceModel{Credits: 50},
		ProductionTime: 60,
		Space:          1,
	}

	// No military building present to unlock infantry -> not available
	if err := base.QueueArmy(army, 1); err == nil {
		t.Fatalf("expected error when queueing army without required military building")
	}
	if len(base.ArmiesPending) != 0 || len(base.ArmiesInProduction) != 0 {
		t.Fatalf("expected no pending or in-production armies after failed QueueArmy")
	}
	if base.Stats.Credits != 10_000 {
		t.Fatalf("expected credits to remain unchanged after failed QueueArmy, got %+v", base.Stats)
	}
	if events := base.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no events after failed QueueArmy, got %v", events)
	}
}

func TestArmy_CancelPendingArmyByID_RefundsAndEmitsEvent(t *testing.T) {
	SetTestNow(t, 6_000)
	base := newBaseWithDefaults(21)

	// Simulate a pending batch of 5 infantry that has already been paid for.
	proto := ArmyItemPrototype{
		ID:             102,
		Name:           "Infantry",
		Category:       ArmyCategoryInfantry,
		Price:          PriceModel{Credits: 10},
		ProductionTime: 60,
		Space:          1,
	}
	// Mimic post-QueueArmy debit for count=5
	base.Stats.Credits -= proto.Price.Credits * 5
	base.ArmiesPending = []ArmyItemPending{{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype:     proto,
		Count:         5,
	}}
	pendingID := base.ArmiesPending[0].ID

	// Cancel 2 out of 5 -> expect refund for 2 and remaining count=3
	if err := base.CancelPendingArmyByID(pendingID, 2); err != nil {
		t.Fatalf("CancelPendingArmyByID error: %v", err)
	}
	// credits should now reflect payment for 3 units (10_000 - 10*3)
	if base.Stats.Credits != 10_000-10*3 {
		t.Fatalf("unexpected credits after partial cancel: %+v", base.Stats)
	}
	if len(base.ArmiesPending) != 1 || base.ArmiesPending[0].Count != 3 {
		t.Fatalf("expected pending count=3 after partial cancel, got %+v", base.ArmiesPending)
	}
	// Event should be emitted
	events := base.PullEvents()
	foundCancelled := false
	for _, e := range events {
		if _, ok := e.(ArmyProductionCancelledEvent); ok {
			foundCancelled = true
			break
		}
	}
	if !foundCancelled {
		t.Fatalf("expected ArmyProductionCancelledEvent after cancel")
	}

	// Now cancel remaining 3 -> entry removed and credits fully restored
	base.PullEvents() // clear
	if err := base.CancelPendingArmyByID(pendingID, 3); err != nil {
		t.Fatalf("CancelPendingArmyByID (full) error: %v", err)
	}
	if len(base.ArmiesPending) != 0 {
		t.Fatalf("expected no pending armies after full cancel, got %+v", base.ArmiesPending)
	}
	if base.Stats.Credits != 10_000 {
		t.Fatalf("expected credits fully restored after full cancel, got %+v", base.Stats)
	}
}

func TestArmy_MoveArmyQueue_RespectsCategorySlots(t *testing.T) {
	SetTestNow(t, 7_000)
	base := newBaseWithDefaults(22)

	// One barracks -> one production slot for infantry
	base.BuildingsPresent = []BuildItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype: BuildItemPrototype{
			ID:       11,
			Name:     "Barracks",
			Category: BuildCategoryMilitary,
			MilitaryData: &MilitaryBuildingData{
				UnlockArmyCategory: ArmyCategoryInfantry,
			},
		},
	}}

	armyProto := ArmyItemPrototype{
		ID:             103,
		Name:           "Infantry",
		Category:       ArmyCategoryInfantry,
		Price:          PriceModel{},
		ProductionTime: 120,
		Space:          1,
	}

	// Simulate one active production already occupying the only slot
	base.ArmiesInProduction = []ArmyItemInProduction{{
		BaseOwnedItem:     NewBaseOwnedItem(base.ID),
		Prototype:         armyProto,
		StartDate:         NowUnix(),
		CompletionDate:    NowUnix() + 1_000, // not finished yet
		CrystalsSkipPrice: int(armyProto.ProductionTime / 60),
	}}

	// Pending batch that should NOT start because slot is already full
	base.ArmiesPending = []ArmyItemPending{{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype:     armyProto,
		Count:         3,
	}}

	base.PullEvents() // clear
	base.MoveArmyQueue()

	// Still exactly one in production; pending unchanged
	if len(base.ArmiesInProduction) != 1 {
		t.Fatalf("expected 1 army in production with full slots, got %d", len(base.ArmiesInProduction))
	}
	if len(base.ArmiesPending) != 1 || base.ArmiesPending[0].Count != 3 {
		t.Fatalf("expected pending count to remain 3 when slots are full, got %+v", base.ArmiesPending)
	}
	// No new ArmyProductionStartedEvent should be emitted
	for _, e := range base.PullEvents() {
		if _, ok := e.(ArmyProductionStartedEvent); ok {
			t.Fatalf("did not expect ArmyProductionStartedEvent when all slots are full")
		}
	}
}

func TestArmy_GetReadyToDeployArmy_SuccessDoesNotMutate(t *testing.T) {
	base := newBaseWithDefaults(30)

	// Two present stacks
	base.ArmiesPresent = []ArmyItemPresent{
		{
			BaseOwnedItem: NewBaseOwnedItem(base.ID),
			Prototype:     ArmyItemPrototype{ID: 200, Category: ArmyCategoryInfantry, Attack: 1, Defence: 1, Space: 1},
			Count:         5,
		},
		{
			BaseOwnedItem: NewBaseOwnedItem(base.ID),
			Prototype:     ArmyItemPrototype{ID: 201, Category: ArmyCategoryInfantry, Attack: 2, Defence: 2, Space: 1},
			Count:         3,
		},
	}

	req := []ArmyDeploymentRequest{
		{PresentItemID: base.ArmiesPresent[0].ID, Count: 2},
		{PresentItemID: base.ArmiesPresent[1].ID, Count: 1},
	}

	ready, err := base.GetReadyToDeployArmy(req)
	if err != nil {
		t.Fatalf("GetReadyToDeployArmy error: %v", err)
	}
	if len(ready) != 2 {
		t.Fatalf("expected 2 ready items, got %d", len(ready))
	}
	if ready[0].PresentItemID != base.ArmiesPresent[0].ID || ready[0].Count != 2 || ready[0].Prototype.ID != 200 {
		t.Fatalf("unexpected first ready item: %+v", ready[0])
	}
	if ready[1].PresentItemID != base.ArmiesPresent[1].ID || ready[1].Count != 1 || ready[1].Prototype.ID != 201 {
		t.Fatalf("unexpected second ready item: %+v", ready[1])
	}
	// Must not mutate present inventory or deployed state
	if base.ArmiesPresent[0].Count != 5 || base.ArmiesPresent[1].Count != 3 {
		t.Fatalf("expected ArmiesPresent counts unchanged, got %+v", base.ArmiesPresent)
	}
	if len(base.ArmiesDeployed) != 0 {
		t.Fatalf("expected no ArmiesDeployed after GetReadyToDeployArmy, got %+v", base.ArmiesDeployed)
	}
}

func TestArmy_GetReadyToDeployArmy_InvalidCountErrorsAndDoesNotMutate(t *testing.T) {
	base := newBaseWithDefaults(31)
	base.ArmiesPresent = []ArmyItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype:     ArmyItemPrototype{ID: 210, Category: ArmyCategoryInfantry, Space: 1},
		Count:         4,
	}}

	req := []ArmyDeploymentRequest{{PresentItemID: base.ArmiesPresent[0].ID, Count: 5}} // > available
	ready, err := base.GetReadyToDeployArmy(req)
	if err == nil {
		t.Fatalf("expected error for invalid deployment count, got ready=%+v", ready)
	}
	if len(ready) != 0 {
		t.Fatalf("expected no ready items on error, got %+v", ready)
	}
	if base.ArmiesPresent[0].Count != 4 {
		t.Fatalf("expected ArmiesPresent count unchanged on error, got %+v", base.ArmiesPresent[0])
	}
}

func TestArmy_GetReadyToDeployArmy_EmptyRequestsErrors(t *testing.T) {
	base := newBaseWithDefaults(33)
	base.ArmiesPresent = []ArmyItemPresent{ // some present units
		{
			BaseOwnedItem: NewBaseOwnedItem(base.ID),
			Prototype:     ArmyItemPrototype{ID: 230, Category: ArmyCategoryInfantry, Space: 1},
			Count:         3,
		},
	}

	ready, err := base.GetReadyToDeployArmy(nil)
	if err == nil {
		t.Fatalf("expected error for empty deployment requests, got ready=%+v", ready)
	}
	if len(ready) != 0 {
		t.Fatalf("expected no ready items on error, got %+v", ready)
	}
	if base.ArmiesPresent[0].Count != 3 {
		t.Fatalf("expected ArmiesPresent count unchanged on error, got %+v", base.ArmiesPresent[0])
	}
}

func TestArmy_AllocateArmyToOperation_RemovesFromPresentAndMergesDeployed(t *testing.T) {
	base := newBaseWithDefaults(32)
	proto := ArmyItemPrototype{ID: 220, Category: ArmyCategoryInfantry, Space: 1}
	base.ArmiesPresent = []ArmyItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype:     proto,
		Count:         5,
	}}
	presentID := base.ArmiesPresent[0].ID
	opID := 99

	// First allocation of 2 units
	chunk, err := base.AllocateArmyToOperation(ArmyDeploymentRequest{PresentItemID: presentID, Count: 2}, opID)
	if err != nil {
		t.Fatalf("AllocateArmyToOperation (1) error: %v", err)
	}
	if chunk.Count != 2 || chunk.Prototype.ID != proto.ID || chunk.OperationID != opID {
		t.Fatalf("unexpected deployed chunk from first allocation: %+v", chunk)
	}
	if base.ArmiesPresent[0].Count != 3 {
		t.Fatalf("expected present count=3 after first allocation, got %+v", base.ArmiesPresent[0])
	}
	if len(base.ArmiesDeployed) != 1 || base.ArmiesDeployed[0].Count != 2 {
		t.Fatalf("expected 1 deployed entry with count=2, got %+v", base.ArmiesDeployed)
	}

	// Second allocation of 1 unit should merge into same deployed entry
	chunk, err = base.AllocateArmyToOperation(ArmyDeploymentRequest{PresentItemID: presentID, Count: 1}, opID)
	if err != nil {
		t.Fatalf("AllocateArmyToOperation (2) error: %v", err)
	}
	if base.ArmiesPresent[0].Count != 2 {
		t.Fatalf("expected present count=2 after second allocation, got %+v", base.ArmiesPresent[0])
	}
	if len(base.ArmiesDeployed) != 1 || base.ArmiesDeployed[0].Count != 3 {
		t.Fatalf("expected merged deployed count=3 for same operation/prototype, got %+v", base.ArmiesDeployed)
	}
}

func TestArmy_ReturnAllDeployedFromOperation_MergesBackAndCleans(t *testing.T) {
	base := newBaseWithDefaults(33)
	proto1 := ArmyItemPrototype{ID: 230, Category: ArmyCategoryInfantry}
	proto2 := ArmyItemPrototype{ID: 231, Category: ArmyCategoryInfantry}

	// Present already has some proto1 and an unrelated proto3
	base.ArmiesPresent = []ArmyItemPresent{
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto1, Count: 5},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: ArmyItemPrototype{ID: 232}, Count: 1},
	}

	// Deployed for two operations
	base.ArmiesDeployed = []ArmyItemDeployed{
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto1, OperationID: 10, Count: 2},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto2, OperationID: 10, Count: 3},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto1, OperationID: 11, Count: 4},
	}

	base.ReturnAllDeployedFromOperation(10)

	// ArmiesDeployed should retain only opID 11
	if len(base.ArmiesDeployed) != 1 || base.ArmiesDeployed[0].OperationID != 11 || base.ArmiesDeployed[0].Prototype.ID != proto1.ID {
		t.Fatalf("expected only opID 11 deployed to remain, got %+v", base.ArmiesDeployed)
	}

	// Present counts: proto1 increased by 2, proto2 added as new stack, proto3 unchanged
	var count1, count2, count3 int
	for _, ap := range base.ArmiesPresent {
		switch ap.Prototype.ID {
		case 230:
			count1 = ap.Count
		case 231:
			count2 = ap.Count
		case 232:
			count3 = ap.Count
		}
	}
	if count1 != 7 || count2 != 3 || count3 != 1 {
		t.Fatalf("unexpected present counts after ReturnAllDeployedFromOperation: proto1=%d proto2=%d proto3=%d", count1, count2, count3)
	}
}

func TestArmy_TrimDeployedToSurvivors_AdjustsCountsAndKeepsZeroEntries(t *testing.T) {
	base := newBaseWithDefaults(34)
	proto1 := ArmyItemPrototype{ID: 240}
	proto2 := ArmyItemPrototype{ID: 241}

	base.ArmiesDeployed = []ArmyItemDeployed{
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto1, OperationID: 50, Count: 5},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto2, OperationID: 50, Count: 4},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto1, OperationID: 51, Count: 7}, // different op
	}

	survivors := []MilitaryUnit{{PrototypeID: proto1.ID, Count: 3}}
	base.TrimDeployedToSurvivors(50, survivors)

	var op50p1, op50p2, op51p1 int
	for _, d := range base.ArmiesDeployed {
		if d.OperationID == 50 && d.Prototype.ID == proto1.ID {
			op50p1 = d.Count
		}
		if d.OperationID == 50 && d.Prototype.ID == proto2.ID {
			op50p2 = d.Count
		}
		if d.OperationID == 51 && d.Prototype.ID == proto1.ID {
			op51p1 = d.Count
		}
	}
	if op50p1 != 3 || op50p2 != 0 || op51p1 != 7 {
		t.Fatalf("unexpected deployed counts after TrimDeployedToSurvivors: op50p1=%d op50p2=%d op51p1=%d", op50p1, op50p2, op51p1)
	}
}

func TestTech_StartAndSpeedUp_EmitsEvents(t *testing.T) {
	SetTestNow(t, 10_000)
	base := newBaseWithDefaults(3)
	tech := &TechItemPrototype{
		ID:           200,
		Name:         "Improved Mining",
		Category:     TechCategoryBuild,
		Price:        PriceModel{Credits: 250},
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
	// resources should be debited by tech price
	if base.Stats.Credits != 10_000-250 {
		t.Fatalf("unexpected credits after StartTechResearch: %+v", base.Stats)
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
	if len(base.TechnologiesInProgress) != 0 {
		t.Fatalf("expected no technologies in progress after speedup")
	}
	if len(base.TechnologiesDone) != 1 || base.TechnologiesDone[0].Prototype.ID != tech.ID {
		t.Fatalf("expected tech to be moved to done after speedup, got %+v", base.TechnologiesDone)
	}
}

func TestTech_StartTechResearch_NotAvailableWhenAlreadyDone(t *testing.T) {
	SetTestNow(t, 11_000)
	base := newBaseWithDefaults(5)
	tech := &TechItemPrototype{
		ID:           201,
		Name:         "Shielding",
		Category:     TechCategoryBuild,
		Price:        PriceModel{Credits: 100},
		ResearchTime: 30,
	}
	// mark tech as already done so AvailableTechnologies returns empty
	base.TechnologiesDone = []TechItemDone{{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype:     *tech,
	}}

	if err := base.StartTechResearch(tech); err == nil {
		t.Fatalf("expected error when starting research for an already done tech")
	}
	if len(base.TechnologiesInProgress) != 0 {
		t.Fatalf("expected no technologies in progress when start fails")
	}
	if events := base.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no events when StartTechResearch fails, got %v", events)
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

func TestStorage_ActivateBuffTwice_ErrorsAndDoesNotDuplicate(t *testing.T) {
	SetTestNow(t, 21_000)
	base := newBaseWithDefaults(6)
	buff := StorageItemPresent{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype: StorageItemPrototype{
			ID:       400,
			Name:     "Space Booster",
			Category: StorageCategoryBuff,
			BuffData: &BuffStorageData{DurationSeconds: 50},
		},
	}
	base.StorageItemsPresent = []StorageItemPresent{buff}

	// first activation succeeds
	if err := base.ActivateBuffByID(buff.ID); err != nil {
		t.Fatalf("first ActivateBuffByID error: %v", err)
	}
	firstActivatedAt := base.StorageItemsPresent[0].Prototype.BuffData.ActivatedAt
	if firstActivatedAt == nil || *firstActivatedAt != 21_000 {
		t.Fatalf("expected ActivatedAt to be set on first activation, got %+v", firstActivatedAt)
	}
	base.PullEvents() // clear

	// second activation should return error and not change ActivatedAt or emit events
	if err := base.ActivateBuffByID(buff.ID); err == nil {
		t.Fatalf("expected error on second ActivateBuffByID for same buff")
	}
	secondActivatedAt := base.StorageItemsPresent[0].Prototype.BuffData.ActivatedAt
	if secondActivatedAt == nil || *secondActivatedAt != *firstActivatedAt {
		t.Fatalf("expected ActivatedAt to remain unchanged on second activation, got %+v", secondActivatedAt)
	}
	if events := base.PullEvents(); len(events) != 0 {
		t.Fatalf("expected no additional events on second activation, got %v", events)
	}
}

func TestApplyDefenderArmyRemaining_UpdatesCountsAndZeroesMissing(t *testing.T) {
	base := newBaseWithDefaults(40)
	base.ArmiesPresent = []ArmyItemPresent{
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: ArmyItemPrototype{ID: 300}, Count: 5},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: ArmyItemPrototype{ID: 301}, Count: 4},
	}

	remaining := []MilitaryUnit{
		{PrototypeID: 300, Count: 2}, // proto 300 survives with 2
	}
	base.ApplyDefenderArmyRemaining(remaining)

	var count300, count301 int
	for _, ap := range base.ArmiesPresent {
		switch ap.Prototype.ID {
		case 300:
			count300 = ap.Count
		case 301:
			count301 = ap.Count
		}
	}
	if count300 != 2 || count301 != 0 {
		t.Fatalf("unexpected counts after ApplyDefenderArmyRemaining: proto300=%d proto301=%d", count300, count301)
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
	remaining := []DefenseStructure{{PrototypeID: 400, Count: 2}}
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
