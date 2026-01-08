package domain

import "testing"

func TestCrystalSpendingService_SpeedUpBuildingProduction_DeductsCrystalsAndSpeedsUp(t *testing.T) {
	SetTestNow(t, 1_000)

	user := &User{ID: 1, Crystals: 10}
	base := newBaseWithDefaults(1)

	// Prepare a building in production with known timings and skip price
	proto := BuildItemPrototype{
		ID:             1,
		Name:           "Mine",
		Category:       BuildCategoryResources,
		ProductionTime: 500,
	}
	item := BuildItemInProduction{
		BaseOwnedItem:     NewBaseOwnedItem(base.ID),
		Prototype:         proto,
		StartDate:         900,
		CompletionDate:    1_200, // total duration 300s, remaining at t=1000 is 200s -> fraction 2/3
		CrystalsSkipPrice: 9,
	}
	base.BuildingsInProduction = []BuildItemInProduction{item}
	buildingItemID := item.BaseOwnedItem.ID

	service := NewCrystalSpendingService()
	if err := service.SpeedUpBuildingProduction(user, base, buildingItemID); err != nil {
		t.Fatalf("SpeedUpBuildingProduction error: %v", err)
	}

	// fraction = remaining/total = 200/300 = 0.666..., price=9 -> crystals ~= 6
	expectedCrystalsSpent := 6
	if user.Crystals != 10-expectedCrystalsSpent {
		t.Fatalf("unexpected user crystals after speedup: got %d, want %d", user.Crystals, 10-expectedCrystalsSpent)
	}

	// Base-level speedup should have completed the building and removed it from production
	if len(base.BuildingsInProduction) != 0 {
		t.Fatalf("expected no buildings in production after crystal speedup, got %d", len(base.BuildingsInProduction))
	}
	if len(base.BuildingsPresent) != 1 || base.BuildingsPresent[0].Prototype.ID != proto.ID {
		t.Fatalf("expected completed building present after speedup, got %+v", base.BuildingsPresent)
	}
}

func TestCrystalSpendingService_SpeedUpOperation_Outbound(t *testing.T) {
	SetTestNow(t, 1_000)

	user := &User{ID: 1, Crystals: 100}
	units := []MilitaryUnitSnap{
		{PrototypeID: 1, Category: ArmyCategoryInfantry, Attack: 10, Defence: 5, Capacity: 0, Stealth: 0, Speed: 100, Count: 5},
	}
	source := Vector2i{X: 0, Y: 0}
	target := Vector2i{X: 3, Y: 4}
	op, err := NewAttackOperation(1, 10, source, target, units)
	if err != nil {
		t.Fatalf("unexpected error from NewAttackOperation: %v", err)
	}

	op.Start()
	if op.Phase != OperationPhaseOutbound {
		t.Fatalf("expected outbound phase after Start, got %s", op.Phase)
	}

	total := op.OutboundArriveAt - op.OutboundDepartAt
	if total <= 0 {
		t.Fatalf("expected positive outbound travel time, got %d", total)
	}

	// Move time to the middle of the outbound leg
	mid := op.OutboundDepartAt + total/2
	SetTestNow(t, mid)

	service := NewCrystalSpendingService()
	beforeCrystals := user.Crystals
	if err := service.SpeedUpOperation(user, op); err != nil {
		t.Fatalf("SpeedUpOperation (outbound) error: %v", err)
	}

	// Crystals should have decreased by at least 1
	if user.Crystals >= beforeCrystals {
		t.Fatalf("expected crystals to decrease after outbound speedup, before=%d after=%d", beforeCrystals, user.Crystals)
	}

	// Operation should now be at target (arrival handler called)
	if op.Phase != OperationPhaseAtTarget {
		t.Fatalf("expected phase AT_TARGET after outbound speedup, got %s", op.Phase)
	}
}

func TestCrystalSpendingService_SpeedUpOperation_Returning(t *testing.T) {
	SetTestNow(t, 2_000)

	user := &User{ID: 1, Crystals: 100}
	units := []MilitaryUnitSnap{
		{PrototypeID: 2, Category: ArmyCategoryInfantry, Attack: 10, Defence: 5, Capacity: 0, Stealth: 0, Speed: 100, Count: 5},
	}
	source := Vector2i{X: 0, Y: 0}
	target := Vector2i{X: 4, Y: 3}
	op, err := NewAttackOperation(1, 10, source, target, units)
	if err != nil {
		t.Fatalf("unexpected error from NewAttackOperation: %v", err)
	}

	op.Start()
	// Jump to arrival and mark at target
	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()
	if op.Phase != OperationPhaseAtTarget {
		t.Fatalf("expected AT_TARGET phase after arrival, got %s", op.Phase)
	}

	// Start return leg (no combat details needed for this test)
	SetTestNow(t, op.OutboundArriveAt)
	op.StartReturn()
	if op.Phase != OperationPhaseReturning {
		t.Fatalf("expected RETURNING phase after StartReturn, got %s", op.Phase)
	}

	total := op.ReturnArriveAt - op.ReturnDepartAt
	if total <= 0 {
		t.Fatalf("expected positive return travel time, got %d", total)
	}

	// Move time to the middle of the return leg
	mid := op.ReturnDepartAt + total/2
	SetTestNow(t, mid)

	service := NewCrystalSpendingService()
	beforeCrystals := user.Crystals
	if err := service.SpeedUpOperation(user, op); err != nil {
		t.Fatalf("SpeedUpOperation (returning) error: %v", err)
	}

	if user.Crystals >= beforeCrystals {
		t.Fatalf("expected crystals to decrease after return speedup, before=%d after=%d", beforeCrystals, user.Crystals)
	}

	// Operation should now be completed (return arrival handler called)
	if op.Phase != OperationPhaseCompleted {
		t.Fatalf("expected COMPLETED phase after return speedup, got %s", op.Phase)
	}
}
