package domain

import "testing"

func TestOperation_Attack_EmptyLocation_PhaseAndEvents(t *testing.T) {
	SetTestNow(t, 1_000)
	units := []MilitaryUnit{
		{PrototypeID: 1, Category: ArmyCategoryInfantry, Attack: 10, Defence: 5, Capacity: 0, Stealth: 0, Speed: 100, Count: 5},
	}
	source := Vector2i{X: 0, Y: 0}
	target := Vector2i{X: 3, Y: 4} // Euclidean distance 5 -> scaled 5000
	op := NewAttackOperation(1, 10, source, target, units)
	if op.Phase != OperationPhasePending {
		t.Fatalf("expected pending phase")
	}

	op.Start()
	if op.Phase != OperationPhaseOutbound {
		t.Fatalf("expected outbound phase after Start, got %s", op.Phase)
	}
	expectedTravel := computeTravelSecondsBetween(source, target, units)
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

	// resolve against empty and start return
	attacker := &UserBaseModel{ID: 10}
	service := NewMilitaryOperationService(op, attacker)
	service.ResolveAgainstEmptySector(&SectorModel{})
	if op.Phase != OperationPhaseReturning {
		t.Fatalf("expected returning phase after resolve, got %s", op.Phase)
	}
	// Return events
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

func TestOperation_Attack_UserBase_LootAndDeduction(t *testing.T) {
	SetTestNow(t, 2_000)
	// Attacking unit with capacity 5
	units := []MilitaryUnit{
		{PrototypeID: 2, Category: ArmyCategoryInfantry, Attack: 10, Defence: 5, Capacity: 5, Stealth: 0, Speed: 200, Count: 1},
	}
	op := NewAttackOperation(1, 10, Vector2i{0, 0}, Vector2i{1, 0}, units)
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
	if lootSum <= 0 || lootSum > 5 {
		t.Fatalf("unexpected loot sum: %d", lootSum)
	}
	// Defender resources reduced by exactly loot sum
	afterSum := def.Stats.Credits + def.Stats.Iron + def.Stats.Titanium + def.Stats.Antimatter
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
	spies := []MilitaryUnit{
		{PrototypeID: 7, Category: ArmyCategorySpy, Attack: 2, Defence: 1, Capacity: 0, Stealth: 4, Speed: 120, Count: 2}, // stealth sum = 8
	}
	op := NewSpyOperation(1, 10, Vector2i{0, 0}, Vector2i{1, 1}, spies)
	op.Start()
	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()

	// Defender with strong cloaking (>= attacker stealth)
	def := &UserBaseModel{ID: 99}
	def.BuildingsPresent = []BuildItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(def.ID),
		Prototype: BuildItemPrototype{
			ID:               300,
			Category:         BuildCategoryIntelligence,
			IntelligenceData: &IntelligenceBuildingData{Subtype: IntelligenceSubtypeCloaking, StealthStrength: 10},
		},
	}}
	// Add some non-spy defenders to ensure they remain unchanged
	def.ArmiesPresent = []ArmyItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(def.ID),
		Prototype:     ArmyItemPrototype{ID: 500, Category: ArmyCategoryInfantry, Attack: 5, Defence: 5},
		Count:         3,
	}}

	svc := NewMilitaryOperationService(op, &UserBaseModel{ID: 10})
	svc.ResolveAgainstUserBase(def)

	if op.SpyResult == nil || op.SpyResult.Outcome != SpyOutcomeBlockedByCloaking {
		t.Fatalf("expected spy outcome BLOCKED_BY_CLOAKING, got %+v", op.SpyResult)
	}
	if op.Result != OperationResultSuccess {
		t.Fatalf("expected operation result SUCCESS, got %s", op.Result)
	}
	if op.Phase != OperationPhaseReturning {
		t.Fatalf("expected returning phase after resolve, got %s", op.Phase)
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
	// Non-spy defenders remain unchanged (ApplyDefenderArmyRemaining only touches spies merge)
	if len(def.ArmiesPresent) != 1 || def.ArmiesPresent[0].Prototype.Category != ArmyCategoryInfantry || def.ArmiesPresent[0].Count != 3 {
		t.Fatalf("unexpected defender armies state: %+v", def.ArmiesPresent)
	}
}

func TestSpy_DefeatedByDefendingSpies_ReturnImmediate(t *testing.T) {
	SetTestNow(t, 3_500)
	// Attacker spies with low attack
	spies := []MilitaryUnit{
		{PrototypeID: 8, Category: ArmyCategorySpy, Attack: 1, Defence: 1, Capacity: 0, Stealth: 1, Speed: 100, Count: 2}, // atk power=2
	}
	op := NewSpyOperation(1, 10, Vector2i{0, 0}, Vector2i{1, 0}, spies)
	op.Start()
	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()

	// Defender spies with higher defence -> attackers lose
	def := &UserBaseModel{ID: 99}
	def.ArmiesPresent = []ArmyItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(def.ID),
		Prototype:     ArmyItemPrototype{ID: 600, Category: ArmyCategorySpy, Attack: 1, Defence: 5},
		Count:         1, // def power = 5 > atk 2
	}}

	svc := NewMilitaryOperationService(op, &UserBaseModel{ID: 10})
	svc.ResolveAgainstUserBase(def)

	if op.SpyResult == nil || op.SpyResult.Outcome != SpyOutcomeDefeatedBySpies {
		t.Fatalf("expected spy outcome DEFEATED_BY_DEFENDING_SPIES, got %+v", op.SpyResult)
	}
	if op.Result != OperationResultFailure {
		t.Fatalf("expected operation result FAILURE, got %s", op.Result)
	}
	// No survivors -> immediate completion
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
	spies := []MilitaryUnit{
		{PrototypeID: 9, Category: ArmyCategorySpy, Attack: 10, Defence: 1, Capacity: 0, Stealth: 2, Speed: 100, Count: 1},
	}
	op := NewSpyOperation(1, 10, Vector2i{0, 0}, Vector2i{2, 0}, spies)
	op.Start()
	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()

	// Defender spies weak
	def := &UserBaseModel{ID: 99}
	def.ArmiesPresent = []ArmyItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(def.ID),
		Prototype:     ArmyItemPrototype{ID: 700, Category: ArmyCategorySpy, Attack: 1, Defence: 2},
		Count:         1, // def power=2 < atk 10 -> attackers win
	}}

	svc := NewMilitaryOperationService(op, &UserBaseModel{ID: 10})
	svc.ResolveAgainstUserBase(def)

	if op.SpyResult == nil || op.SpyResult.Outcome != SpyOutcomeReportProduced {
		t.Fatalf("expected spy outcome REPORT_PRODUCED, got %+v", op.SpyResult)
	}
	if op.Result != OperationResultSuccess {
		t.Fatalf("expected operation result SUCCESS, got %s", op.Result)
	}
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
