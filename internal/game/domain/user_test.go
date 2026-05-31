package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUser_DefaultsAndEvents(t *testing.T) {
	id := uuid.New()
	u := NewUser(id, "Alice")

	if u.ID != id {
		t.Fatalf("unexpected id: got %s, want %s", u.ID, id)
	}
	if u.Name != "Alice" {
		t.Fatalf("unexpected name: got %q, want %q", u.Name, "Alice")
	}
	if u.Crystals != DefaultCrystalsBalance {
		t.Fatalf("unexpected starting balance: got %d, want %d", u.Crystals, DefaultCrystalsBalance)
	}

	events := u.PullEvents()
	if len(events) != 2 {
		t.Fatalf("expected 2 events on creation, got %d", len(events))
	}
	if _, ok := events[0].(UserAccountCreatedEvent); !ok {
		t.Fatalf("expected first event UserAccountCreatedEvent, got %T", events[0])
	}
	credited, ok := events[1].(CrystalsCreditedEvent)
	if !ok {
		t.Fatalf("expected second event CrystalsCreditedEvent, got %T", events[1])
	}
	if credited.Reason != CrystalCreditReasonSignupGrant {
		t.Errorf("unexpected reason: got %s, want %s", credited.Reason, CrystalCreditReasonSignupGrant)
	}
	if credited.Amount != DefaultCrystalsBalance {
		t.Errorf("unexpected amount: got %d, want %d", credited.Amount, DefaultCrystalsBalance)
	}
	if credited.BalanceAfter != DefaultCrystalsBalance {
		t.Errorf("unexpected balance after: got %d, want %d", credited.BalanceAfter, DefaultCrystalsBalance)
	}
	if credited.UserID != id {
		t.Errorf("unexpected user id on event: got %s, want %s", credited.UserID, id)
	}
}

func TestUser_AddCrystals_CreditsAndEmits(t *testing.T) {
	u := &User{ID: uuid.New(), Crystals: 10}

	if err := u.AddCrystals(15, CrystalCreditReasonPackPurchase, "order-123"); err != nil {
		t.Fatalf("AddCrystals error: %v", err)
	}
	if u.Crystals != 25 {
		t.Fatalf("unexpected balance: got %d, want %d", u.Crystals, 25)
	}

	events := u.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	ev, ok := events[0].(CrystalsCreditedEvent)
	if !ok {
		t.Fatalf("expected CrystalsCreditedEvent, got %T", events[0])
	}
	if ev.Amount != 15 {
		t.Errorf("unexpected amount: got %d, want %d", ev.Amount, 15)
	}
	if ev.Reason != CrystalCreditReasonPackPurchase {
		t.Errorf("unexpected reason: got %s, want %s", ev.Reason, CrystalCreditReasonPackPurchase)
	}
	if ev.Reference != "order-123" {
		t.Errorf("unexpected reference: got %q, want %q", ev.Reference, "order-123")
	}
	if ev.BalanceAfter != 25 {
		t.Errorf("unexpected balance after: got %d, want %d", ev.BalanceAfter, 25)
	}
}

func TestUser_AddCrystals_RejectsNonPositive(t *testing.T) {
	for _, amount := range []int{0, -5} {
		u := &User{ID: uuid.New(), Crystals: 10}
		if err := u.AddCrystals(amount, CrystalCreditReasonPackPurchase, ""); err == nil {
			t.Fatalf("expected error for amount %d", amount)
		}
		if u.Crystals != 10 {
			t.Errorf("balance changed on rejected add (amount %d): got %d, want 10", amount, u.Crystals)
		}
		if len(u.PullEvents()) != 0 {
			t.Errorf("expected no events on rejected add (amount %d)", amount)
		}
	}
}

func TestUser_SpendCrystals_DeductsAndEmits(t *testing.T) {
	u := &User{ID: uuid.New(), Crystals: 30}

	if err := u.SpendCrystals(12, CrystalSpendReasonBlackMarketArmy, "offer-7"); err != nil {
		t.Fatalf("SpendCrystals error: %v", err)
	}
	if u.Crystals != 18 {
		t.Fatalf("unexpected balance: got %d, want %d", u.Crystals, 18)
	}

	events := u.PullEvents()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	ev, ok := events[0].(CrystalsSpentEvent)
	if !ok {
		t.Fatalf("expected CrystalsSpentEvent, got %T", events[0])
	}
	if ev.Amount != 12 {
		t.Errorf("unexpected amount: got %d, want %d", ev.Amount, 12)
	}
	if ev.Reason != CrystalSpendReasonBlackMarketArmy {
		t.Errorf("unexpected reason: got %s, want %s", ev.Reason, CrystalSpendReasonBlackMarketArmy)
	}
	if ev.Reference != "offer-7" {
		t.Errorf("unexpected reference: got %q, want %q", ev.Reference, "offer-7")
	}
	if ev.BalanceAfter != 18 {
		t.Errorf("unexpected balance after: got %d, want %d", ev.BalanceAfter, 18)
	}
}

func TestUser_SpendCrystals_RejectsNonPositive(t *testing.T) {
	for _, amount := range []int{0, -3} {
		u := &User{ID: uuid.New(), Crystals: 10}
		if err := u.SpendCrystals(amount, CrystalSpendReasonSpeedupBuilding, ""); err == nil {
			t.Fatalf("expected error for amount %d", amount)
		}
		if u.Crystals != 10 {
			t.Errorf("balance changed on rejected spend (amount %d): got %d, want 10", amount, u.Crystals)
		}
		if len(u.PullEvents()) != 0 {
			t.Errorf("expected no events on rejected spend (amount %d)", amount)
		}
	}
}

func TestUser_SpendCrystals_RejectsInsufficientBalance(t *testing.T) {
	u := &User{ID: uuid.New(), Crystals: 5}

	if err := u.SpendCrystals(6, CrystalSpendReasonSpeedupTech, ""); err == nil {
		t.Fatal("expected error when spending more than balance")
	}
	if u.Crystals != 5 {
		t.Errorf("balance changed on insufficient spend: got %d, want 5", u.Crystals)
	}
	if len(u.PullEvents()) != 0 {
		t.Error("expected no events on insufficient spend")
	}
}
