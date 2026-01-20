package domain

import "math/rand"

// Faction represents the origin of an army unit or location defenders.
type Faction string

const (
	FactionExoCoalition      Faction = "EXO_COALITION"   // Playable (Human)
	FactionMarauders         Faction = "MARAUDERS"       // NPC: Credits
	FactionFerrousSwarm      Faction = "FERROUS_SWARM"   // NPC: Iron
	FactionTitanArachnids    Faction = "TITAN_ARACHNIDS" // NPC: Titanium
	FactionVoidEcho          Faction = "VOID_ECHO"       // NPC: Antimatter
	FactionCustodianProtocol Faction = "CUSTODIAN"       // NPC: Dangerous (Artifacts)
	FactionScorchWalkers     Faction = "SCORCH_WALKERS"  // NPC: Dangerous (Buffs)
	FactionObsidianSentinels Faction = "OBSIDIAN"        // NPC: Dangerous (Trophies)
	FactionNeuralWormApex    Faction = "NEURAL_WORM"     // NPC: Dangerous (Intel)
)

const (
	WorthCredit        = 1.0
	WorthIron          = 4.0   // 1000 / 250
	WorthTitanium      = 20.0  // 1000 / 50
	WorthAntimatter    = 333.3 // 1000 / 3
	WorthDefenderPower = 50.0  // 50 credits = 1 defense power
)

// LocationResourceStats represents the available resources at a non-user-base location
// (e.g., resource nodes or dangerous locations). This is similar in spirit to
// UserBaseStats but intentionally scoped down to just the available resource pool
// at that region, which can be looted or otherwise consumed by operations.
type LocationResourceStats struct {
	Credits    int
	Iron       int
	Titanium   int
	Antimatter int

	// Optional bookkeeping to support time-based accumulation if needed later.
	// For now, this lets us compute deltas similarly to how bases do, without
	// introducing production/capacity semantics until required.
	CalculationTimestamp int64 // Unix timestamp of last resource calculation
}

// IsEmpty returns true if all resource counts are zero.
func (stats LocationResourceStats) IsEmpty() bool {
	return stats.Credits == 0 && stats.Iron == 0 && stats.Titanium == 0 && stats.Antimatter == 0
}

// FillFromBudget populates the LocationResourceStats based on a total budget (in credit-equivalents).
// If primaryType is set (not empty/0), that resource gets primaryRatio of the budget.
func (stats *LocationResourceStats) FillFromBudget(totalBudget float64, primaryType ResourceType, primaryRatio float64) {
	stats.CalculationTimestamp = NowUnix()
	if totalBudget <= 0 {
		return
	}

	mainBudget := totalBudget
	othersBudget := 0.0

	// If a primary type is specified, we split the budget.
	if primaryType != "" {
		mainBudget = totalBudget * primaryRatio
		othersBudget = totalBudget - mainBudget
	}

	// 1. Assign the main budget
	switch primaryType {
	case ResourceTypeCredits:
		stats.Credits += int(mainBudget / WorthCredit)
	case ResourceTypeIron:
		stats.Iron += int(mainBudget / WorthIron)
	case ResourceTypeTitanium:
		stats.Titanium += int(mainBudget / WorthTitanium)
	case ResourceTypeAntimatter:
		stats.Antimatter += int(mainBudget / WorthAntimatter)
	default:
		// No primary type, treat everything as "others"
		othersBudget = totalBudget
	}

	// 2. Distribute the "others" budget across all remaining resource types.
	if othersBudget > 0 {
		// Figure out which types are "others"
		othersList := []ResourceType{ResourceTypeCredits, ResourceTypeIron, ResourceTypeTitanium, ResourceTypeAntimatter}
		validOthers := make([]ResourceType, 0, 4)
		for _, t := range othersList {
			if t != primaryType {
				validOthers = append(validOthers, t)
			}
		}

		if len(validOthers) > 0 {
			perTypeBudget := othersBudget / float64(len(validOthers))
			for _, t := range validOthers {
				switch t {
				case ResourceTypeCredits:
					stats.Credits += int(perTypeBudget / WorthCredit)
				case ResourceTypeIron:
					stats.Iron += int(perTypeBudget / WorthIron)
				case ResourceTypeTitanium:
					stats.Titanium += int(perTypeBudget / WorthTitanium)
				case ResourceTypeAntimatter:
					stats.Antimatter += int(perTypeBudget / WorthAntimatter)
				}
			}
		}
	}
}

// FillDefenders populates the provided army and structure stacks based on a total worth budget.
// It uses WorthDefenderPower to calculate the target total combat power (using Defence values).
func FillDefenders(
	armies *[]ArmyStack,
	structures *[]DefenseStack,
	faction Faction,
	totalWorth int,
	armyProtos []*ArmyItemPrototype,
	buildProtos []*BuildItemPrototype,
) {
	if totalWorth <= 0 {
		return
	}

	targetPower := float64(totalWorth) / WorthDefenderPower
	if targetPower < 1 {
		targetPower = 1
	}

	// 1. Filter prototypes by faction
	factionArmies := make([]*ArmyItemPrototype, 0)
	for _, p := range armyProtos {
		if p.Faction == faction {
			factionArmies = append(factionArmies, p)
		}
	}

	factionBuilds := make([]*BuildItemPrototype, 0)
	for _, p := range buildProtos {
		// Only consider buildings that provide defense power (turrets, etc.)
		if p.Faction == faction && p.DefenseData != nil && p.DefenseData.DefenceBonus > 0 {
			factionBuilds = append(factionBuilds, p)
		}
	}

	if len(factionArmies) == 0 && len(factionBuilds) == 0 {
		return
	}

	currentPower := 0.0
	// To avoid infinite loops if power is 0 (though we filtered builds)
	maxIterations := 100
	iterations := 0

	// We'll use maps to aggregate stacks before finalizing
	armyCounts := make(map[int]int)
	buildCounts := make(map[int]int)

	for currentPower < targetPower && iterations < maxIterations {
		iterations++

		// Randomly pick between army and build if both available
		useArmy := len(factionArmies) > 0
		if len(factionArmies) > 0 && len(factionBuilds) > 0 {
			// Weighted towards army 80/20
			useArmy = rand.Intn(10) < 8
		}

		if useArmy {
			p := factionArmies[rand.Intn(len(factionArmies))]
			pwr := p.Defence // Only defence counts for defender power calculation
			if pwr <= 0 {
				pwr = 1 // fallback
			}
			armyCounts[p.ID]++
			currentPower += float64(pwr)
		} else if len(factionBuilds) > 0 {
			p := factionBuilds[rand.Intn(len(factionBuilds))]
			pwr := 0
			if p.DefenseData != nil {
				pwr = p.DefenseData.DefenceBonus
			}
			if pwr <= 0 {
				pwr = 1
			}
			buildCounts[p.ID]++
			currentPower += float64(pwr)
		}
	}

	// Convert maps back to stacks
	for id, count := range armyCounts {
		var proto *ArmyItemPrototype
		for _, p := range factionArmies {
			if p.ID == id {
				proto = p
				break
			}
		}
		if proto != nil {
			*armies = append(*armies, ArmyStack{Prototype: *proto, Count: count})
		}
	}

	for id, count := range buildCounts {
		var proto *BuildItemPrototype
		for _, p := range factionBuilds {
			if p.ID == id {
				proto = p
				break
			}
		}
		if proto != nil {
			*structures = append(*structures, DefenseStack{Prototype: *proto, Count: count})
		}
	}
}
