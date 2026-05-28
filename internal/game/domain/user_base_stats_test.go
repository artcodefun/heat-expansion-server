package domain

import (
	"testing"
)

func TestStats_Recalculate_DefaultValues(t *testing.T) {
	base := newBaseWithDefaults(1)
	stats := base.Stats

	if stats.Credits != 1000 {
		t.Errorf("expected 1000 credits, got %v", stats.Credits)
	}
	if stats.CreditsCapacity != DefaultCreditsCapacity {
		t.Errorf("expected %d capacity, got %d", DefaultCreditsCapacity, stats.CreditsCapacity)
	}
	if stats.MaxSpace != DefaultMaxSpace {
		t.Errorf("expected %d max space, got %d", DefaultMaxSpace, stats.MaxSpace)
	}
}

func TestStats_Recalculate_WithBuildings(t *testing.T) {
	base := newBaseWithDefaults(1)

	// Add a warehouse to increase iron capacity
	base.BuildingsPresent = append(base.BuildingsPresent, BuildItemPresent{
		Prototype: BuildItemPrototype{
			ID:    1,
			Space: 5,
			ResourcesData: &ResourcesBuildingData{
				IronCapacity: 1000,
			},
		},
	})

	base.recalculateStats()

	expectedCapacity := DefaultIronCapacity + 1000
	if base.Stats.IronCapacity != expectedCapacity {
		t.Errorf("expected iron capacity %d, got %d", expectedCapacity, base.Stats.IronCapacity)
	}

	// Space should be increased by building space
	if base.Stats.Space != 5 {
		t.Errorf("expected space 5, got %d", base.Stats.Space)
	}
}

func TestStats_ProductionOverTime(t *testing.T) {
	base := newBaseWithDefaults(1)
	base.Stats.Credits = 0
	base.Stats.CalculationTimestamp = NowUnix() - 10 // 10 seconds ago

	// Add building with production
	base.BuildingsPresent = append(base.BuildingsPresent, BuildItemPresent{
		Prototype: BuildItemPrototype{
			ResourcesData: &ResourcesBuildingData{
				CreditsProduction: 10.0,
			},
		},
	})

	base.recalculateStats()

	// 10 seconds * 10 credits/sec = 100 credits
	if base.Stats.Credits != 100 {
		t.Errorf("expected 100 credits, got %v", base.Stats.Credits)
	}
}

func TestStats_Production_ClampedToCapacity(t *testing.T) {
	base := newBaseWithDefaults(1)
	base.Stats.Credits = float64(DefaultCreditsCapacity - 50)
	base.Stats.CalculationTimestamp = NowUnix() - 100 // a long time ago

	base.BuildingsPresent = append(base.BuildingsPresent, BuildItemPresent{
		Prototype: BuildItemPrototype{
			ResourcesData: &ResourcesBuildingData{
				CreditsProduction: 10.0,
			},
		},
	})

	base.recalculateStats()

	if base.Stats.Credits != float64(DefaultCreditsCapacity) {
		t.Errorf("expected credits clamped to %d, got %v", DefaultCreditsCapacity, base.Stats.Credits)
	}
}

func TestStats_CheckResources(t *testing.T) {
	stats := UserBaseStats{
		Credits: 100,
		Iron:    50,
	}

	price := PriceModel{Credits: 150}
	if err := stats.CheckResources(price); err == nil {
		t.Error("expected error for insufficient credits")
	}

	price = PriceModel{Credits: 50, Iron: 20}
	if err := stats.CheckResources(price); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestStats_CreditAndDeductLoot(t *testing.T) {
	base := newBaseWithDefaults(1)
	base.Stats.CreditsCapacity = 2000
	base.Stats.IronCapacity = 2000
	base.Stats.Credits = 1000
	base.Stats.Iron = 1000

	loot := PriceModel{Credits: 500, Iron: 200}

	base.DeductLoot(loot)
	if base.Stats.Credits != 500 || base.Stats.Iron != 800 {
		t.Errorf("DeductLoot failed: credits=%v, iron=%v", base.Stats.Credits, base.Stats.Iron)
	}

	base.CreditLoot(loot)
	if base.Stats.Credits != 1000 || base.Stats.Iron != 1000 {
		t.Errorf("CreditLoot failed: credits=%v, iron=%v", base.Stats.Credits, base.Stats.Iron)
	}
}

func TestReceiveResource_RespectsCapacity(t *testing.T) {
	base := newBaseWithDefaults(1)
	base.Stats.Credits = 100
	base.Stats.CreditsCapacity = 150

	if err := base.ReceiveResource(ResourceTypeCredits, 50); err != nil {
		t.Fatalf("expected credits to fit, got %v", err)
	}
	if base.Stats.Credits != 150 {
		t.Fatalf("expected credits to reach capacity, got %v", base.Stats.Credits)
	}
	if err := base.ReceiveResource(ResourceTypeCredits, 1); err == nil {
		t.Fatal("expected credits overflow to fail")
	}
	if base.Stats.Credits != 150 {
		t.Fatalf("expected credits unchanged after overflow, got %v", base.Stats.Credits)
	}

	base.Stats.Iron = 50
	base.Stats.IronCapacity = 60
	if err := base.ReceiveResource(ResourceTypeIron, 11); err == nil {
		t.Fatal("expected iron overflow to fail")
	}
	if base.Stats.Iron != 50 {
		t.Fatalf("expected iron unchanged after overflow, got %v", base.Stats.Iron)
	}

	base.Stats.Titanium = 10
	base.Stats.TitaniumCapacity = 20
	if err := base.ReceiveResource(ResourceTypeTitanium, 10); err != nil {
		t.Fatalf("expected titanium to fit exactly, got %v", err)
	}
	if base.Stats.Titanium != 20 {
		t.Fatalf("expected titanium to reach capacity, got %v", base.Stats.Titanium)
	}

	base.Stats.Antimatter = 1
	base.Stats.AntimatterCapacity = 2
	if err := base.ReceiveResource(ResourceTypeAntimatter, 2); err == nil {
		t.Fatal("expected antimatter overflow to fail")
	}
	if base.Stats.Antimatter != 1 {
		t.Fatalf("expected antimatter unchanged after overflow, got %v", base.Stats.Antimatter)
	}
}

func TestReceiveResource_RejectsInvalidInputs(t *testing.T) {
	base := newBaseWithDefaults(1)

	if err := base.ReceiveResource(ResourceTypeCredits, 0); err == nil {
		t.Fatal("expected zero amount to fail")
	}
	if err := base.ReceiveResource(ResourceType("INVALID"), 1); err == nil {
		t.Fatal("expected invalid resource type to fail")
	}
}
