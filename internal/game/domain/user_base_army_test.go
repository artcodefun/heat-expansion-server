package domain

import (
	"strings"
	"testing"
)

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
			Faction:      FactionExoCoalition,
			MilitaryData: &MilitaryBuildingData{UnlockArmyCategory: ArmyCategoryInfantry},
			Space:        0,
		},
	}}

	army := &ArmyItemPrototype{
		ID:             100,
		Name:           "Infantry",
		Category:       ArmyCategoryInfantry,
		Faction:        FactionExoCoalition,
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
	// by adding a building that consumes most of the default capacity (100)

	// Unlock infantry via present military building
	base.BuildingsPresent = []BuildItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(base.ID),
		Prototype: BuildItemPrototype{
			ID:       12,
			Name:     "Barracks",
			Category: BuildCategoryMilitary,
			Faction:  FactionExoCoalition,
			MilitaryData: &MilitaryBuildingData{
				UnlockArmyCategory: ArmyCategoryInfantry,
			},
			Space: 99,
		},
	}}

	army := &ArmyItemPrototype{
		ID:             105,
		Name:           "Heavy Infantry",
		Category:       ArmyCategoryInfantry,
		Faction:        FactionExoCoalition,
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
	if base.Stats.Credits != 1000 {
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
		Faction:        FactionExoCoalition,
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
	if base.Stats.Credits != 1000 {
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
		Faction:        FactionExoCoalition,
		Price:          PriceModel{Credits: 10},
		ProductionTime: 60,
		Space:          1,
	}
	// Mimic post-QueueArmy debit for count=5
	base.Stats.Credits -= float64(proto.Price.Credits * 5)
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
	// credits should now reflect payment for 3 units (1000 - 10*3)
	if base.Stats.Credits != 1000-10*3 {
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
	if base.Stats.Credits != 1000 {
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
			Faction:  FactionExoCoalition,
			MilitaryData: &MilitaryBuildingData{
				UnlockArmyCategory: ArmyCategoryInfantry,
			},
		},
	}}

	armyProto := ArmyItemPrototype{
		ID:             103,
		Name:           "Infantry",
		Category:       ArmyCategoryInfantry,
		Faction:        FactionExoCoalition,
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
			Prototype:     ArmyItemPrototype{ID: 200, Category: ArmyCategoryInfantry, Faction: FactionExoCoalition, Attack: 1, Defence: 1, Space: 1},
			Count:         5,
		},
		{
			BaseOwnedItem: NewBaseOwnedItem(base.ID),
			Prototype:     ArmyItemPrototype{ID: 201, Category: ArmyCategoryInfantry, Faction: FactionExoCoalition, Attack: 2, Defence: 2, Space: 1},
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
		Prototype:     ArmyItemPrototype{ID: 210, Category: ArmyCategoryInfantry, Faction: FactionExoCoalition, Space: 1},
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
			Prototype:     ArmyItemPrototype{ID: 230, Category: ArmyCategoryInfantry, Faction: FactionExoCoalition, Space: 1},
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
	proto := ArmyItemPrototype{ID: 220, Category: ArmyCategoryInfantry, Faction: FactionExoCoalition, Space: 1}
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

func TestArmy_AllocateArmyToOperation_RespectsMaxOperations(t *testing.T) {
	base := newBaseWithDefaults(34) // recalculateStats() sets DefaultMaxOperations (2)
	proto := ArmyItemPrototype{ID: 240, Category: ArmyCategoryInfantry, Faction: FactionExoCoalition, Space: 1}

	base.ArmiesPresent = []ArmyItemPresent{
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto, Count: 10},
	}
	presentID := base.ArmiesPresent[0].ID

	// Use up 2 operation slots
	for i := 1; i <= 2; i++ {
		if _, err := base.AllocateArmyToOperation(ArmyDeploymentRequest{PresentItemID: presentID, Count: 1}, i); err != nil {
			t.Fatalf("allocation for op %d failed: %v", i, err)
		}
	}

	// 3rd operation should fail
	if _, err := base.AllocateArmyToOperation(ArmyDeploymentRequest{PresentItemID: presentID, Count: 1}, 3); err == nil {
		t.Errorf("expected error for exceeding MaxOperations (2), got nil")
	} else if !strings.HasPrefix(err.Error(), "error.domain.operation.max_reached") {
		t.Errorf("unexpected error message: %v", err)
	}

	// Adding more units to an existing operation (say op 2) should still succeed
	if _, err := base.AllocateArmyToOperation(ArmyDeploymentRequest{PresentItemID: presentID, Count: 1}, 2); err != nil {
		t.Errorf("failed to add units to existing operation 2: %v", err)
	}
}

func TestArmy_ReturnAllDeployedFromOperation_MergesBackAndCleans(t *testing.T) {
	base := newBaseWithDefaults(33)
	proto1 := ArmyItemPrototype{ID: 230, Category: ArmyCategoryInfantry, Faction: FactionExoCoalition}
	proto2 := ArmyItemPrototype{ID: 231, Category: ArmyCategoryInfantry, Faction: FactionExoCoalition}

	// Present already has some proto1 and an unrelated proto3
	base.ArmiesPresent = []ArmyItemPresent{
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto1, Count: 5},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: ArmyItemPrototype{ID: 232, Faction: FactionExoCoalition}, Count: 1},
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
	proto1 := ArmyItemPrototype{ID: 240, Faction: FactionExoCoalition}
	proto2 := ArmyItemPrototype{ID: 241, Faction: FactionExoCoalition}

	base.ArmiesDeployed = []ArmyItemDeployed{
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto1, OperationID: 50, Count: 5},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto2, OperationID: 50, Count: 4},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: proto1, OperationID: 51, Count: 7}, // different op
	}

	survivors := []MilitaryUnitSnap{{PrototypeID: proto1.ID, Count: 3}}
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

func TestApplyDefenderArmyRemaining_UpdatesCountsAndZeroesMissing(t *testing.T) {
	base := newBaseWithDefaults(40)
	base.ArmiesPresent = []ArmyItemPresent{
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: ArmyItemPrototype{ID: 300, Faction: FactionExoCoalition}, Count: 5},
		{BaseOwnedItem: NewBaseOwnedItem(base.ID), Prototype: ArmyItemPrototype{ID: 301, Faction: FactionExoCoalition}, Count: 4},
	}

	remaining := []MilitaryUnitSnap{
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
