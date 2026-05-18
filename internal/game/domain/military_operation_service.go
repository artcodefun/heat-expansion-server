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

// BuildMilitaryOperationForCreation validates deployment requests against the attacker base,
// converts them to unit snapshots, and constructs a MilitaryOperation aggregate ready for
// persistence. Mirrors BuildTradeOperationForCreation in its contract.
func BuildMilitaryOperationForCreation(
	opType MilitaryOperationType,
	attacker *UserBaseModel,
	targetCoords Vector2i,
	deploymentRequests []ArmyDeploymentRequest,
) (*MilitaryOperation, error) {
	readyToDeploy, err := attacker.GetReadyToDeployArmy(deploymentRequests)
	if err != nil {
		return nil, err
	}
	snaps := attacker.ActiveStorageSnaps()
	units := MilitaryUnitsFromDeployed(readyToDeploy)

	switch opType {
	case MilitaryOperationTypeAttack:
		return NewAttackOperation(attacker.UserID, attacker.ID, attacker.Coordinates, targetCoords, units, snaps)
	case MilitaryOperationTypeSpy:
		return NewSpyOperation(attacker.UserID, attacker.ID, attacker.Coordinates, targetCoords, units, snaps)
	default:
		return nil, NewError("error.domain.operation.invalid_type", nil)
	}
}

// StartOperationAndCommitAttacker allocates the attacker's army units to the operation and starts
// outbound travel. Must be called after the operation has been persisted (so it has an ID).
// Mirrors CommitSenderForTradeCreation in its contract.
func (s MilitaryOperationService) StartOperationAndCommitAttacker() error {
	if err := allocateMilitaryArmyByUnitSnaps(s.Attacker, s.Operation.Units, s.Operation.ID); err != nil {
		return err
	}
	s.Operation.Start()
	return nil
}

// allocateMilitaryArmyByUnitSnaps allocates army units from a base to a military operation
// by matching present stacks against the unit snapshot prototype IDs.
func allocateMilitaryArmyByUnitSnaps(base *UserBaseModel, units []MilitaryUnitSnap, operationID int) error {
	for _, req := range units {
		remaining := req.Count
		for remaining > 0 {
			idx := -1
			for i, p := range base.ArmiesPresent {
				if p.Prototype.ID == req.PrototypeID && p.Count > 0 {
					idx = i
					break
				}
			}
			if idx == -1 {
				return NewError("error.domain.operation.insufficient_army_units", H{"prototype_id": req.PrototypeID, "required": req.Count, "available": req.Count - remaining})
			}

			take := remaining
			if base.ArmiesPresent[idx].Count < take {
				take = base.ArmiesPresent[idx].Count
			}

			if _, err := base.AllocateArmyToOperation(ArmyDeploymentRequest{PresentItemID: base.ArmiesPresent[idx].ID, Count: take}, OperationKindMilitary, operationID); err != nil {
				return err
			}
			remaining -= take
		}
	}
	return nil
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
		defenders := MilitaryUnitsFromPresent(defender.ArmiesPresent)
		structures := DefenseStructuresFromBuildings(defender.BuildingsPresent)
		defenderSnaps := defender.ActiveStorageSnaps()

		// Resolve using operation's domain logic
		// Compute available resource pool at target for loot calculation inside operation
		available := PriceModel{
			Credits:    int(max(0, defender.Stats.Credits)),
			Iron:       int(max(0, defender.Stats.Iron)),
			Titanium:   int(max(0, defender.Stats.Titanium)),
			Antimatter: int(max(0, defender.Stats.Antimatter)),
		}
		res := s.Operation.ResolveAttack(defenders, structures, available, nil, defenderSnaps)

		// Apply effects
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(OperationKindMilitary, s.Operation.ID, res.AttackerRemaining)
			defender.ApplyDefenderArmyRemaining(res.DefenderRemaining)
			defender.ApplyRemainingDefensiveStructures(res.RemainingStructures)
			// Deduct loot from defender resources
			defender.DeductLoot(res.Loot)
		}
	case MilitaryOperationTypeSpy:
		// Spy: no loot, but still compute survivors and adjust only spy defenders
		defenderSnaps := defender.ActiveStorageSnaps()
		allDefenders := MilitaryUnitsFromPresent(defender.ArmiesPresent)
		defendingSpies := filterSpyUnits(allDefenders)
		cloak := cloakingStrengthFromBuildings(defender.BuildingsPresent)
		res := s.Operation.ResolveSpy(cloak, defendingSpies, defenderSnaps)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(OperationKindMilitary, s.Operation.ID, res.AttackerRemaining)
			if res.Outcome == SpyOutcomeBlockedByCloaking {
				defender.ApplyDefenderArmyRemaining(allDefenders)
			} else {
				// Merge spy remaining into full defenders snapshot to avoid wiping non-spy units.
				merged := mergeSpyRemaining(allDefenders, res.DefenderRemaining)
				defender.ApplyDefenderArmyRemaining(merged)
			}
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
			Credits:    max(loc.Resources.Credits, 0),
			Iron:       max(loc.Resources.Iron, 0),
			Titanium:   max(loc.Resources.Titanium, 0),
			Antimatter: max(loc.Resources.Antimatter, 0),
		}
		res := s.Operation.ResolveAttack(defenders, structures, available, nil, nil)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(OperationKindMilitary, s.Operation.ID, res.AttackerRemaining)
			loc.ApplyDefenderArmyRemaining(res.DefenderRemaining)
			loc.ApplyRemainingDefensiveStructures(res.RemainingStructures)
			// Deduct loot from location resources
			loc.DeductLoot(res.Loot)
		}
	case MilitaryOperationTypeSpy:
		allDefenders := loc.MaterializeDefenderArmySnapshot()
		defendingSpies := filterSpyUnits(allDefenders)
		// No cloaking at resource locations for now
		res := s.Operation.ResolveSpy(0, defendingSpies, nil)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(OperationKindMilitary, s.Operation.ID, res.AttackerRemaining)
			if res.Outcome == SpyOutcomeBlockedByCloaking {
				loc.ApplyDefenderArmyRemaining(allDefenders)
			} else {
				merged := mergeSpyRemaining(allDefenders, res.DefenderRemaining)
				loc.ApplyDefenderArmyRemaining(merged)
			}
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
			Credits:    max(loc.Resources.Credits, 0),
			Iron:       max(loc.Resources.Iron, 0),
			Titanium:   max(loc.Resources.Titanium, 0),
			Antimatter: max(loc.Resources.Antimatter, 0),
		}
		res := s.Operation.ResolveAttack(defenders, structures, available, loc.Trophies, nil)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(OperationKindMilitary, s.Operation.ID, res.AttackerRemaining)
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
		res := s.Operation.ResolveSpy(0, defendingSpies, nil)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(OperationKindMilitary, s.Operation.ID, res.AttackerRemaining)
			if res.Outcome == SpyOutcomeBlockedByCloaking {
				loc.ApplyDefenderArmyRemaining(allDefenders)
			} else {
				merged := mergeSpyRemaining(allDefenders, res.DefenderRemaining)
				loc.ApplyDefenderArmyRemaining(merged)
			}
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
		res := s.Operation.ResolveAttack(defenders, structures, PriceModel{}, nil, nil)
		if res != nil {
			// available loot is zero; res.Loot remains empty
			s.Attacker.TrimDeployedToSurvivors(OperationKindMilitary, s.Operation.ID, res.AttackerRemaining)
		}
	case MilitaryOperationTypeSpy:
		// No defenders, no cloaking
		res := s.Operation.ResolveSpy(0, nil, nil)
		if res != nil {
			s.Attacker.TrimDeployedToSurvivors(OperationKindMilitary, s.Operation.ID, res.AttackerRemaining)
		}
	}
	s.Operation.StartReturn()
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
