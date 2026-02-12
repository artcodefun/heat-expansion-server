package domain

import "testing"

func TestOperation_Attack_EmptyLocation_PhaseAndEvents(t *testing.T) {
	SetTestNow(t, 1_000)
	units := []MilitaryUnitSnap{
		{PrototypeID: 1, Category: ArmyCategoryInfantry, Attack: 10, Defence: 5, Capacity: 0, Stealth: 0, Speed: 100, Count: 5},
	}
	source := Vector2i{X: 0, Y: 0}
	target := Vector2i{X: 3, Y: 4} // Euclidean distance 5 -> scaled 5000
	op, err := NewAttackOperation(1, 10, source, target, units, nil)
	if err != nil {
		t.Fatalf("unexpected error from NewAttackOperation: %v", err)
	}
	if op.Phase != OperationPhasePending {
		t.Fatalf("expected pending phase")
	}

	op.Start()
	if op.Phase != OperationPhaseOutbound {
		t.Fatalf("expected outbound phase after Start, got %s", op.Phase)
	}
	expectedTravel := computeTravelSecondsBetween(source, target, units, MilitaryModifiersFromSnaps(nil))
	if op.OutboundArriveAt != 1_000+expectedTravel {
		t.Fatalf("unexpected OutboundArriveAt, want %d got %d", 1_000+expectedTravel, op.OutboundArriveAt)
	}
	// Started event present
	started := false
	for _, e := range op.PullEvents() {
		if _, ok := e.(MilitaryOperationStartedEvent); ok {
			started = true
		}
	}
	if !started {
		t.Fatalf("expected MilitaryOperationStartedEvent after Start")
	}

	// not yet arrived
	SetTestNow(t, 1_000+expectedTravel-1)
	op.UpdatePhaseBasedOnTime()
	if op.Phase != OperationPhaseOutbound {
		t.Fatalf("should still be outbound before arrival threshold")
	}
	// arrive
	SetTestNow(t, 1_000+expectedTravel)
	op.UpdatePhaseBasedOnTime()
	if op.Phase != OperationPhaseAtTarget {
		t.Fatalf("expected at target after arrival, got %s", op.Phase)
	}
	arrived := false
	for _, e := range op.PullEvents() {
		if _, ok := e.(MilitaryOperationArrivedEvent); ok {
			arrived = true
		}
	}
	if !arrived {
		t.Fatalf("expected MilitaryOperationArrivedEvent on arrival")
	}

	// Resolve directly against an empty location (no defenders, no loot)
	res := op.ResolveAttack(nil, nil, PriceModel{}, nil, nil)
	if res == nil {
		t.Fatalf("expected non-nil AttackResult when resolving against empty location")
	}
	if op.Phase != OperationPhaseResolving {
		t.Fatalf("expected resolving phase after ResolveAttack, got %s", op.Phase)
	}

	// Start return leg from the operation itself
	op.StartReturn()
	if op.Phase != OperationPhaseReturning {
		t.Fatalf("expected returning phase after StartReturn, got %s", op.Phase)
	}
	// Return-started event
	returnStarted := false
	for _, e := range op.PullEvents() {
		if _, ok := e.(MilitaryOperationReturnStartedEvent); ok {
			returnStarted = true
		}
	}
	if !returnStarted {
		t.Fatalf("expected MilitaryOperationReturnStartedEvent")
	}

	// complete return
	SetTestNow(t, op.ReturnArriveAt)
	op.UpdatePhaseBasedOnTime()
	if op.Phase != OperationPhaseCompleted {
		t.Fatalf("expected completed phase after return arrival, got %s", op.Phase)
	}
	returnArrived := false
	for _, e := range op.PullEvents() {
		if _, ok := e.(MilitaryOperationReturnArrivedEvent); ok {
			returnArrived = true
		}
	}
	if !returnArrived {
		t.Fatalf("expected MilitaryOperationReturnArrivedEvent")
	}
}

func TestOperation_TimeBeforeEntersCircle(t *testing.T) {
	SetTestNow(t, 1_000)
	units := []MilitaryUnitSnap{
		{PrototypeID: 1, Category: ArmyCategoryInfantry, Attack: 10, Defence: 5, Capacity: 0, Stealth: 0, Speed: 100, Count: 1},
	}
	source := Vector2i{X: 0, Y: 0}
	target := Vector2i{X: 10, Y: 0} // Distance 10 (scaled 10000)
	// Speed 100 -> Total travel time 10000 / 100 = 100s

	op, _ := NewAttackOperation(1, 1, source, target, units, nil)
	op.Start() // OutboundDepartAt = 1000, OutboundArriveAt = 1100

	center := Vector2i{X: 10, Y: 0}
	radius := 4
	// Edge of circle at X=6. Travels from X=0 to X=6.
	// That's 6/10 of total distance.
	// So 6/10 of total time (100s) = 60s.
	// Enter at 1000 + 60 = 1060.

	enterAt, err := op.TimeBeforeEntersCircle(center, radius)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if enterAt != 1060 {
		t.Errorf("expected enterAt 1060, got %d", enterAt)
	}

	// Test already inside
	centerInside := Vector2i{X: 1, Y: 0}
	radiusInside := 2
	enterAtInside, err := op.TimeBeforeEntersCircle(centerInside, radiusInside)
	if err != nil {
		t.Fatalf("unexpected error for inside center: %v", err)
	}
	if enterAtInside != 1000 {
		t.Errorf("expected enterAt 1000 for starting inside, got %d", enterAtInside)
	}

	// Test never enters
	centerFar := Vector2i{X: 0, Y: 10}
	radiusFar := 2
	_, err = op.TimeBeforeEntersCircle(centerFar, radiusFar)
	if err == nil {
		t.Errorf("expected error for never entering circle")
	}
}

func TestOperation_TotalStealth(t *testing.T) {
	op := &MilitaryOperation{
		Units: []MilitaryUnitSnap{
			{Stealth: 10, Count: 2},
			{Stealth: 5, Count: 3},
			{Stealth: 0, Count: 10},
		},
		TotalModifiers: MilitaryModifiers{
			AttackMul:   1.0,
			DefenceMul:  1.0,
			StealthMul:  1.0,
			CapacityMul: 1.0,
			SpeedMul:    1.0,
		},
	}
	expected := 2*10 + 3*5 // 35
	if got := op.TotalStealth(); got != expected {
		t.Errorf("expected total stealth %d, got %d", expected, got)
	}
}

func TestOperation_Attack_UserBase_LootAndDeduction(t *testing.T) {
	SetTestNow(t, 2_000)
	// Attacking unit with capacity 5
	units := []MilitaryUnitSnap{
		{PrototypeID: 2, Category: ArmyCategoryInfantry, Attack: 10, Defence: 5, Capacity: 7, Stealth: 0, Speed: 200, Count: 1},
	}
	op, err := NewAttackOperation(1, 10, Vector2i{0, 0}, Vector2i{1, 0}, units, nil)
	if err != nil {
		t.Fatalf("unexpected error from NewAttackOperation: %v", err)
	}
	op.Start()
	// arrive immediately for simplicity
	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()

	// Defender with some resources
	def := &UserBaseModel{ID: 99}
	def.Stats = UserBaseStats{Credits: 10, Iron: 8, Titanium: 4, Antimatter: 2, CalculationTimestamp: NowUnix()}
	// No defenders/structures for a guaranteed win and no losses

	attackerBase := &UserBaseModel{ID: 10}
	service := NewMilitaryOperationService(op, attackerBase)
	service.ResolveAgainstUserBase(def)

	if op.AttackResult == nil {
		t.Fatalf("expected AttackResult to be set")
	}
	loot := op.AttackResult.Loot
	lootSum := loot.Credits + loot.Iron + loot.Titanium + loot.Antimatter
	// With Capacity 7 and WorthCapacityMultiplier=10, we have 70 volume points.
	// Greedy looting (Low to High value):
	// - Credits (Worth 1): takes 10 (vol used: 10, rem: 60)
	// - Iron (Worth 4): takes 8 (vol used: 32, rem: 28)
	// - Titanium (Worth 20): fits 1
	// - Antimatter (Worth 333.3): fits 0
	// Total items = 10 + 8 + 1 = 19.
	if lootSum != 19 {
		t.Fatalf("unexpected loot sum: %d (expected 19 due to least-expensive-first looting)", lootSum)
	}
	// Defender resources reduced by exactly loot sum
	afterSum := int(def.Stats.Credits + def.Stats.Iron + def.Stats.Titanium + def.Stats.Antimatter)
	if before := 10 + 8 + 4 + 2; before-afterSum != lootSum {
		t.Fatalf("expected defender resource sum decrease by %d, got %d", lootSum, before-afterSum)
	}
	// Should have started return
	if op.Phase != OperationPhaseReturning && op.Phase != OperationPhaseCompleted {
		t.Fatalf("expected returning or completed phase, got %s", op.Phase)
	}
}

func TestSpy_BlockedByCloaking_OutcomeAndReturn(t *testing.T) {
	SetTestNow(t, 3_000)
	// Attacker spies
	spies := []MilitaryUnitSnap{
		{PrototypeID: 7, Category: ArmyCategorySpy, Attack: 2, Defence: 1, Capacity: 0, Stealth: 4, Speed: 120, Count: 2}, // stealth sum = 8
	}
	op, err := NewSpyOperation(1, 10, Vector2i{0, 0}, Vector2i{1, 1}, spies, nil)
	if err != nil {
		t.Fatalf("unexpected error from NewSpyOperation: %v", err)
	}
	op.Start()
	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()

	// Strong cloaking at target (>= attacker stealth) should block the spy
	res := op.ResolveSpy(10, nil, nil)
	if res == nil || res.Outcome != SpyOutcomeBlockedByCloaking {
		t.Fatalf("expected spy outcome BLOCKED_BY_CLOAKING, got %+v", op.SpyResult)
	}
	if op.Result != OperationResultSuccess {
		t.Fatalf("expected operation result SUCCESS, got %s", op.Result)
	}

	op.StartReturn()
	if op.Phase != OperationPhaseReturning {
		t.Fatalf("expected returning phase after StartReturn, got %s", op.Phase)
	}
	// Ensure return started event exists
	hasResolved, hasReturnStarted := false, false
	for _, e := range op.PullEvents() {
		switch e.(type) {
		case MilitaryOperationResolvedEvent:
			hasResolved = true
		case MilitaryOperationReturnStartedEvent:
			hasReturnStarted = true
		}
	}
	if !hasResolved || !hasReturnStarted {
		t.Fatalf("expected resolved and return-started events; got resolved=%v returnStarted=%v", hasResolved, hasReturnStarted)
	}
}

func TestSpy_DefeatedByDefendingSpies_ReturnImmediate(t *testing.T) {
	SetTestNow(t, 3_500)
	// Attacker spies with low attack
	spies := []MilitaryUnitSnap{
		{PrototypeID: 8, Category: ArmyCategorySpy, Attack: 1, Defence: 1, Capacity: 0, Stealth: 1, Speed: 100, Count: 2}, // atk power=2
	}
	op, err := NewSpyOperation(1, 10, Vector2i{0, 0}, Vector2i{1, 0}, spies, nil)
	if err != nil {
		t.Fatalf("unexpected error from NewSpyOperation: %v", err)
	}
	op.Start()
	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()

	// Defender spies with higher defence -> attackers lose
	defendingSpies := []MilitaryUnitSnap{{PrototypeID: 600, Category: ArmyCategorySpy, Attack: 1, Defence: 5, Count: 1}}
	res := op.ResolveSpy(0, defendingSpies, nil)
	if res == nil || res.Outcome != SpyOutcomeDefeatedBySpies {
		t.Fatalf("expected spy outcome DEFEATED_BY_DEFENDING_SPIES, got %+v", op.SpyResult)
	}
	if op.Result != OperationResultFailure {
		t.Fatalf("expected operation result FAILURE, got %s", op.Result)
	}

	// No survivors -> StartReturn should immediately complete the operation
	op.StartReturn()
	if op.Phase != OperationPhaseCompleted {
		t.Fatalf("expected completed phase due to zero survivors, got %s", op.Phase)
	}
	// Ensure resolved + return-arrived events exist
	hasResolved, hasReturnArrived := false, false
	for _, e := range op.PullEvents() {
		switch e.(type) {
		case MilitaryOperationResolvedEvent:
			hasResolved = true
		case MilitaryOperationReturnArrivedEvent:
			hasReturnArrived = true
		}
	}
	if !hasResolved || !hasReturnArrived {
		t.Fatalf("expected resolved and return-arrived events; got resolved=%v returnArrived=%v", hasResolved, hasReturnArrived)
	}
}

func TestSpy_ReportProduced_OutcomeAndReturn(t *testing.T) {
	SetTestNow(t, 4_000)
	// Attacker spies strong enough to win
	spies := []MilitaryUnitSnap{
		{PrototypeID: 9, Category: ArmyCategorySpy, Attack: 10, Defence: 1, Capacity: 0, Stealth: 2, Speed: 100, Count: 1},
	}
	op, err := NewSpyOperation(1, 10, Vector2i{0, 0}, Vector2i{2, 0}, spies, nil)
	if err != nil {
		t.Fatalf("unexpected error from NewSpyOperation: %v", err)
	}

	// existing test continues
	op.Start()
	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()

	// Defender spies weak -> attackers win and produce a report
	defendingSpies := []MilitaryUnitSnap{{PrototypeID: 700, Category: ArmyCategorySpy, Attack: 1, Defence: 2, Count: 1}}
	res := op.ResolveSpy(0, defendingSpies, nil)
	if res == nil || res.Outcome != SpyOutcomeReportProduced {
		t.Fatalf("expected spy outcome REPORT_PRODUCED, got %+v", op.SpyResult)
	}
	if op.Result != OperationResultSuccess {
		t.Fatalf("expected operation result SUCCESS, got %s", op.Result)
	}

	op.StartReturn()
	if op.Phase != OperationPhaseReturning {
		t.Fatalf("expected returning phase, got %s", op.Phase)
	}
	// Ensure resolved + return-started events exist
	hasResolved, hasReturnStarted := false, false
	for _, e := range op.PullEvents() {
		switch e.(type) {
		case MilitaryOperationResolvedEvent:
			hasResolved = true
		case MilitaryOperationReturnStartedEvent:
			hasReturnStarted = true
		}
	}
	if !hasResolved || !hasReturnStarted {
		t.Fatalf("expected resolved and return-started events; got resolved=%v returnStarted=%v", hasResolved, hasReturnStarted)
	}
}

func TestNewAttackOperation_RejectsNoUnits(t *testing.T) {
	_, err := NewAttackOperation(1, 10, Vector2i{0, 0}, Vector2i{1, 0}, nil, nil)
	if err == nil {
		t.Fatalf("expected error when creating attack operation with no units")
	}
}

func TestNewAttackOperation_RejectsSameCoordinates(t *testing.T) {
	units := []MilitaryUnitSnap{{PrototypeID: 1, Category: ArmyCategoryInfantry, Count: 1}}
	_, err := NewAttackOperation(1, 10, Vector2i{5, 5}, Vector2i{5, 5}, units, nil)
	if err == nil {
		t.Fatalf("expected error when creating attack operation with same source/target")
	}
}

func TestNewSpyOperation_RejectsNoUnits(t *testing.T) {
	_, err := NewSpyOperation(1, 10, Vector2i{0, 0}, Vector2i{1, 0}, nil, nil)
	if err == nil {
		t.Fatalf("expected error when creating spy operation with no units")
	}
}

func TestNewSpyOperation_RejectsSameCoordinates(t *testing.T) {
	units := []MilitaryUnitSnap{{PrototypeID: 7, Category: ArmyCategorySpy, Count: 1}}
	_, err := NewSpyOperation(1, 10, Vector2i{5, 5}, Vector2i{5, 5}, units, nil)
	if err == nil {
		t.Fatalf("expected error when creating spy operation with same source/target")
	}
}

func TestNewSpyOperation_RejectsNonSpyUnits(t *testing.T) {
	units := []MilitaryUnitSnap{{PrototypeID: 1, Category: ArmyCategoryInfantry, Count: 1}}
	_, err := NewSpyOperation(1, 10, Vector2i{0, 0}, Vector2i{1, 0}, units, nil)
	if err == nil {
		t.Fatalf("expected error when creating spy operation with non-spy units")
	}
}

func TestOperation_WithStorageSnaps_AffectsResolution(t *testing.T) {
	SetTestNow(t, 5_000)

	// 1. Attacker has a 20% Attack buff
	atkSnaps := []StorageItemSnap{
		{
			Buff: &BuffStorageData{Type: BuffTypeAttackIncrease, Value: 1.2},
		},
	}
	atkUnits := []MilitaryUnitSnap{
		{PrototypeID: 1, Category: ArmyCategoryInfantry, Attack: 100, Defence: 50, Count: 10},
	}
	// Base attack: 100 * 10 = 1000. With 1.2 multiplier = 1200.

	op, err := NewAttackOperation(1, 10, Vector2i{0, 0}, Vector2i{1, 0}, atkUnits, atkSnaps)
	if err != nil {
		t.Fatalf("failed to create operation: %v", err)
	}

	op.Start()
	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()

	if got := op.TotalAttack(); got != 1200 {
		t.Errorf("expected boosted attack 1200, got %d", got)
	}

	// 2. Defender has a 50% Defence artifact
	defSnaps := []StorageItemSnap{
		{
			Artifact: &ArtifactStorageData{Type: ArtifactEffectTypeDefenceIncrease, Value: 1.5},
		},
	}
	defUnits := []MilitaryUnitSnap{
		{PrototypeID: 2, Category: ArmyCategoryInfantry, Attack: 50, Defence: 100, Count: 10},
	}
	// Base defence: 100 * 10 = 1000. With 1.5 multiplier = 1500.

	// 3. Resolve
	// Attacker power (1200) < Defender power (1500) -> Defender should hold
	res := op.ResolveAttack(defUnits, nil, PriceModel{}, nil, defSnaps)

	if res.Outcome != AttackOutcomeDefenderHeld {
		t.Errorf("expected defender to hold due to artifact bonus, but outcome was %s", res.Outcome)
	}

	// Double check that the snaps are stored/used correctly
	if len(op.StorageSnaps) != 1 {
		t.Fatalf("expected 1 attacker storage snap in operation, got %d", len(op.StorageSnaps))
	}
	if op.StorageSnaps[0].Buff == nil {
		t.Fatalf("expected attacker storage snap to have a buff")
	}
	if op.StorageSnaps[0].Buff.Value != 1.2 {
		t.Errorf("expected buff value 1.2, got %v", op.StorageSnaps[0].Buff.Value)
	}

	if len(res.DefenderStorageSnaps) != 1 {
		t.Fatalf("expected 1 defender storage snap, got %d", len(res.DefenderStorageSnaps))
	}
	if res.DefenderStorageSnaps[0].Artifact == nil {
		t.Fatalf("expected defender storage snap to have an artifact")
	}
	if res.DefenderStorageSnaps[0].Artifact.Value != 1.5 {
		t.Errorf("expected artifact value 1.5, got %v", res.DefenderStorageSnaps[0].Artifact.Value)
	}
}
