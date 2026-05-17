package domain

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestBlackMarketService_PurchaseResources_SpendsCrystalsAndCreditsBase(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 10}
	base := newBaseWithDefaults(1)
	base.Stats.Credits = 0

	if err := service.PurchaseResources(user, base, ResourceTypeCredits, 2); err != nil {
		t.Fatalf("PurchaseResources error: %v", err)
	}

	if user.Crystals != 8 {
		t.Fatalf("expected 8 crystals remaining, got %d", user.Crystals)
	}
	if base.Stats.Credits != 200 {
		t.Fatalf("expected credits to increase by 200, got %v", base.Stats.Credits)
	}
}

func TestBlackMarketService_PurchaseResources_RejectsWhenBaseCannotReceive(t *testing.T) {
	service := NewBlackMarketService()
	user := &User{ID: uuid.New(), Crystals: 10}
	base := newBaseWithDefaults(2)
	base.Stats.Iron = float64(base.Stats.IronCapacity)
	startingCrystals := user.Crystals
	startingIron := base.Stats.Iron

	err := service.PurchaseResources(user, base, ResourceTypeIron, 1)
	if err == nil {
		t.Fatal("expected capacity error")
	}
	var domainErr Error
	if !errors.As(err, &domainErr) {
		t.Fatalf("expected domain error, got %T", err)
	}
	if domainErr.Key != "error.domain.black_market.resource_capacity_reached" {
		t.Fatalf("expected black market capacity key, got %s", domainErr.Key)
	}
	if user.Crystals != startingCrystals {
		t.Fatalf("expected crystals unchanged, got %d", user.Crystals)
	}
	if base.Stats.Iron != startingIron {
		t.Fatalf("expected iron unchanged, got %v", base.Stats.Iron)
	}
}
