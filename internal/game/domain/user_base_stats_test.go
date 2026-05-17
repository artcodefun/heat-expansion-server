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

func TestStats_CanReceiveResourceAmount(t *testing.T) {
	stats := UserBaseStats{
		Credits:            100,
		CreditsCapacity:    150,
		Iron:               50,
		IronCapacity:       60,
		Titanium:           10,
		TitaniumCapacity:   20,
		Antimatter:         1,
		AntimatterCapacity: 2,
	}

	if !stats.CanReceiveResourceAmount(ResourceTypeCredits, 50) {
		t.Fatal("expected credits to fit")
	}
	if stats.CanReceiveResourceAmount(ResourceTypeIron, 11) {
		t.Fatal("expected iron to exceed capacity")
	}
	if !stats.CanReceiveResourceAmount(ResourceTypeTitanium, 10) {
		t.Fatal("expected titanium to fit exactly")
	}
	if stats.CanReceiveResourceAmount(ResourceTypeAntimatter, 2) {
		t.Fatal("expected antimatter to exceed capacity")
	}
}
