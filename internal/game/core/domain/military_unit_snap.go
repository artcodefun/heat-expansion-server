package domain

import "math"

// MilitaryUnitSnap captures a snapshot of an army unit participating in an operation.
// Using a snapshot decouples battle resolution from prototype changes over time.
type MilitaryUnitSnap struct {
	PrototypeID int
	Category    ArmyCategory
	Attack      int
	Defence     int
	Capacity    int
	Stealth     int
	Speed       int
	Count       int
}

// DefenseStructureSnap is a simplified snapshot of a defending structure (e.g., turrets, shields).
// This provides a lightweight input for attack resolution without coupling to build prototypes.
type DefenseStructureSnap struct {
	PrototypeID int
	Defence     int
	Count       int
}

// MilitaryUnitFromPresent converts a present army item to an MilitaryUnitSnap snapshot.
func MilitaryUnitFromPresent(p ArmyItemPresent) MilitaryUnitSnap {
	return MilitaryUnitSnap{
		PrototypeID: p.Prototype.ID,
		Category:    p.Prototype.Category,
		Attack:      p.Prototype.Attack,
		Defence:     p.Prototype.Defence,
		Capacity:    p.Prototype.Capacity,
		Stealth:     p.Prototype.Stealth,
		Speed:       p.Prototype.Speed,
		Count:       p.Count,
	}
}

// MilitaryUnitFromStack converts an army stack into an operational snapshot.
func MilitaryUnitFromStack(stack ArmyStack) MilitaryUnitSnap {
	return MilitaryUnitSnap{
		PrototypeID: stack.Prototype.ID,
		Category:    stack.Prototype.Category,
		Attack:      stack.Prototype.Attack,
		Defence:     stack.Prototype.Defence,
		Capacity:    stack.Prototype.Capacity,
		Stealth:     stack.Prototype.Stealth,
		Speed:       stack.Prototype.Speed,
		Count:       stack.Count,
	}
}

// MilitaryUnitsFromPresent maps present army items to MilitaryUnitSnap snapshots.
func MilitaryUnitsFromPresent(items []ArmyItemPresent) []MilitaryUnitSnap {
	if len(items) == 0 {
		return nil
	}
	out := make([]MilitaryUnitSnap, 0, len(items))
	for _, it := range items {
		if it.Count > 0 {
			out = append(out, MilitaryUnitFromPresent(it))
		}
	}
	return out
}

// MilitaryUnitFromDeployed converts a single deploy-ready stack into an military unit snapshot.
func MilitaryUnitFromDeployed(d DeploymentReadyItem) MilitaryUnitSnap {
	return MilitaryUnitSnap{
		PrototypeID: d.Prototype.ID,
		Category:    d.Prototype.Category,
		Attack:      d.Prototype.Attack,
		Defence:     d.Prototype.Defence,
		Capacity:    d.Prototype.Capacity,
		Stealth:     d.Prototype.Stealth,
		Speed:       d.Prototype.Speed,
		Count:       d.Count,
	}
}

// MilitaryUnitsFromDeployed returns military units for a list of deploy-ready stacks.
func MilitaryUnitsFromDeployed(items []DeploymentReadyItem) []MilitaryUnitSnap {
	if len(items) == 0 {
		return nil
	}
	out := make([]MilitaryUnitSnap, 0, len(items))
	for _, d := range items {
		out = append(out, MilitaryUnitFromDeployed(d))
	}
	return out
}

// DefenseStructuresFromBuildings returns defense structures based on present buildings with defense data.
// For now we map each defensive building to a single structure with Defence equal to DefenceBonus.
func DefenseStructuresFromBuildings(buildings []BuildItemPresent) []DefenseStructureSnap {
	if len(buildings) == 0 {
		return nil
	}
	out := make([]DefenseStructureSnap, 0, len(buildings))
	for _, b := range buildings {
		if b.Prototype.DefenseData != nil {
			out = append(out, DefenseStructureSnap{
				PrototypeID: b.Prototype.ID,
				Defence:     b.Prototype.DefenseData.DefenceBonus,
				Count:       1,
			})
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// DefenseStructureFromStack converts a defense stack into a simplified snapshot.
func DefenseStructureFromStack(stack DefenseStack) DefenseStructureSnap {
	defence := 0
	if stack.Prototype.DefenseData != nil {
		defence = stack.Prototype.DefenseData.DefenceBonus
	}
	return DefenseStructureSnap{
		PrototypeID: stack.Prototype.ID,
		Defence:     defence,
		Count:       stack.Count,
	}
}

// Helpers

func cloneUnits(src []MilitaryUnitSnap) []MilitaryUnitSnap {
	if len(src) == 0 {
		return nil
	}
	out := make([]MilitaryUnitSnap, len(src))
	copy(out, src)
	return out
}

func cloneStorageSnaps(src []StorageItemSnap) []StorageItemSnap {
	if len(src) == 0 {
		return nil
	}
	out := make([]StorageItemSnap, len(src))
	copy(out, src)
	return out
}

func cloneStructures(src []DefenseStructureSnap) []DefenseStructureSnap {
	if len(src) == 0 {
		return nil
	}
	out := make([]DefenseStructureSnap, len(src))
	copy(out, src)
	return out
}

func slowestSpeed(units []MilitaryUnitSnap) int {
	if len(units) == 0 {
		return 0
	}
	min := 0
	for _, u := range units {
		if u.Count <= 0 {
			continue
		}
		if min == 0 || u.Speed < min {
			min = u.Speed
		}
	}
	return min
}

func sumStealth(units []MilitaryUnitSnap) int {
	total := 0
	for _, u := range units {
		total += u.Stealth * u.Count
	}
	return total
}

func sumAttack(units []MilitaryUnitSnap) int {
	total := 0
	for _, u := range units {
		total += u.Attack * u.Count
	}
	return total
}

func sumDefence(units []MilitaryUnitSnap) int {
	total := 0
	for _, u := range units {
		total += u.Defence * u.Count
	}
	return total
}

func sumCapacity(units []MilitaryUnitSnap) int {
	total := 0
	for _, u := range units {
		total += u.Capacity * u.Count
	}
	return total
}

func sumStructureDefence(structs []DefenseStructureSnap) int {
	total := 0
	for _, s := range structs {
		total += s.Defence * s.Count
	}
	return total
}

// Effective stat helpers with multipliers

func SumEffectiveAttack(units []MilitaryUnitSnap, mul float64) float64 {
	return float64(sumAttack(units)) * mul
}

func SumEffectiveDefence(units []MilitaryUnitSnap, mul float64) float64 {
	return float64(sumDefence(units)) * mul
}

func SumEffectiveStealth(units []MilitaryUnitSnap, mul float64) float64 {
	return float64(sumStealth(units)) * mul
}

func SumEffectiveCapacity(units []MilitaryUnitSnap, mul float64) float64 {
	return float64(sumCapacity(units)) * mul
}

func SumEffectiveStructureDefence(structs []DefenseStructureSnap, mul float64) float64 {
	return float64(sumStructureDefence(structs)) * mul
}

func GetEffectiveSpeed(units []MilitaryUnitSnap, mul float64) int {
	base := slowestSpeed(units)
	if base <= 0 {
		return 0
	}
	return int(math.Round(float64(base) * mul))
}
