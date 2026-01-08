package domain

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

// Helpers

func cloneUnits(src []MilitaryUnitSnap) []MilitaryUnitSnap {
	if len(src) == 0 {
		return nil
	}
	out := make([]MilitaryUnitSnap, len(src))
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
