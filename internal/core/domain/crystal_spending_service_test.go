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
