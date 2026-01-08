package domain

// MilitaryUnit captures a snapshot of an army unit participating in an operation.
// Using a snapshot decouples battle resolution from prototype changes over time.
type MilitaryUnit struct {
	PrototypeID int
	Category    ArmyCategory
	Attack      int
	Defence     int
	Capacity    int
	Stealth     int
	Speed       int
	Count       int
}

// DefenseStructure is a simplified snapshot of a defending structure (e.g., turrets, shields).
// This provides a lightweight input for attack resolution without coupling to build prototypes.
type DefenseStructure struct {
	PrototypeID int
	Defence     int
	Count       int
}

// Helpers

func cloneUnits(src []MilitaryUnit) []MilitaryUnit {
	if len(src) == 0 {
		return nil
	}
	out := make([]MilitaryUnit, len(src))
	copy(out, src)
	return out
}

func cloneStructures(src []DefenseStructure) []DefenseStructure {
	if len(src) == 0 {
		return nil
	}
	out := make([]DefenseStructure, len(src))
	copy(out, src)
	return out
}

func slowestSpeed(units []MilitaryUnit) int {
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

func sumStealth(units []MilitaryUnit) int {
	total := 0
	for _, u := range units {
		total += u.Stealth * u.Count
	}
	return total
}

func sumAttack(units []MilitaryUnit) int {
	total := 0
	for _, u := range units {
		total += u.Attack * u.Count
	}
	return total
}

func sumDefence(units []MilitaryUnit) int {
	total := 0
	for _, u := range units {
		total += u.Defence * u.Count
	}
	return total
}

func sumCapacity(units []MilitaryUnit) int {
	total := 0
	for _, u := range units {
		total += u.Capacity * u.Count
	}
	return total
}

func sumStructureDefence(structs []DefenseStructure) int {
	total := 0
	for _, s := range structs {
		total += s.Defence * s.Count
	}
	return total
}
