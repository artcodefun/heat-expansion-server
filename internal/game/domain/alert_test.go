package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewActivityAlert_DefenseAttack(t *testing.T) {
	SetTestNow(t, 1_000)

	event := NewActivityCreatedEvent(uuid.New(), uuid.New(), 42, ActivityKindDefense, string(DefenseActivitySubtypeAttack))
	alert, ok := NewActivityAlert(event)
	if !ok {
		t.Fatalf("expected activity alert to be created")
	}
	if alert.Kind != AlertKindCombat {
		t.Fatalf("expected combat alert, got %s", alert.Kind)
	}
	if alert.Title != "alert.combat.attack.title" {
		t.Fatalf("unexpected title key: %s", alert.Title)
	}
	if alert.Content != "alert.combat.attack.content" {
		t.Fatalf("unexpected content key: %s", alert.Content)
	}
	if alert.BaseID == nil || *alert.BaseID != 42 {
		t.Fatalf("expected base id 42, got %+v", alert.BaseID)
	}
}

func TestNewDiplomaticRequestAlert_CeasefireProposal(t *testing.T) {
	SetTestNow(t, 2_000)
	baseID := 7
	event := NewDiplomaticRequestCreatedEvent(uuid.New(), uuid.New(), uuid.New(), &baseID, DiplomaticRequestKindCeasefireProposal)

	alert, ok := NewDiplomaticRequestAlert(event)
	if !ok {
		t.Fatalf("expected diplomatic request alert to be created")
	}
	if alert.Title != "alert.diplomacy.ceasefire_proposal.title" {
		t.Fatalf("unexpected title key: %s", alert.Title)
	}
	if alert.Content != "alert.diplomacy.ceasefire_proposal.content" {
		t.Fatalf("unexpected content key: %s", alert.Content)
	}
}
