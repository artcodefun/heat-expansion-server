package domain

import (
	"math"

	"github.com/google/uuid"
)

// MilitaryOperationType represents the type of a military operation.
type MilitaryOperationType string

const (
	MilitaryOperationTypeAttack     MilitaryOperationType = "ATTACK"
	MilitaryOperationTypeSpy        MilitaryOperationType = "SPY"
	MilitaryOperationTypeOccupation MilitaryOperationType = "OCCUPATION" // reserved for future use
)

// MilitaryOperationPhase represents the movement/lifecycle stage of an operation.
type MilitaryOperationPhase string

const (
	OperationPhasePending   MilitaryOperationPhase = "PENDING"   // created, not yet departed
	OperationPhaseOutbound  MilitaryOperationPhase = "OUTBOUND"  // traveling to target
	OperationPhaseAtTarget  MilitaryOperationPhase = "AT_TARGET" // arrived at target
	OperationPhaseResolving MilitaryOperationPhase = "RESOLVING" // resolving at target
	OperationPhaseReturning MilitaryOperationPhase = "RETURNING" // traveling back to source
	OperationPhaseCompleted MilitaryOperationPhase = "COMPLETED" // returned and finished
)

// MilitaryOperationResult represents the outcome once completed/canceled.
type MilitaryOperationResult string

const (
	OperationResultUnknown  MilitaryOperationResult = "UNKNOWN"
	OperationResultSuccess  MilitaryOperationResult = "SUCCESS"
	OperationResultFailure  MilitaryOperationResult = "FAILURE"
	OperationResultCanceled MilitaryOperationResult = "CANCELED"
)

// SpyOutcome enumerates possible results of a spy operation.
type SpyOutcome string

const (
	SpyOutcomeBlockedByCloaking SpyOutcome = "BLOCKED_BY_CLOAKING_EMPTY_REPORT" // target cloaking >= attackers stealth; units unharmed
	SpyOutcomeDefeatedBySpies   SpyOutcome = "DEFEATED_BY_DEFENDING_SPIES"      // attackers lost skirmish with defending spies
	SpyOutcomeReportProduced    SpyOutcome = "REPORT_PRODUCED"                  // successful intel report created
)

// AttackOutcome enumerates possible results of an attack operation.
type AttackOutcome string

const (
	AttackOutcomeAttackerWon  AttackOutcome = "ATTACKER_WON"
	AttackOutcomeDefenderHeld AttackOutcome = "DEFENDER_HELD"
)

// Spy now uses SectorScanReport as the intelligence artifact.

type SpyResult struct {
	Outcome           SpyOutcome
	AttackerRemaining []MilitaryUnitSnap
	DefenderRemaining []MilitaryUnitSnap
	// New: snapshot of defenders before resolution (for UI diffs)
	DefendersBefore []MilitaryUnitSnap
	// Buffs/artifacts active during resolution
	DefenderStorageSnaps []StorageItemSnap
	// Aggregated multiplier for defenders
	TotalDefenderModifiers MilitaryModifiers
}

type AttackResult struct {
	Outcome             AttackOutcome
	AttackerRemaining   []MilitaryUnitSnap
	DefenderRemaining   []MilitaryUnitSnap
	RemainingStructures []DefenseStructureSnap
	Loot                PriceModel          // what attackers managed to carry back; computed elsewhere
	Trophies            []TrophyStorageItem // special items collected from dangerous locations
	// New: snapshots for UI to show casualties/damage
	DefendersBefore  []MilitaryUnitSnap
	StructuresBefore []DefenseStructureSnap
	// Buffs/artifacts active during resolution
	DefenderStorageSnaps []StorageItemSnap
	// Aggregated multiplier for defenders
	TotalDefenderModifiers MilitaryModifiers
}

// MilitaryOperation models an attack or spy op traveling between sectors and resolving on arrival.
type MilitaryOperation struct {
	EventProducer
	ID           int
	UUID         uuid.UUID
	Type         MilitaryOperationType
	OwnerUserID  uuid.UUID
	SourceBaseID int

	// Coordinates snapshot (for travel calculations)
	SourceCoordinates Vector2i
	TargetCoordinates Vector2i

	// Timeline
	OutboundDepartAt int64
	OutboundArriveAt int64
	ReturnDepartAt   int64
	ReturnArriveAt   int64
	CompletedAt      int64

	// CrystalsSkipPrice is the base crystal cost to skip the current travel leg.
	// It is recomputed whenever the operation starts outbound or return travel.
	CrystalsSkipPrice int

	Phase  MilitaryOperationPhase
	Result MilitaryOperationResult

	// Snapshot of attacking units
	Units []MilitaryUnitSnap

	// Snapshot of attacking storage items (buffs/artifacts) active at start
	StorageSnaps []StorageItemSnap

	// TotalModifiers is the aggregated multiplier from StorageSnaps.
	TotalModifiers MilitaryModifiers

	// Results (only one will be populated depending on Type)
	SpyResult    *SpyResult
	AttackResult *AttackResult
}

// NewAttackOperation creates an ATTACK operation in transit.
// It validates that at least one unit is provided and that source/target are different.
func NewAttackOperation(ownerUserID uuid.UUID, sourceBaseID int, source, target Vector2i, units []MilitaryUnitSnap, storageSnaps []StorageItemSnap) (*MilitaryOperation, error) {
	if source == target {
		return nil, NewError("error.domain.operation.source_equals_target", nil)
	}
	if len(units) == 0 {
		return nil, NewError("error.domain.operation.no_units_provided", nil)
	}
	op := &MilitaryOperation{
		UUID:              uuid.Must(uuid.NewV7()),
		Type:              MilitaryOperationTypeAttack,
		OwnerUserID:       ownerUserID,
		SourceBaseID:      sourceBaseID,
		SourceCoordinates: source,
		TargetCoordinates: target,
		OutboundDepartAt:  0,
		OutboundArriveAt:  0,
		Phase:             OperationPhasePending,
		Result:            OperationResultUnknown,
		Units:             cloneUnits(units),
		StorageSnaps:      cloneStorageSnaps(storageSnaps),
		TotalModifiers:    MilitaryModifiersFromSnaps(storageSnaps),
	}
	return op, nil
}

// NewSpyOperation creates a SPY operation in transit.
// It validates that at least one unit is provided, targeting a different sector, and that all units are spies.
func NewSpyOperation(ownerUserID uuid.UUID, sourceBaseID int, source, target Vector2i, spies []MilitaryUnitSnap, storageSnaps []StorageItemSnap) (*MilitaryOperation, error) {
	if source == target {
		return nil, NewError("error.domain.operation.invalid_coordinates", nil)
	}
	if len(spies) == 0 {
		return nil, NewError("error.domain.operation.no_spy_units", nil)
	}
	for _, u := range spies {
		if u.Category != ArmyCategorySpy {
			return nil, NewError("error.domain.operation.only_spy_units_allowed", nil)
		}
	}
	op := &MilitaryOperation{
		UUID:              uuid.Must(uuid.NewV7()),
		Type:              MilitaryOperationTypeSpy,
		OwnerUserID:       ownerUserID,
		SourceBaseID:      sourceBaseID,
		SourceCoordinates: source,
		TargetCoordinates: target,
		OutboundDepartAt:  0,
		OutboundArriveAt:  0,
		Phase:             OperationPhasePending,
		Result:            OperationResultUnknown,
		Units:             cloneUnits(spies),
		StorageSnaps:      cloneStorageSnaps(storageSnaps),
		TotalModifiers:    MilitaryModifiersFromSnaps(storageSnaps),
	}
	return op, nil
}

// Start begins the operation's outbound travel if it is currently pending.
// It sets departure/arrival timestamps, switches phase to OUTBOUND, and emits a Started event.
func (op *MilitaryOperation) Start() {
	if op.Phase != OperationPhasePending {
		return
	}
	now := NowUnix()
	travelSeconds := computeTravelSecondsBetween(op.SourceCoordinates, op.TargetCoordinates, op.Units, op.TotalModifiers)
	op.OutboundDepartAt = now
	op.OutboundArriveAt = now + travelSeconds
	// Base skip price proportional to total travel time (minimum 1 crystal)
	op.CrystalsSkipPrice = max(1, int(travelSeconds/60))
	op.Phase = OperationPhaseOutbound
	op.AddEvent(NewMilitaryOperationStartedEvent(op.ID, op.OutboundArriveAt))
}

// OnArrive marks the operation as arrived if it was in transit.
func (op *MilitaryOperation) OnArrive() {
	if op.Phase != OperationPhaseOutbound {
		return
	}
	now := NowUnix()
	op.Phase = OperationPhaseAtTarget
	op.OutboundArriveAt = now
	op.AddEvent(NewMilitaryOperationArrivedEvent(op.ID))
}

// ResolveSpy resolves a spy operation by first checking target cloaking vs attacker stealth,
// then (if not blocked) resolving a skirmish versus defending spies.
// targetCloakingStrength: sum of StealthStrength across target cloaking buildings
// defendingSpies: snapshot of defending spy units present in the sector
func (op *MilitaryOperation) ResolveSpy(targetCloakingStrength int, defendingSpies []MilitaryUnitSnap, defenderStorageSnaps []StorageItemSnap) *SpyResult {
	if op.Type != MilitaryOperationTypeSpy || (op.Phase != OperationPhaseAtTarget && op.Phase != OperationPhaseOutbound) {
		return op.SpyResult
	}
	op.Phase = OperationPhaseResolving

	atkMods := op.TotalModifiers
	defMods := MilitaryModifiersFromSnaps(defenderStorageSnaps)

	// Ensure we consider only spy-category units on the attacker side
	attackers := filterUnitsByCategory(op.Units, ArmyCategorySpy)
	attackerStealth := SumEffectiveStealth(attackers, atkMods.StealthMul)

	// Snapshot defenders before resolving for UI diffs
	defBefore := cloneUnits(defendingSpies)

	// 1) Cloaking check: if target cloaking >= attacker stealth, the report is empty and no defender intel is exposed.
	if float64(targetCloakingStrength) >= attackerStealth {
		res := &SpyResult{
			Outcome: SpyOutcomeBlockedByCloaking,
			// Everyone survives in this outcome; no defender intel should be exposed.
			AttackerRemaining: cloneUnits(attackers),
		}
		op.SpyResult = res
		op.Result = OperationResultSuccess
		op.AddEvent(NewMilitaryOperationResolvedEvent(op.ID, op.Result))
		return res
	}

	// 2) Skirmish with defending spies
	atkRemaining, defRemaining, attackerWon := resolveSpySkirmish(attackers, atkMods, defendingSpies, defMods)
	if !attackerWon {
		res := &SpyResult{
			Outcome:                SpyOutcomeDefeatedBySpies,
			AttackerRemaining:      atkRemaining,
			DefenderRemaining:      defRemaining,
			DefendersBefore:        defBefore,
			DefenderStorageSnaps:   cloneStorageSnaps(defenderStorageSnaps),
			TotalDefenderModifiers: defMods,
		}
		op.SpyResult = res
		op.Result = OperationResultFailure
		op.AddEvent(NewMilitaryOperationResolvedEvent(op.ID, op.Result))
		return res
	}

	// 3) Successful intelligence report
	res := &SpyResult{
		Outcome:                SpyOutcomeReportProduced,
		AttackerRemaining:      atkRemaining,
		DefenderRemaining:      defRemaining,
		DefendersBefore:        defBefore,
		DefenderStorageSnaps:   cloneStorageSnaps(defenderStorageSnaps),
		TotalDefenderModifiers: defMods,
	}
	op.SpyResult = res
	op.Result = OperationResultSuccess
	op.AddEvent(NewMilitaryOperationResolvedEvent(op.ID, op.Result))
	return res
}

// ResolveAttack resolves an attack using a simplified power comparison as a placeholder.
// A more detailed sequential algorithm can replace this later.
func (op *MilitaryOperation) ResolveAttack(defenders []MilitaryUnitSnap, structures []DefenseStructureSnap, availableResourcePool PriceModel, trophies []TrophyStorageItem, defenderStorageSnaps []StorageItemSnap) *AttackResult {
	if op.Type != MilitaryOperationTypeAttack || (op.Phase != OperationPhaseAtTarget && op.Phase != OperationPhaseOutbound) {
		return op.AttackResult
	}
	op.Phase = OperationPhaseResolving

	atkMods := op.TotalModifiers
	defMods := MilitaryModifiersFromSnaps(defenderStorageSnaps)

	// capture "before" snapshots for UI diffs
	defBefore := cloneUnits(defenders)
	structBefore := cloneStructures(structures)

	atkRemain, defRemain, structRemain, attackerWon := resolveAttackCombat(cloneUnits(op.Units), atkMods, cloneUnits(defenders), defMods, cloneStructures(structures))

	result := &AttackResult{
		DefendersBefore:        defBefore,
		StructuresBefore:       structBefore,
		AttackerRemaining:      atkRemain,
		DefenderRemaining:      defRemain,
		RemainingStructures:    structRemain,
		DefenderStorageSnaps:   cloneStorageSnaps(defenderStorageSnaps),
		TotalDefenderModifiers: defMods,
		// Compute loot here based on remaining attackers' capacity and available resources at target
		Loot: computeLoadFromLocation(atkRemain, atkMods, availableResourcePool),
	}
	if attackerWon {
		result.Outcome = AttackOutcomeAttackerWon
		op.Result = OperationResultSuccess
		// Award trophies only if attacker won
		result.Trophies = trophies
	} else {
		result.Outcome = AttackOutcomeDefenderHeld
		op.Result = OperationResultFailure
	}

	op.AttackResult = result
	op.AddEvent(NewMilitaryOperationResolvedEvent(op.ID, op.Result))
	return result
}

// resolveAttackCombat performs a simple round-based combat simulation between attackers and defenders.
// Each round:
//   - Attackers deal total damage sum(Attack*Count) to defenders first, spilling over to structures
//   - Defenders deal total damage [sum(Defence*Count) of defending units + sum(Defence*Count) of structures] to attackers
//
// HP model:
//   - Attacking units use their Attack as per-instance HP when taking damage
//   - Defending units use their Defence as per-instance HP when taking damage
//   - Defensive structures use their Defence as per-structure HP
//
// Damage removes full instances proportionally to float HP.
// Returns remaining stacks and whether attackers won (defenders and structures eliminated).
func resolveAttackCombat(attackers []MilitaryUnitSnap, atkMods MilitaryModifiers, defenders []MilitaryUnitSnap, defMods MilitaryModifiers, structures []DefenseStructureSnap) (atkRemaining []MilitaryUnitSnap, defRemaining []MilitaryUnitSnap, structRemaining []DefenseStructureSnap, attackerWon bool) {
	// Normalize: drop zero-count stacks
	attackers = filterZeroCountUnits(attackers)
	defenders = filterZeroCountUnits(defenders)
	structures = filterZeroCountStructures(structures)

	// Safety cap to avoid infinite loops in degenerate configs
	const maxRounds = 1000
	for round := 0; round < maxRounds; round++ {
		// Check end conditions
		if totalCount(attackers) == 0 {
			attackerWon = false
			break
		}
		if totalCount(defenders) == 0 && totalStructCount(structures) == 0 {
			attackerWon = true
			break
		}

		// Compute damages from current forces using floats
		attDmg := SumEffectiveAttack(attackers, atkMods.AttackMul)
		// Defenders and structures both contribute damage using their Defence as power
		defDmg := SumEffectiveDefence(defenders, defMods.DefenceMul) + SumEffectiveStructureDefence(structures, defMods.DefenceMul)

		// If neither side can deal damage, stop to avoid infinite loop
		if attDmg <= 0 && defDmg <= 0 {
			attackerWon = false
			break
		}

		// Apply attacker damage to defenders then structures
		var attOverflow float64
		defenders, attOverflow = applyDamageToUnits(defenders, attDmg, true, defMods.DefenceMul)
		structures, _ = applyDamageToStructures(structures, attOverflow, defMods.DefenceMul)

		// Apply defender damage to attackers
		attackers, _ = applyDamageToUnits(attackers, defDmg, false, atkMods.AttackMul)

		// Optional normalization could be applied here if needed
	}

	return attackers, defenders, structures, attackerWon
}

// applyDamageToUnits applies damage to units. When isDefensive=true, per-instance HP is Defence; otherwise HP is Attack.
// Damage is float64; per-instance HP is also scaled by mul.
// Returns updated units and any leftover damage not consumed.
func applyDamageToUnits(units []MilitaryUnitSnap, damage float64, isDefensive bool, mul float64) ([]MilitaryUnitSnap, float64) {
	if damage <= 0 || len(units) == 0 {
		return units, damage
	}
	out := make([]MilitaryUnitSnap, 0, len(units))
	remainingDamage := damage
	for _, u := range units {
		if u.Count <= 0 {
			continue
		}
		var baseHp int
		if isDefensive {
			baseHp = u.Defence
		} else {
			baseHp = u.Attack
		}
		hp := float64(baseHp) * mul
		// Guard against non-positive HP creating immortal stacks
		if hp <= 0 {
			hp = 1.0
		}
		// kills = min(count, floor(damage / hpPerUnit))
		possibleKills := int(math.Floor(remainingDamage / hp))
		if possibleKills > 0 {
			if possibleKills >= u.Count {
				remainingDamage -= float64(u.Count) * hp
				// all killed -> skip append
				continue
			}
			// partial kills
			u.Count -= possibleKills
			remainingDamage -= float64(possibleKills) * hp
		}
		if u.Count > 0 {
			out = append(out, u)
		}
	}
	return out, remainingDamage
}

// (No separate variants; isDefensive flag selects HP stat.)

// applyDamageToStructures applies damage to structures similarly based on Defence HP per structure.
func applyDamageToStructures(structs []DefenseStructureSnap, damage float64, mul float64) ([]DefenseStructureSnap, float64) {
	if damage <= 0 || len(structs) == 0 {
		return structs, damage
	}
	out := make([]DefenseStructureSnap, 0, len(structs))
	remainingDamage := damage
	for _, s := range structs {
		if s.Count <= 0 {
			continue
		}
		hp := float64(s.Defence) * mul
		if hp <= 0 {
			hp = 1.0
		}
		possibleKills := int(math.Floor(remainingDamage / hp))
		if possibleKills > 0 {
			if possibleKills >= s.Count {
				remainingDamage -= float64(s.Count) * hp
				continue
			}
			s.Count -= possibleKills
			remainingDamage -= float64(possibleKills) * hp
		}
		if s.Count > 0 {
			out = append(out, s)
		}
	}
	return out, remainingDamage
}

// Helpers to filter and count
func filterZeroCountUnits(units []MilitaryUnitSnap) []MilitaryUnitSnap {
	if len(units) == 0 {
		return nil
	}
	out := make([]MilitaryUnitSnap, 0, len(units))
	for _, u := range units {
		if u.Count > 0 {
			out = append(out, u)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func filterZeroCountStructures(structs []DefenseStructureSnap) []DefenseStructureSnap {
	if len(structs) == 0 {
		return nil
	}
	out := make([]DefenseStructureSnap, 0, len(structs))
	for _, s := range structs {
		if s.Count > 0 {
			out = append(out, s)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func totalCount(units []MilitaryUnitSnap) int {
	if len(units) == 0 {
		return 0
	}
	sum := 0
	for _, u := range units {
		sum += u.Count
	}
	return sum
}

func totalStructCount(structs []DefenseStructureSnap) int {
	if len(structs) == 0 {
		return 0
	}
	sum := 0
	for _, s := range structs {
		sum += s.Count
	}
	return sum
}

// StartReturn initiates the return leg after resolution (or cancel).
func (op *MilitaryOperation) StartReturn() {
	if op.Phase != OperationPhaseResolving && op.Phase != OperationPhaseAtTarget && op.Phase != OperationPhaseOutbound {
		return
	}
	now := NowUnix()
	// Determine survivors for the return leg
	var returningUnits []MilitaryUnitSnap
	if op.AttackResult != nil {
		returningUnits = op.AttackResult.AttackerRemaining
	} else if op.SpyResult != nil {
		returningUnits = op.SpyResult.AttackerRemaining
	} else {
		returningUnits = op.Units
	}

	// If no survivors remain, complete immediately (no return travel)
	survivors := 0
	for _, u := range returningUnits {
		if u.Count > 0 {
			survivors += u.Count
		}
	}
	if survivors == 0 {
		op.ReturnDepartAt = now
		op.ReturnArriveAt = now
		op.CompletedAt = now
		op.Phase = OperationPhaseCompleted
		// Emit arrival event immediately to notify listeners of completion
		op.AddEvent(NewMilitaryOperationReturnArrivedEvent(op.ID))
		return
	}

	travelSeconds := computeTravelSecondsBetween(op.TargetCoordinates, op.SourceCoordinates, returningUnits, op.TotalModifiers)
	op.ReturnDepartAt = now
	op.ReturnArriveAt = now + travelSeconds
	// Recompute skip price for the return leg (minimum 1 crystal)
	op.CrystalsSkipPrice = max(1, int(travelSeconds/60))
	op.Phase = OperationPhaseReturning
	op.AddEvent(NewMilitaryOperationReturnStartedEvent(op.ID, op.ReturnArriveAt))
}

// OnReturnArrive completes the operation on reaching the source base.
func (op *MilitaryOperation) OnReturnArrive() {
	if op.Phase != OperationPhaseReturning {
		return
	}
	now := NowUnix()
	op.ReturnArriveAt = now
	op.CompletedAt = now
	op.Phase = OperationPhaseCompleted
	op.AddEvent(NewMilitaryOperationReturnArrivedEvent(op.ID))
}

// UpdatePhaseBasedOnTime advances the operation's phase based on current time and scheduled timestamps.
// It is idempotent and safe to call multiple times.
func (op *MilitaryOperation) UpdatePhaseBasedOnTime() {
	now := NowUnix()
	switch op.Phase {
	case OperationPhaseOutbound:
		if op.OutboundArriveAt > 0 && now >= op.OutboundArriveAt {
			op.OnArrive()
		}
	case OperationPhaseReturning:
		if op.ReturnArriveAt > 0 && now >= op.ReturnArriveAt {
			op.OnReturnArrive()
		}
	default:
		// no-op for other phases
	}
}

// Cancel marks the operation canceled if it is still in transit (Pending or Outbound).
// Once it reaches the target or starts resolving, it can no longer be canceled.
func (op *MilitaryOperation) Cancel() error {
	if op.Phase != OperationPhasePending && op.Phase != OperationPhaseOutbound {
		return NewError("error.domain.operation.cannot_cancel_in_phase", H{"phase": string(op.Phase)})
	}

	op.Result = OperationResultCanceled
	op.AddEvent(NewMilitaryOperationCancelledEvent(op.ID))

	switch op.Phase {
	case OperationPhasePending:
		// If it hasn't even started moving, it just finishes immediately.
		op.Phase = OperationPhaseCompleted
		op.CompletedAt = NowUnix()
		// This event handles returning units to the base in the command handler.
		op.AddEvent(NewMilitaryOperationReturnArrivedEvent(op.ID))
	case OperationPhaseOutbound:
		op.StartReturn()
	}

	return nil
}

func (op *MilitaryOperation) CurrentCoordinates() Vector2i {
	at := NowUnix()
	switch op.Phase {
	case OperationPhasePending:
		return op.SourceCoordinates
	case OperationPhaseOutbound:
		return lerpCoordinates(op.SourceCoordinates, op.TargetCoordinates, op.OutboundDepartAt, op.OutboundArriveAt, at)
	case OperationPhaseReturning:
		return lerpCoordinates(op.TargetCoordinates, op.SourceCoordinates, op.ReturnDepartAt, op.ReturnArriveAt, at)
	case OperationPhaseAtTarget, OperationPhaseResolving:
		return op.TargetCoordinates
	case OperationPhaseCompleted:
		return op.SourceCoordinates
	default:
		return op.SourceCoordinates
	}
}

func lerpCoordinates(s, t Vector2i, startT, endT, at int64) Vector2i {
	if at <= startT {
		return s
	}
	if at >= endT {
		return t
	}
	duration := endT - startT
	if duration <= 0 {
		return t
	}
	progress := float64(at-startT) / float64(duration)
	return Vector2i{
		X: int(math.Round(float64(s.X) + float64(t.X-s.X)*progress)),
		Y: int(math.Round(float64(s.Y) + float64(t.Y-s.Y)*progress)),
	}
}

// Travel time helpers

// computeTravelSecondsBetween computes travel time as ceil(distance / slowestSpeed).
// Distance is Manhattan on sector coordinates for now; speed is taken from the slowest unit.
func computeTravelSecondsBetween(a, b Vector2i, units []MilitaryUnitSnap, mods MilitaryModifiers) int64 {
	// Euclidean distance scaled by 1000, then divide by slowest speed and ceil
	dist := euclideanScaled(a, b)
	if dist <= 0 {
		return 0
	}
	s := GetEffectiveSpeed(units, mods.SpeedMul)
	if s <= 0 {
		s = 1 // guard against invalid data
	}
	// ceil(dist / s)
	return int64((dist + s - 1) / s)
}

// euclideanScaled returns ceil( sqrt(dx^2+dy^2) * 1000 ) as an integer distance unit
func euclideanScaled(a, b Vector2i) int {
	dx := a.X - b.X
	dy := a.Y - b.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}
	// avoid floats by using integer sqrt approximation if desired, but here we use float64 for clarity
	fdx := float64(dx)
	fdy := float64(dy)
	d := math.Sqrt(fdx*fdx+fdy*fdy) * 1000.0
	return int(math.Ceil(d))
}

// TimeBeforeEntersCircle returns the unix timestamp when the operation first enters the circle
// defined by center and radius (in coordinate units) during its outbound journey.
// Returns an error if the operation is not outbound or if it never enters the circle.
func (op *MilitaryOperation) TimeBeforeEntersCircle(center Vector2i, radius int) (int64, error) {
	if op.Phase != OperationPhaseOutbound {
		return 0, NewError("error.domain.operation.not_in_outbound", nil)
	}

	// Quadratic intersection for P(k) = S + k(T-S), k in [0, 1]
	// ||P(k) - C||^2 = R^2
	sx, sy := float64(op.SourceCoordinates.X), float64(op.SourceCoordinates.Y)
	tx, ty := float64(op.TargetCoordinates.X), float64(op.TargetCoordinates.Y)
	cx, cy := float64(center.X), float64(center.Y)
	r := float64(radius)

	dx := tx - sx
	dy := ty - sy
	wx := sx - cx
	wy := sy - cy

	a := dx*dx + dy*dy
	if a == 0 {
		return 0, NewError("error.domain.operation.no_movement", nil)
	}
	b := 2 * (dx*wx + dy*wy)
	c := wx*wx + wy*wy - r*r

	// If already inside (c <= 0), it entered at the start of travel
	if c <= 0 {
		return op.OutboundDepartAt, nil
	}

	delta := b*b - 4*a*c
	if delta < 0 {
		return 0, NewError("error.domain.operation.radar_never_enters_circle", nil)
	}

	sqrtDelta := math.Sqrt(delta)
	k1 := (-b - sqrtDelta) / (2 * a)

	if k1 < 0 {
		// Moving away from the circle since we started outside (c > 0)
		return 0, NewError("error.domain.operation.radar_moving_away", nil)
	}
	if k1 > 1.0 {
		return 0, NewError("error.domain.operation.radar_reached_before_entering", nil)
	}

	travelDuration := op.OutboundArriveAt - op.OutboundDepartAt
	enterAt := op.OutboundDepartAt + int64(float64(travelDuration)*k1)
	return enterAt, nil
}

func filterUnitsByCategory(units []MilitaryUnitSnap, cat ArmyCategory) []MilitaryUnitSnap {
	if len(units) == 0 {
		return nil
	}
	out := make([]MilitaryUnitSnap, 0, len(units))
	for _, u := range units {
		if u.Category == cat && u.Count > 0 {
			out = append(out, u)
		}
	}
	return out
}

func (op *MilitaryOperation) TotalAttack() int {
	mods := op.TotalModifiers
	return int(math.Round(SumEffectiveAttack(op.Units, mods.AttackMul)))
}

func (op *MilitaryOperation) TotalStealth() int {
	mods := op.TotalModifiers
	return int(math.Round(SumEffectiveStealth(op.Units, mods.StealthMul)))
}

func (op *MilitaryOperation) TotalSpeed() int {
	mods := op.TotalModifiers
	return GetEffectiveSpeed(op.Units, mods.SpeedMul)
}

func (op *MilitaryOperation) TotalCapacity() int {
	mods := op.TotalModifiers
	return int(math.Round(SumEffectiveCapacity(op.Units, mods.CapacityMul)))
}

func (op *MilitaryOperation) ProducedVisibleIntel() bool {
	if op.Type != MilitaryOperationTypeSpy {
		return true
	}
	return op.SpyResult == nil || op.SpyResult.Outcome != SpyOutcomeBlockedByCloaking
}

// computeLoadFromLocation fills available carrying capacity using the available
// resource pool, prioritizing lower value resources first (least expensive to most expensive).
// This incentivizes players to fully loot a location.
// Carrying capacity is volume-based (WorthCapacityMultiplier * total unit capacity).
func computeLoadFromLocation(remaining []MilitaryUnitSnap, mods MilitaryModifiers, available PriceModel) PriceModel {
	capacity := int(math.Round(SumEffectiveCapacity(remaining, mods.CapacityMul)))
	if capacity <= 0 {
		return PriceModel{}
	}

	remainingVolume := float64(capacity) * WorthCapacityMultiplier
	loot := PriceModel{}

	// Use local copies since available is passed by value anyway
	poolCredits := max(0, available.Credits)
	poolIron := max(0, available.Iron)
	poolTitanium := max(0, available.Titanium)
	poolAntimatter := max(0, available.Antimatter)

	// Take resources in order of increasing value density (Low to High):
	// Credits (1 worth) -> Iron (4 worth) -> Titanium (20 worth) -> Antimatter (333.3 worth)

	// 1. Credits
	if poolCredits > 0 && remainingVolume >= WorthCredit {
		maxAmt := int(remainingVolume / WorthCredit)
		take := min(poolCredits, maxAmt)
		loot.Credits = take
		remainingVolume -= float64(take) * WorthCredit
	}

	// 2. Iron
	if poolIron > 0 && remainingVolume >= WorthIron {
		maxAmt := int(remainingVolume / WorthIron)
		take := min(poolIron, maxAmt)
		loot.Iron = take
		remainingVolume -= float64(take) * WorthIron
	}

	// 3. Titanium
	if poolTitanium > 0 && remainingVolume >= WorthTitanium {
		maxAmt := int(remainingVolume / WorthTitanium)
		take := min(poolTitanium, maxAmt)
		loot.Titanium = take
		remainingVolume -= float64(take) * WorthTitanium
	}

	// 4. Antimatter
	if poolAntimatter > 0 && remainingVolume >= WorthAntimatter {
		maxAmt := int(remainingVolume / WorthAntimatter)
		take := min(poolAntimatter, maxAmt)
		loot.Antimatter = take
		remainingVolume -= float64(take) * WorthAntimatter
	}

	return loot
}

// resolveSpySkirmish performs a very simple skirmish resolution between spy units as a placeholder.
// Later, replace with the step-by-step subtraction algorithm across items.
func resolveSpySkirmish(attackers []MilitaryUnitSnap, atkMods MilitaryModifiers, defenders []MilitaryUnitSnap, defMods MilitaryModifiers) (atkRemaining []MilitaryUnitSnap, defRemaining []MilitaryUnitSnap, attackerWon bool) {
	atk := cloneUnits(attackers)
	def := cloneUnits(defenders)

	atkPower := SumEffectiveAttack(atk, atkMods.AttackMul)
	defPower := SumEffectiveDefence(def, defMods.DefenceMul)

	// Simple outcome: higher total attack vs defence wins
	attackerWon = atkPower >= defPower
	if attackerWon {
		// Placeholder: attackers take no losses
		def = nil
	} else {
		// Placeholder: defenders take no losses
		atk = nil
	}
	return atk, def, attackerWon
}
