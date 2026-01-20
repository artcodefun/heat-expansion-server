package domain

// MilitaryOperationService performs battle resolution and applies effects to aggregates.
// It has no repository dependencies; callers are responsible for loading and persisting aggregates.
// The service holds the operation and attacker context to avoid threading them through each call.
type MilitaryOperationService struct {
	Operation *MilitaryOperation
	Attacker  *UserBaseModel
}

// NewMilitaryOperationService constructs a new service instance bound to a specific
// operation and attacker user base.
func NewMilitaryOperationService(op *MilitaryOperation, attacker *UserBaseModel) MilitaryOperationService {
	return MilitaryOperationService{Operation: op, Attacker: attacker}
}

// ResolveAgainstUserBase resolves an operation against a defending user base and mutates
// the operation, the attacker base (deployed survivors), and the defender base (remaining defenders/structures).
// Loot is computed internally based on defender resources and attacker capacity.
func (s MilitaryOperationService) ResolveAgainstUserBase(defender *UserBaseModel) {
	// Ensure arrival state
	s.Operation.OnArrive()

	switch s.Operation.Type {
	case MilitaryOperationTypeAttack:
		// Build defender snapshot
		mods := defender.ActiveModifiers()
		defenders := MilitaryUnitsFromPresent(defender.ArmiesPresent, mods)
		structures := DefenseStructuresFromBuildings(defender.BuildingsPresent, mods)

		// Resolve using operation's domain logic
		// Compute available resource pool at target for loot calculation inside operation
		available := PriceModel{
			Credits:    maxInt(defender.Stats.Credits, 0),
			Iron:       maxInt(defender.Stats.Iron, 0),
			Titanium:   maxInt(defender.Stats.Titanium, 0),
			Antimatter: maxInt(defender.Stats.Antimatter, 0),
		}
		res := s.Operation.ResolveAttack(defenders, structures, available, nil)

		// Apply effects
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(s.Operation.ID, res.AttackerRemaining)
			defender.ApplyDefenderArmyRemaining(res.DefenderRemaining)
			defender.ApplyRemainingDefensiveStructures(res.RemainingStructures)
			// Deduct loot from defender resources
			defender.DeductLoot(res.Loot)
		}
	case MilitaryOperationTypeSpy:
		// Spy: no loot, but still compute survivors and adjust only spy defenders
		mods := defender.ActiveModifiers()
		allDefenders := MilitaryUnitsFromPresent(defender.ArmiesPresent, mods)
		defendingSpies := filterSpyUnits(allDefenders)
		cloak := cloakingStrengthFromBuildings(defender.BuildingsPresent)
		res := s.Operation.ResolveSpy(cloak, defendingSpies)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(s.Operation.ID, res.AttackerRemaining)
			// Merge spy remaining into full defenders snapshot to avoid wiping non-spy units
			merged := mergeSpyRemaining(allDefenders, res.DefenderRemaining)
			defender.ApplyDefenderArmyRemaining(merged)
		}
	default:
		// Fallback: do nothing special
	}

	// Start return leg (time computed internally)
	s.Operation.StartReturn()
}

// ResolveAgainstResourceLocation resolves an operation against a resource location and mutates
// the operation, the attacker base (deployed survivors), and the location defenders.
// Loot is computed internally based on location resources and attacker capacity.
func (s MilitaryOperationService) ResolveAgainstResourceLocation(loc *ResourceLocationModel) {
	s.Operation.OnArrive()

	switch s.Operation.Type {
	case MilitaryOperationTypeAttack:
		defenders := loc.MaterializeDefenderArmySnapshot()
		structures := loc.MaterializeDefenderStructureSnapshot()
		available := PriceModel{
			Credits:    maxInt(loc.Resources.Credits, 0),
			Iron:       maxInt(loc.Resources.Iron, 0),
			Titanium:   maxInt(loc.Resources.Titanium, 0),
			Antimatter: maxInt(loc.Resources.Antimatter, 0),
		}
		res := s.Operation.ResolveAttack(defenders, structures, available, nil)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(s.Operation.ID, res.AttackerRemaining)
			loc.ApplyDefenderArmyRemaining(res.DefenderRemaining)
			loc.ApplyRemainingDefensiveStructures(res.RemainingStructures)
			// Deduct loot from location resources
			loc.DeductLoot(res.Loot)
		}
	case MilitaryOperationTypeSpy:
		allDefenders := loc.MaterializeDefenderArmySnapshot()
		defendingSpies := filterSpyUnits(allDefenders)
		// No cloaking at resource locations for now
		res := s.Operation.ResolveSpy(0, defendingSpies)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(s.Operation.ID, res.AttackerRemaining)
			merged := mergeSpyRemaining(allDefenders, res.DefenderRemaining)
			loc.ApplyDefenderArmyRemaining(merged)
		}
	}
	s.Operation.StartReturn()
}

// ResolveAgainstDangerousLocation resolves an operation against a dangerous location and mutates
// the operation, the attacker base (deployed survivors), and the location defenders.
// Loot is computed internally based on location resources and attacker capacity.
func (s MilitaryOperationService) ResolveAgainstDangerousLocation(loc *DangerousLocationModel) {
	s.Operation.OnArrive()

	switch s.Operation.Type {
	case MilitaryOperationTypeAttack:
		defenders := loc.MaterializeDefenderArmySnapshot()
		structures := loc.MaterializeDefenderStructureSnapshot()
		available := PriceModel{
			Credits:    maxInt(loc.Resources.Credits, 0),
			Iron:       maxInt(loc.Resources.Iron, 0),
			Titanium:   maxInt(loc.Resources.Titanium, 0),
			Antimatter: maxInt(loc.Resources.Antimatter, 0),
		}
		res := s.Operation.ResolveAttack(defenders, structures, available, loc.Trophies)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(s.Operation.ID, res.AttackerRemaining)
			loc.ApplyDefenderArmyRemaining(res.DefenderRemaining)
			loc.ApplyRemainingDefensiveStructures(res.RemainingStructures)

			// Clear trophies if they were taken
			if len(res.Trophies) > 0 {
				loc.Trophies = nil
			}
			// Deduct loot from location resources (emits Drained if empty)
			loc.DeductLoot(res.Loot)
		}
	case MilitaryOperationTypeSpy:
		allDefenders := loc.MaterializeDefenderArmySnapshot()
		defendingSpies := filterSpyUnits(allDefenders)
		res := s.Operation.ResolveSpy(0, defendingSpies)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(s.Operation.ID, res.AttackerRemaining)
			merged := mergeSpyRemaining(allDefenders, res.DefenderRemaining)
			loc.ApplyDefenderArmyRemaining(merged)
		}
	}
	s.Operation.StartReturn()
}

// ResolveAgainstEmptyLocation resolves an operation against an empty sector (no defenders, no resources).
// It trims deployed to survivors (which will be all units in this placeholder resolution) and starts return.
func (s MilitaryOperationService) ResolveAgainstEmptySector(sector *SectorModel) {
	s.Operation.OnArrive()

	switch s.Operation.Type {
	case MilitaryOperationTypeAttack:
		// No defenders or structures, zero available loot
		var defenders []MilitaryUnitSnap
		var structures []DefenseStructureSnap
		res := s.Operation.ResolveAttack(defenders, structures, PriceModel{}, nil)
		if res != nil {
			// available loot is zero; res.Loot remains empty
			s.Attacker.TrimDeployedToSurvivors(s.Operation.ID, res.AttackerRemaining)
		}
	case MilitaryOperationTypeSpy:
		// No defenders, no cloaking
		res := s.Operation.ResolveSpy(0, nil)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(s.Operation.ID, res.AttackerRemaining)
		}
	}
	s.Operation.StartReturn()
}

// --- internal helpers (value-object helpers retained: capacity + loot capping) ---

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// --- local helpers for spy resolution ---

// filterSpyUnits returns only units with Spy category.
func filterSpyUnits(units []MilitaryUnitSnap) []MilitaryUnitSnap {
	if len(units) == 0 {
		return nil
	}
	out := make([]MilitaryUnitSnap, 0, len(units))
	for _, u := range units {
		if u.Category == ArmyCategorySpy && u.Count > 0 {
			out = append(out, u)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

// mergeSpyRemaining merges spy remaining snapshot into the original defenders snapshot,
// preserving non-spy units unchanged and applying new counts for spy prototypes.
func mergeSpyRemaining(original []MilitaryUnitSnap, spyRemaining []MilitaryUnitSnap) []MilitaryUnitSnap {
	if len(original) == 0 && len(spyRemaining) == 0 {
		return nil
	}
	merged := make([]MilitaryUnitSnap, 0, len(original)+len(spyRemaining))
	// keep non-spy units as-is
	for _, u := range original {
		if u.Category != ArmyCategorySpy && u.Count > 0 {
			merged = append(merged, u)
		}
	}
	// append remaining spies directly (they already carry proper metadata)
	for _, s := range spyRemaining {
		if s.Count > 0 {
			merged = append(merged, s)
		}
	}
	if len(merged) == 0 {
		return nil
	}
	return merged
}

// cloakingStrengthFromBuildings sums StealthStrength from intelligence buildings with CLOAKING subtype.
func cloakingStrengthFromBuildings(buildings []BuildItemPresent) int {
	if len(buildings) == 0 {
		return 0
	}
	total := 0
	for _, b := range buildings {
		if b.Prototype.IntelligenceData != nil && b.Prototype.IntelligenceData.Subtype == IntelligenceSubtypeCloaking {
			total += b.Prototype.IntelligenceData.StealthStrength
		}
	}
	return total
}
