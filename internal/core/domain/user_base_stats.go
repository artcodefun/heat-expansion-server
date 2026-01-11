package domain

import (
	"fmt"
)

// Default capacities and stats for UserBaseStats
const (
	DefaultCreditsCapacity       = 10000
	DefaultIronCapacity          = 5000
	DefaultTitaniumCapacity      = 2500
	DefaultAntimatterCapacity    = 1000
	DefaultMaxSpace              = 50
	DefaultMaxOperations         = 3
	DefaultMaxActiveBuffs        = 2
	DefaultMaxActiveArtifacts    = 1
	DefaultMaxBuildingProduction = 1
	DefaultMaxActiveRestorations = 1
	DefaultMaxActiveDecryptions  = 1
)

// UserBaseStats represents current properties of a base.
type UserBaseStats struct {
	Credits               int
	CreditsCapacity       int
	CreditsProduction     float64
	Iron                  int
	IronCapacity          int
	IronProduction        float64
	Titanium              int
	TitaniumCapacity      int
	TitaniumProduction    float64
	Antimatter            int
	AntimatterCapacity    int
	AntimatterProduction  float64
	Defence               int
	Attack                int
	Space                 int
	MaxSpace              int
	MaxOperations         int
	MaxActiveBuffs        int
	MaxActiveArtifacts    int
	MaxBuildingProduction int
	MaxActiveRestorations int
	MaxActiveDecryptions  int
	CalculationTimestamp  int64 // Unix timestamp of last resource calculation
}

func (s *UserBaseStats) CheckResources(price PriceModel) error {
	if price.Credits > s.Credits {
		return fmt.Errorf("insufficient credits")
	}
	if price.Iron > s.Iron {
		return fmt.Errorf("insufficient iron")
	}
	if price.Titanium > s.Titanium {
		return fmt.Errorf("insufficient titanium")
	}
	if price.Antimatter > s.Antimatter {
		return fmt.Errorf("insufficient antimatter")
	}
	return nil
}

func (s *UserBaseStats) SubtractResources(price PriceModel) {
	s.Credits -= price.Credits
	s.Iron -= price.Iron
	s.Titanium -= price.Titanium
	s.Antimatter -= price.Antimatter
}

// RecalculateStats updates the UserBaseStats based on present items and default constants.
func (ub *UserBaseModel) recalculateStats() {
	stats := UserBaseStats{}
	// Set default capacities
	stats.CreditsCapacity = DefaultCreditsCapacity
	stats.IronCapacity = DefaultIronCapacity
	stats.TitaniumCapacity = DefaultTitaniumCapacity
	stats.AntimatterCapacity = DefaultAntimatterCapacity
	stats.MaxSpace = DefaultMaxSpace
	stats.MaxOperations = DefaultMaxOperations
	stats.MaxActiveBuffs = DefaultMaxActiveBuffs
	stats.MaxActiveArtifacts = DefaultMaxActiveArtifacts
	stats.MaxBuildingProduction = DefaultMaxBuildingProduction
	stats.MaxActiveRestorations = DefaultMaxActiveRestorations
	stats.MaxActiveDecryptions = DefaultMaxActiveDecryptions

	// Aggregate bonuses from present buildings
	for _, b := range ub.BuildingsPresent {
		proto := b.Prototype
		// Resources buildings
		if proto.ResourcesData != nil {
			stats.CreditsCapacity += proto.ResourcesData.CreditsCapacity
			stats.IronCapacity += proto.ResourcesData.IronCapacity
			stats.TitaniumCapacity += proto.ResourcesData.TitaniumCapacity
			stats.AntimatterCapacity += proto.ResourcesData.AntimatterCapacity
			stats.CreditsProduction += proto.ResourcesData.CreditsProduction
			stats.IronProduction += proto.ResourcesData.IronProduction
			stats.TitaniumProduction += proto.ResourcesData.TitaniumProduction
			stats.AntimatterProduction += proto.ResourcesData.AntimatterProduction
		}
		// Defense buildings
		if proto.DefenseData != nil {
			stats.Defence += proto.DefenseData.DefenceBonus
		}
		// Space is always added
		stats.Space += proto.Space
	}

	// Include space from buildings in production
	for _, b := range ub.BuildingsInProduction {
		stats.Space += b.Prototype.Space
	}

	// Include space from buildings pending
	for _, b := range ub.BuildingsPending {
		stats.Space += b.Prototype.Space
	}
	// Include space from armies present
	for _, a := range ub.ArmiesPresent {
		stats.Space += a.Prototype.Space * a.Count
	}

	// Include space from armies deployed (still occupy capacity)
	for _, d := range ub.ArmiesDeployed {
		stats.Space += d.Prototype.Space * d.Count
	}

	// Include space from armies in production
	for _, a := range ub.ArmiesInProduction {
		stats.Space += a.Prototype.Space
	}

	// Include space from armies pending
	for _, a := range ub.ArmiesPending {
		stats.Space += a.Prototype.Space * a.Count
	}

	// Aggregate power from present armies
	for _, a := range ub.ArmiesPresent {
		stats.Defence += a.Prototype.Defence
		stats.Attack += a.Prototype.Attack
	}

	// Apply researched technology improvements (additive bonuses scaling with level)
	for _, tech := range ub.TechnologiesDone {
		imp := tech.Prototype.Improvement
		if imp == nil {
			continue
		}
		value := imp.Value * tech.Level
		switch imp.Type {
		case ImprovementTypeSpaceCapacity:
			stats.MaxSpace += value
		case ImprovementTypeOperationsCount:
			stats.MaxOperations += value
		case ImprovementTypeActiveBuffsCount:
			stats.MaxActiveBuffs += value
		case ImprovementTypeActiveArtifactsCount:
			stats.MaxActiveArtifacts += value
		case ImprovementTypeActiveRestorationsCount:
			stats.MaxActiveRestorations += value
		case ImprovementTypeBuildingProductionCount:
			stats.MaxBuildingProduction += value
		case ImprovementTypeActiveDecryptionsCount:
			stats.MaxActiveDecryptions += value
		}
	}

	// Apply modifiers from storage items (buffs and artifacts)
	mods := ub.ActiveModifiers()
	stats.CreditsProduction *= mods.CreditsProdMul
	stats.IronProduction *= mods.IronProdMul
	stats.TitaniumProduction *= mods.TitaniumProdMul
	// (Antimatter production doesn't have a specific multiplier in the current BuffTypes/ArtifactTypes)

	stats.Attack = mulInt(stats.Attack, mods.AttackMul)
	stats.Defence = mulInt(stats.Defence, mods.DefenceMul)

	// Calculate current resources based on previous value, production rate, and elapsed time
	prevStats := ub.Stats
	now := NowUnix()
	delta := now - prevStats.CalculationTimestamp

	if delta > 0 {
		stats.Credits = prevStats.Credits + int(stats.CreditsProduction*float64(delta))
		if stats.Credits > stats.CreditsCapacity {
			stats.Credits = stats.CreditsCapacity
		}

		stats.Iron = prevStats.Iron + int(stats.IronProduction*float64(delta))
		if stats.Iron > stats.IronCapacity {
			stats.Iron = stats.IronCapacity
		}

		stats.Titanium = prevStats.Titanium + int(stats.TitaniumProduction*float64(delta))
		if stats.Titanium > stats.TitaniumCapacity {
			stats.Titanium = stats.TitaniumCapacity
		}

		stats.Antimatter = prevStats.Antimatter + int(stats.AntimatterProduction*float64(delta))
		if stats.Antimatter > stats.AntimatterCapacity {
			stats.Antimatter = stats.AntimatterCapacity
		}
	} else {
		stats.Credits = prevStats.Credits
		stats.Iron = prevStats.Iron
		stats.Titanium = prevStats.Titanium
		stats.Antimatter = prevStats.Antimatter
	}

	stats.CalculationTimestamp = now
	ub.Stats = stats
}

// DeductLoot subtracts the provided loot from the base's resources, clamped at zero.
func (ub *UserBaseModel) DeductLoot(loot PriceModel) {
	if loot.Credits > 0 {
		ub.Stats.Credits = maxInt(ub.Stats.Credits-loot.Credits, 0)
	}
	if loot.Iron > 0 {
		ub.Stats.Iron = maxInt(ub.Stats.Iron-loot.Iron, 0)
	}
	if loot.Titanium > 0 {
		ub.Stats.Titanium = maxInt(ub.Stats.Titanium-loot.Titanium, 0)
	}
	if loot.Antimatter > 0 {
		ub.Stats.Antimatter = maxInt(ub.Stats.Antimatter-loot.Antimatter, 0)
	}
}

// CreditLoot adds the provided loot to the base's resources, clamped by capacities.
func (ub *UserBaseModel) CreditLoot(loot PriceModel) {
	if loot.Credits > 0 {
		ub.Stats.Credits = min(ub.Stats.Credits+loot.Credits, ub.Stats.CreditsCapacity)
	}
	if loot.Iron > 0 {
		ub.Stats.Iron = min(ub.Stats.Iron+loot.Iron, ub.Stats.IronCapacity)
	}
	if loot.Titanium > 0 {
		ub.Stats.Titanium = min(ub.Stats.Titanium+loot.Titanium, ub.Stats.TitaniumCapacity)
	}
	if loot.Antimatter > 0 {
		ub.Stats.Antimatter = min(ub.Stats.Antimatter+loot.Antimatter, ub.Stats.AntimatterCapacity)
	}
}
