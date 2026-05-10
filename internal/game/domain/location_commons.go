package domain

import (
	"math/rand"
)

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

// FactionForResourceType returns the NPC faction that guards locations of the given resource type.
func FactionForResourceType(r ResourceType) Faction {
	switch r {
	case ResourceTypeIron:
		return FactionFerrousSwarm
	case ResourceTypeTitanium:
		return FactionTitanArachnids
	case ResourceTypeAntimatter:
		return FactionVoidEcho
	default:
		return FactionMarauders
	}
}

const (
	// Base defense power for each location type at a fresh player's progression level.
	baseResourcefulDefense = 10.0
	baseDangerousDefense   = 40.0
	// Attack benchmark for a new player filling starting space with basic infantry (e.g. 100 riflemen = 100 attack).
	spawnStartingPower = 50.0
)

// AppropriateLocationDefense returns the target defense power for a newly spawned
// location near a base. It scales with the player's actual military strength,
// using MaxSpace as a floor for when armies are deployed away.
func AppropriateLocationDefense(stats UserBaseStats, locType LocationType) float64 {
	var base float64
	switch locType {
	case LocationTypeResourceful:
		base = baseResourcefulDefense
	case LocationTypeDangerous:
		base = baseDangerousDefense
	default:
		return 0
	}
	spaceFactor := max(1.0, float64(stats.MaxSpace)/float64(DefaultMaxSpace))
	armyFactor := float64(stats.Attack) / spawnStartingPower
	return base * max(spaceFactor, armyFactor)
}

const (
	WorthCredit             = 1.0
	WorthIron               = 4.0   // 1000 / 250
	WorthTitanium           = 20.0  // 1000 / 50
	WorthAntimatter         = 333.3 // 1000 / 3
	WorthDefenderPower      = 50.0  // 50 credits = 1 defense power
	WorthCapacityMultiplier = 10.0  // 1 capacity point = 10 credits worth of volume
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

// FillDefenders populates the provided army and structure stacks until their combined
// defence power reaches targetPower.
func FillDefenders(
	armies *[]ArmyStack,
	structures *[]DefenseStack,
	faction Faction,
	targetPower float64,
	armyProtos []*ArmyItemPrototype,
	buildProtos []*BuildItemPrototype,
) {
	if targetPower <= 0 {
		return
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

	// 2. Filter out units that are "too strong" for the budget (> 75% of target power).
	// This ensures a diverse force rather than a single super-unit for low-budget locations.
	powerCap := targetPower * 0.75
	filteredArmies := make([]*ArmyItemPrototype, 0)
	for _, p := range factionArmies {
		if float64(p.Defence) <= powerCap {
			filteredArmies = append(filteredArmies, p)
		}
	}
	filteredBuilds := make([]*BuildItemPrototype, 0)
	for _, p := range factionBuilds {
		if float64(p.DefenseData.DefenceBonus) <= powerCap {
			filteredBuilds = append(filteredBuilds, p)
		}
	}
	// Fallback to full lists if filtering was too aggressive (e.g. at very low power levels where all units exceed the cap)
	if len(filteredArmies) > 0 || len(filteredBuilds) > 0 {
		factionArmies = filteredArmies
		factionBuilds = filteredBuilds
	}

	currentPower := 0.0
	armyCounts := make(map[int]int)
	buildCounts := make(map[int]int)

	for currentPower < targetPower {

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
