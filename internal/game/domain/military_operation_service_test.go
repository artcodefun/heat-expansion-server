package domain

import (
	"testing"

	"github.com/google/uuid"
)

// TestResolveAgainstUserBase_Attack_AppliesLootAndSurvivors verifies that
// MilitaryOperationService resolves an attack against a user base by:
// - delegating combat to the operation
// - trimming deployed units to survivors on the attacker
// - applying remaining defenders/structures to the defender
// - deducting loot from defender resources
// - starting the return leg.
func TestResolveAgainstUserBase_Attack_AppliesLootAndSurvivors(t *testing.T) {
	SetTestNow(t, 10_000)

	// Simple attacking force with some capacity for loot
	attackUnits := []MilitaryUnitSnap{
		{PrototypeID: 1, Category: ArmyCategoryInfantry, Attack: 5, Defence: 3, Capacity: 2, Stealth: 0, Speed: 100, Count: 3},
	}

	op, err := NewAttackOperation(uuid.New(), 10, Vector2i{X: 0, Y: 0}, Vector2i{X: 1, Y: 0}, attackUnits, nil)
	if err != nil {
		t.Fatalf("unexpected error from NewAttackOperation: %v", err)
	}
	op.Start()

	// Fast-forward to arrival so ResolveAgainstUserBase can call OnArrive safely.
	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()

	// Attacker base with some deployed units (all units deployed to this operation).
	attacker := &UserBaseModel{ID: 10, UserID: uuid.New()}
	attacker.ArmiesDeployed = []ArmyItemDeployed{
		{
			BaseOwnedItem: NewBaseOwnedItem(attacker.ID),
			Prototype: ArmyItemPrototype{
				ID:       1,
				Category: ArmyCategoryInfantry,
				Attack:   5,
				Defence:  3,
				Capacity: 2,
				Speed:    100,
			},
			OperationKind: OperationKindMilitary,
			OperationID:   op.ID,
			Count:         3,
		},
	}

	// Defender with modest resources and a small defending force.
	defender := &UserBaseModel{ID: 99, UserID: uuid.New()}
	defender.Stats = UserBaseStats{Credits: 10, Iron: 5, Titanium: 3, Antimatter: 2, CalculationTimestamp: NowUnix()}
	defender.ArmiesPresent = []ArmyItemPresent{
		{
			BaseOwnedItem: NewBaseOwnedItem(defender.ID),
			Prototype:     ArmyItemPrototype{ID: 2, Category: ArmyCategoryInfantry, Attack: 3, Defence: 2},
			Count:         2,
		},
	}

	// No defensive structures for this scenario.

	service := NewMilitaryOperationService(op, attacker)
	service.ResolveAgainstUserBase(defender)

	// After resolution we expect an AttackResult set.
	if op.AttackResult == nil {
		t.Fatalf("expected AttackResult to be set")
	}

	// Loot must be non-negative and not exceed defender's initial total resources.
	loot := op.AttackResult.Loot
	lootTotal := loot.Credits + loot.Iron + loot.Titanium + loot.Antimatter
	if lootTotal < 0 {
		t.Fatalf("loot total should be non-negative, got %d", lootTotal)
	}
	initial := 10 + 5 + 3 + 2
	remaining := int(defender.Stats.Credits + defender.Stats.Iron + defender.Stats.Titanium + defender.Stats.Antimatter)
	if initial-remaining != lootTotal {
		t.Fatalf("expected defender resource sum decrease by %d, got %d", lootTotal, initial-remaining)
	}

	// Attacker's deployed units should have been trimmed to survivors of this operation.
	// We don't assert exact numbers (combat is more complex), but we at least ensure
	// that the operation ID entry still exists and count is between 0 and original.
	found := false
	for _, d := range attacker.ArmiesDeployed {
		if d.OperationID == op.ID {
			found = true
			if d.Count < 0 || d.Count > 3 {
				t.Fatalf("unexpected survivor count %d for deployed units", d.Count)
			}
		}
	}
	if !found {
		t.Fatalf("expected deployed entry for operation %d to remain after trimming", op.ID)
	}

	// Operation should have started its return leg.
	if op.Phase != OperationPhaseReturning && op.Phase != OperationPhaseCompleted {
		t.Fatalf("expected operation to be returning or completed, got %s", op.Phase)
	}
}

// TestResolveAgainstUserBase_Spy_BlockedByCloaking_PreservesNonSpyDefenders verifies that
// the service correctly uses cloaking strength and preserves non-spy defenders
// when a spy operation is blocked by cloaking.
func TestResolveAgainstUserBase_Spy_BlockedByCloaking_PreservesNonSpyDefenders(t *testing.T) {
	SetTestNow(t, 20_000)

	// Attacker spies (stealth sum = 8)
	spies := []MilitaryUnitSnap{
		{PrototypeID: 7, Category: ArmyCategorySpy, Attack: 2, Defence: 1, Capacity: 0, Stealth: 4, Speed: 120, Count: 2},
	}
	op, err := NewSpyOperation(uuid.New(), 10, Vector2i{X: 0, Y: 0}, Vector2i{X: 1, Y: 1}, spies, nil)
	if err != nil {
		t.Fatalf("unexpected error from NewSpyOperation: %v", err)
	}
	op.Start()
	SetTestNow(t, op.OutboundArriveAt)
	op.UpdatePhaseBasedOnTime()

	// Defender with strong cloaking (>= attacker stealth) and a non-spy garrison
	defender := &UserBaseModel{ID: 99, UserID: uuid.New()}
	defender.BuildingsPresent = []BuildItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(defender.ID),
		Prototype: BuildItemPrototype{
			ID:       300,
			Category: BuildCategoryIntelligence,
			IntelligenceData: &IntelligenceBuildingData{
				Subtype:         IntelligenceSubtypeCloaking,
				StealthStrength: 10,
			},
		},
	}}
	defender.ArmiesPresent = []ArmyItemPresent{{
		BaseOwnedItem: NewBaseOwnedItem(defender.ID),
		Prototype:     ArmyItemPrototype{ID: 500, Category: ArmyCategoryInfantry, Attack: 5, Defence: 5},
		Count:         3,
	}}

	service := NewMilitaryOperationService(op, &UserBaseModel{ID: 10, UserID: uuid.New()})
	service.ResolveAgainstUserBase(defender)

	if op.SpyResult == nil || op.SpyResult.Outcome != SpyOutcomeBlockedByCloaking {
		t.Fatalf("expected spy outcome BLOCKED_BY_CLOAKING, got %+v", op.SpyResult)
	}
	if op.Result != OperationResultSuccess {
		t.Fatalf("expected operation result SUCCESS, got %s", op.Result)
	}
	if op.Phase != OperationPhaseReturning {
		t.Fatalf("expected returning phase after resolve, got %s", op.Phase)
	}

	// Non-spy defenders must remain unchanged (spy merge only touches spies)
	if len(defender.ArmiesPresent) != 1 || defender.ArmiesPresent[0].Prototype.Category != ArmyCategoryInfantry || defender.ArmiesPresent[0].Count != 3 {
		t.Fatalf("unexpected defender armies state: %+v", defender.ArmiesPresent)
	}
}
