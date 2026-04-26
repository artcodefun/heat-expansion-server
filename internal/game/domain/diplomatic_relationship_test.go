package domain

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func newNeutralRelationshipForTest(t *testing.T, userA, userB, changedBy uuid.UUID) *DiplomaticRelationship {
	t.Helper()
	rel, err := NewUnknownRelationship(userA, userB)
	if err != nil {
		t.Fatalf("unexpected unknown relationship init error: %v", err)
	}
	if err := rel.EstablishContact(changedBy); err != nil {
		t.Fatalf("unexpected establish contact error: %v", err)
	}
	return rel
}

func TestCanPerformAttackOperation_AllowsUnknownRelationship(t *testing.T) {
	SetTestNow(t, 5_000)
	attackerUserID := uuid.New()
	defenderUserID := uuid.New()

	rel, err := NewUnknownRelationship(attackerUserID, defenderUserID)
	if err != nil {
		t.Fatalf("unexpected unknown relationship init error: %v", err)
	}
	err = rel.CanPerformAttackOperation()
	if err != nil {
		t.Fatalf("expected unknown relationship to allow attack operation, got %v", err)
	}
}

func TestCanPerformAttackOperation_AllowsWar(t *testing.T) {
	SetTestNow(t, 6_000)
	attackerUserID := uuid.New()
	defenderUserID := uuid.New()
	rel := newNeutralRelationshipForTest(t, attackerUserID, defenderUserID, attackerUserID)
	if err := rel.DeclareWar(attackerUserID); err != nil {
		t.Fatalf("unexpected war declaration error: %v", err)
	}
	SetTestNow(t, *rel.WarAttacksAllowedAt)

	err := rel.CanPerformAttackOperation()
	if err != nil {
		t.Fatalf("expected war relationship to allow attack, got %v", err)
	}
}

func TestCanPerformTradeOperation_RequiresAlliance(t *testing.T) {
	SetTestNow(t, 6_100)
	userA := uuid.New()
	userB := uuid.New()
	rel := newNeutralRelationshipForTest(t, userA, userB, userA)

	err := rel.CanPerformTradeOperation()
	if err == nil {
		t.Fatalf("expected neutral relationship to reject trade operation")
	}
	if !strings.HasPrefix(err.Error(), "error.domain.diplomacy.trade_requires_alliance") {
		t.Fatalf("expected dedicated trade alliance error, got %v", err)
	}
}

func TestCanPerformTradeOperation_AllowsAlliance(t *testing.T) {
	SetTestNow(t, 6_200)
	userA := uuid.New()
	userB := uuid.New()
	rel := newNeutralRelationshipForTest(t, userA, userB, userA)
	rel.FormAlliance(userA)

	if err := rel.CanPerformTradeOperation(); err != nil {
		t.Fatalf("expected allied relationship to allow trade operation, got %v", err)
	}
}

func TestCanDeclareWar_RejectsMissingRelationship(t *testing.T) {
	SetTestNow(t, 7_000)
	rel, err := NewUnknownRelationship(uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected unknown relationship init error: %v", err)
	}

	if err := rel.CanDeclareWar(); err == nil {
		t.Fatalf("expected missing relationship declaration to be rejected")
	}
}

func TestEscalateToWar_UnknownRelationshipBecomesWar(t *testing.T) {
	SetTestNow(t, 7_100)
	rel, err := NewUnknownRelationship(uuid.New(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected unknown relationship init error: %v", err)
	}

	attackerUserID := uuid.New()
	if err := rel.EscalateToWar(attackerUserID); err != nil {
		t.Fatalf("expected unknown relationship escalation to succeed, got %v", err)
	}
	if rel.Status != DiplomaticStatusWar {
		t.Fatalf("expected relationship to become war, got %s", rel.Status)
	}
	if rel.ID == uuid.Nil {
		t.Fatalf("expected transformed war relationship to receive persistent id")
	}
	if rel.ChangedByUserID != attackerUserID {
		t.Fatalf("expected transformed war relationship to record attacker user id")
	}
	if rel.WarAttacksAllowedAt == nil || *rel.WarAttacksAllowedAt != 7_100 {
		t.Fatalf("expected transformed war relationship to allow attacks immediately")
	}
}

func TestEscalateToWar_RejectsEstablishedRelationship(t *testing.T) {
	SetTestNow(t, 7_200)
	rel := newNeutralRelationshipForTest(t, uuid.New(), uuid.New(), uuid.New())

	if err := rel.EscalateToWar(uuid.New()); err == nil {
		t.Fatalf("expected established relationship escalation to be rejected")
	}
}

func TestCanDeclareWar_RejectsAlliedRelationship(t *testing.T) {
	SetTestNow(t, 8_000)
	userA := uuid.New()
	userB := uuid.New()
	rel := newNeutralRelationshipForTest(t, userA, userB, userA)
	rel.FormAlliance(userA)

	if err := rel.CanDeclareWar(); err == nil {
		t.Fatalf("expected allied relationship to reject war declaration")
	}
}

func TestBreakAlliance_RequiresAlliance(t *testing.T) {
	SetTestNow(t, 9_000)
	userA := uuid.New()
	userB := uuid.New()
	rel := newNeutralRelationshipForTest(t, userA, userB, userA)

	if err := rel.BreakAlliance(userA); err == nil {
		t.Fatalf("expected non-allied relationship to reject alliance break")
	}
}

func TestBreakAlliance_ProtectsNeutrality(t *testing.T) {
	SetTestNow(t, 9_100)
	userA := uuid.New()
	userB := uuid.New()
	rel := newNeutralRelationshipForTest(t, userA, userB, userA)
	rel.FormAlliance(userA)

	if err := rel.BreakAlliance(userA); err != nil {
		t.Fatalf("expected allied relationship break to succeed, got %v", err)
	}
	if rel.Status != DiplomaticStatusNeutral {
		t.Fatalf("expected relationship to become neutral, got %s", rel.Status)
	}
	if rel.NeutralityProtectedUntil == nil || *rel.NeutralityProtectedUntil != 9_100+int64(DiplomaticCeasefireProtectionDuration.Seconds()) {
		t.Fatalf("expected neutrality protection to be set, got %+v", rel.NeutralityProtectedUntil)
	}

	if err := rel.CanDeclareWar(); err == nil {
		t.Fatalf("expected neutrality protection to block war declaration")
	}

	SetTestNow(t, *rel.NeutralityProtectedUntil)
	if err := rel.CanDeclareWar(); err != nil {
		t.Fatalf("expected war declaration after protection expiry, got %v", err)
	}
}
