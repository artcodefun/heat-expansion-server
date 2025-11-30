package domain

// CloneOperationUnits returns a deep copy (by value) of the given operation units slice.
// OperationUnit is a value object, so a shallow element copy is sufficient.
func CloneOperationUnits(src []MilitaryUnit) []MilitaryUnit {
	if len(src) == 0 {
		return nil
	}
	out := make([]MilitaryUnit, len(src))
	copy(out, src)
	return out
}

// CloneDefenseStructures returns a deep copy (by value) of the given defense structures slice.
// DefenseStructure is a value object, so a shallow element copy is sufficient.
func CloneDefenseStructures(src []DefenseStructure) []DefenseStructure {
	if len(src) == 0 {
		return nil
	}
	out := make([]DefenseStructure, len(src))
	copy(out, src)
	return out
}

// OperationUnitFromPresent converts a present army stack to an OperationUnit snapshot.
func OperationUnitFromPresent(p ArmyItemPresent) MilitaryUnit {
	return MilitaryUnit{
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

// OperationUnitsFromPresent maps present army stacks to operation unit snapshots.
func OperationUnitsFromPresent(items []ArmyItemPresent) []MilitaryUnit {
	if len(items) == 0 {
		return nil
	}
	out := make([]MilitaryUnit, 0, len(items))
	for _, it := range items {
		if it.Count > 0 {
			out = append(out, OperationUnitFromPresent(it))
		}
	}
	return out
}

// DefenseStructuresFromBuildings returns defense structures based on present buildings with defense data.
// For now we map each defensive building to a single structure with Defence equal to ShieldStrength.
func DefenseStructuresFromBuildings(buildings []BuildItemPresent) []DefenseStructure {
	if len(buildings) == 0 {
		return nil
	}
	out := make([]DefenseStructure, 0, len(buildings))
	for _, b := range buildings {
		if b.Prototype.DefenseData != nil {
			out = append(out, DefenseStructure{
				PrototypeID: b.Prototype.ID,
				Defence:     b.Prototype.DefenseData.ShieldStrength,
				Count:       1,
			})
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
