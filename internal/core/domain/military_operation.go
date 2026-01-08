package domain

import (
	"fmt"
	"math"
	"math/rand"
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
	AttackerRemaining []MilitaryUnit
	DefenderRemaining []MilitaryUnit
	// New: snapshot of defenders before resolution (for UI diffs)
	DefendersBefore []MilitaryUnit
}

type AttackResult struct {
	Outcome             AttackOutcome
	AttackerRemaining   []MilitaryUnit
	DefenderRemaining   []MilitaryUnit
	RemainingStructures []DefenseStructure
	Loot                PriceModel // what attackers managed to carry back; computed elsewhere
	// New: snapshots for UI to show casualties/damage
	DefendersBefore  []MilitaryUnit
	StructuresBefore []DefenseStructure
}

// MilitaryOperation models an attack or spy op traveling between sectors and resolving on arrival.
type MilitaryOperation struct {
	EventProducer
	ID           int
	Type         MilitaryOperationType
	OwnerUserID  int
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
	Units []MilitaryUnit

	// Results (only one will be populated depending on Type)
	SpyResult    *SpyResult
	AttackResult *AttackResult
}

// NewAttackOperation creates an ATTACK operation in transit.
// It validates that at least one unit is provided and that source/target are different.
func NewAttackOperation(ownerUserID, sourceBaseID int, source, target Vector2i, units []MilitaryUnit) (*MilitaryOperation, error) {
	if source == target {
		return nil, fmt.Errorf("source and target coordinates must be different")
	}
	if len(units) == 0 {
		return nil, fmt.Errorf("no units provided for attack operation")
	}
	op := &MilitaryOperation{
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
	}
	return op, nil
}

// NewSpyOperation creates a SPY operation in transit.
// It validates that at least one unit is provided, targeting a different sector, and that all units are spies.
func NewSpyOperation(ownerUserID, sourceBaseID int, source, target Vector2i, spies []MilitaryUnit) (*MilitaryOperation, error) {
	if source == target {
		return nil, fmt.Errorf("source and target coordinates must be different")
	}
	if len(spies) == 0 {
		return nil, fmt.Errorf("no units provided for spy operation")
	}
	for _, u := range spies {
		if u.Category != ArmyCategorySpy {
			return nil, fmt.Errorf("spy operations require only spy units")
		}
	}
	op := &MilitaryOperation{
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
	travelSeconds := computeTravelSecondsBetween(op.SourceCoordinates, op.TargetCoordinates, op.Units)
	op.OutboundDepartAt = now
	op.OutboundArriveAt = now + travelSeconds
	// Base skip price proportional to total travel time (similar to production queues)
	op.CrystalsSkipPrice = int(travelSeconds / 60)
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
func (op *MilitaryOperation) ResolveSpy(targetCloakingStrength int, defendingSpies []MilitaryUnit) *SpyResult {
	if op.Type != MilitaryOperationTypeSpy || (op.Phase != OperationPhaseAtTarget && op.Phase != OperationPhaseOutbound) {
		return op.SpyResult
	}
	op.Phase = OperationPhaseResolving

	// Ensure we consider only spy-category units on the attacker side
	attackers := filterUnitsByCategory(op.Units, ArmyCategorySpy)
	attackerStealth := sumStealth(attackers)

	// Snapshot defenders before resolving for UI diffs
	defBefore := cloneUnits(defendingSpies)

	// 1) Cloaking check: if target cloaking >= attacker stealth -> empty report; units unharmed.
	if targetCloakingStrength >= attackerStealth {
		res := &SpyResult{
			Outcome: SpyOutcomeBlockedByCloaking,
			// Everyone survives in this outcome
			AttackerRemaining: cloneUnits(attackers),
			DefenderRemaining: cloneUnits(defendingSpies),
			DefendersBefore:   defBefore,
		}
		op.SpyResult = res
		op.Result = OperationResultSuccess
		op.AddEvent(NewMilitaryOperationResolvedEvent(op.ID, op.Result))
		return res
	}

	// 2) Skirmish with defending spies
	atkRemaining, defRemaining, attackerWon := resolveSpySkirmish(attackers, defendingSpies)
	if !attackerWon {
		res := &SpyResult{
			Outcome:           SpyOutcomeDefeatedBySpies,
			AttackerRemaining: atkRemaining,
			DefenderRemaining: defRemaining,
			DefendersBefore:   defBefore,
		}
		op.SpyResult = res
		op.Result = OperationResultFailure
		op.AddEvent(NewMilitaryOperationResolvedEvent(op.ID, op.Result))
		return res
	}

	// 3) Successful intelligence report
	res := &SpyResult{
		Outcome:           SpyOutcomeReportProduced,
		AttackerRemaining: atkRemaining,
		DefenderRemaining: defRemaining,
		DefendersBefore:   defBefore,
	}
	op.SpyResult = res
	op.Result = OperationResultSuccess
	op.AddEvent(NewMilitaryOperationResolvedEvent(op.ID, op.Result))
	return res
}

// ResolveAttack resolves an attack using a simplified power comparison as a placeholder.
// A more detailed sequential algorithm can replace this later.
func (op *MilitaryOperation) ResolveAttack(defenders []MilitaryUnit, structures []DefenseStructure, availableResourcePool PriceModel) *AttackResult {
	if op.Type != MilitaryOperationTypeAttack || (op.Phase != OperationPhaseAtTarget && op.Phase != OperationPhaseOutbound) {
		return op.AttackResult
	}
	op.Phase = OperationPhaseResolving

	// capture "before" snapshots for UI diffs
	defBefore := cloneUnits(defenders)
	structBefore := cloneStructures(structures)

	atkRemain, defRemain, structRemain, attackerWon := resolveAttackCombat(cloneUnits(op.Units), cloneUnits(defenders), cloneStructures(structures))

	result := &AttackResult{
		DefendersBefore:     defBefore,
		StructuresBefore:    structBefore,
		AttackerRemaining:   atkRemain,
		DefenderRemaining:   defRemain,
		RemainingStructures: structRemain,
		// Compute loot here based on remaining attackers' capacity and available resources at target
		Loot: computeLoadFromLocation(atkRemain, availableResourcePool),
	}
	if attackerWon {
		result.Outcome = AttackOutcomeAttackerWon
		op.Result = OperationResultSuccess
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
// Damage removes full instances only.
// Returns remaining stacks and whether attackers won (defenders and structures eliminated).
func resolveAttackCombat(attackers []MilitaryUnit, defenders []MilitaryUnit, structures []DefenseStructure) (atkRemaining []MilitaryUnit, defRemaining []MilitaryUnit, structRemaining []DefenseStructure, attackerWon bool) {
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

		// Compute damages from current forces
		attDmg := sumAttack(attackers)
		// Defenders and structures both contribute damage using their Defence as power
		defDmg := sumDefence(defenders) + sumStructureDefence(structures)

		// If neither side can deal damage, stop to avoid infinite loop
		if attDmg == 0 && defDmg == 0 {
			attackerWon = false
			break
		}

		// Apply attacker damage to defenders then structures
		var attOverflow int
		defenders, attOverflow = applyDamageToUnits(defenders, attDmg, true)
		structures, _ = applyDamageToStructures(structures, attOverflow)

		// Apply defender damage to attackers
		attackers, _ = applyDamageToUnits(attackers, defDmg, false)

		// Optional normalization could be applied here if needed
	}

	return attackers, defenders, structures, attackerWon
}

// applyDamageToUnits applies damage to units. When isDefensive=true, per-instance HP is Defence; otherwise HP is Attack.
// Returns updated units and any leftover damage not consumed.
func applyDamageToUnits(units []MilitaryUnit, damage int, isDefensive bool) ([]MilitaryUnit, int) {
	if damage <= 0 || len(units) == 0 {
		return units, damage
	}
	out := make([]MilitaryUnit, 0, len(units))
	remainingDamage := damage
	for _, u := range units {
		if u.Count <= 0 {
			continue
		}
		var hp int
		if isDefensive {
			hp = u.Defence
		} else {
			hp = u.Attack
		}
		// Guard against non-positive HP creating immortal stacks
		if hp <= 0 {
			hp = 1
		}
		// kills = min(count, floor(damage / hpPerUnit))
		possibleKills := remainingDamage / hp
		if possibleKills > 0 {
			if possibleKills >= u.Count {
				remainingDamage -= u.Count * hp
				// all killed -> skip append
				continue
			}
			// partial kills
			u.Count -= possibleKills
			remainingDamage -= possibleKills * hp
		}
		if u.Count > 0 {
			out = append(out, u)
		}
	}
	return out, remainingDamage
}

// (No separate variants; isDefensive flag selects HP stat.)

// applyDamageToStructures applies damage to structures similarly based on Defence HP per structure.
func applyDamageToStructures(structs []DefenseStructure, damage int) ([]DefenseStructure, int) {
	if damage <= 0 || len(structs) == 0 {
		return structs, damage
	}
	out := make([]DefenseStructure, 0, len(structs))
	remainingDamage := damage
	for _, s := range structs {
		if s.Count <= 0 {
			continue
		}
		hp := s.Defence
		if hp <= 0 {
			hp = 1
		}
		possibleKills := remainingDamage / hp
		if possibleKills > 0 {
			if possibleKills >= s.Count {
				remainingDamage -= s.Count * hp
				continue
			}
			s.Count -= possibleKills
			remainingDamage -= possibleKills * hp
		}
		if s.Count > 0 {
			out = append(out, s)
		}
	}
	return out, remainingDamage
}

// Helpers to filter and count
func filterZeroCountUnits(units []MilitaryUnit) []MilitaryUnit {
	if len(units) == 0 {
		return nil
	}
	out := make([]MilitaryUnit, 0, len(units))
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

func filterZeroCountStructures(structs []DefenseStructure) []DefenseStructure {
	if len(structs) == 0 {
		return nil
	}
	out := make([]DefenseStructure, 0, len(structs))
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

func totalCount(units []MilitaryUnit) int {
	if len(units) == 0 {
		return 0
	}
	sum := 0
	for _, u := range units {
		sum += u.Count
	}
	return sum
}

func totalStructCount(structs []DefenseStructure) int {
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
	var returningUnits []MilitaryUnit
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

	travelSeconds := computeTravelSecondsBetween(op.TargetCoordinates, op.SourceCoordinates, returningUnits)
	op.ReturnDepartAt = now
	op.ReturnArriveAt = now + travelSeconds
	// Recompute skip price for the return leg
	op.CrystalsSkipPrice = int(travelSeconds / 60)
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

// Cancel marks the operation canceled and optionally starts immediate return if outbound or at target.
func (op *MilitaryOperation) Cancel(startReturn bool) {
	op.Result = OperationResultCanceled
	if startReturn && (op.Phase == OperationPhaseOutbound || op.Phase == OperationPhaseAtTarget || op.Phase == OperationPhaseResolving) {
		op.StartReturn()
	}
}

// Travel time helpers

// computeTravelSecondsBetween computes travel time as ceil(distance / slowestSpeed).
// Distance is Manhattan on sector coordinates for now; speed is taken from the slowest unit.
func computeTravelSecondsBetween(a, b Vector2i, units []MilitaryUnit) int64 {
	// Euclidean distance scaled by 1000, then divide by slowest speed and ceil
	dist := euclideanScaled(a, b)
	if dist <= 0 {
		return 0
	}
	s := slowestSpeed(units)
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
		return 0, fmt.Errorf("operation is not in outbound travel")
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
		return 0, fmt.Errorf("operation has no movement")
	}
	b := 2 * (dx*wx + dy*wy)
	c := wx*wx + wy*wy - r*r

	// If already inside (c <= 0), it entered at the start of travel
	if c <= 0 {
		return op.OutboundDepartAt, nil
	}

	delta := b*b - 4*a*c
	if delta < 0 {
		return 0, fmt.Errorf("operation never enters circle")
	}

	sqrtDelta := math.Sqrt(delta)
	k1 := (-b - sqrtDelta) / (2 * a)

	if k1 < 0 {
		// Moving away from the circle since we started outside (c > 0)
		return 0, fmt.Errorf("operation moving away from circle")
	}
	if k1 > 1.0 {
		return 0, fmt.Errorf("operation reaches target before entering circle")
	}

	travelDuration := op.OutboundArriveAt - op.OutboundDepartAt
	enterAt := op.OutboundDepartAt + int64(float64(travelDuration)*k1)
	return enterAt, nil
}

func filterUnitsByCategory(units []MilitaryUnit, cat ArmyCategory) []MilitaryUnit {
	if len(units) == 0 {
		return nil
	}
	out := make([]MilitaryUnit, 0, len(units))
	for _, u := range units {
		if u.Category == cat && u.Count > 0 {
			out = append(out, u)
		}
	}
	return out
}

func (op *MilitaryOperation) TotalStealth() int {
	return sumStealth(op.Units)
}

func (op *MilitaryOperation) TotalSpeed() int {
	return slowestSpeed(op.Units)
}

func (op *MilitaryOperation) TotalCapacity() int {
	return sumCapacity(op.Units)
}

// computeLoadFromLocation randomly fills available carrying capacity using the available
// resource pool. Allocation is random but bounded by both capacity and available amounts.
// This is used after resolution, based on attacker survivors.
func computeLoadFromLocation(remaining []MilitaryUnit, available PriceModel) PriceModel {
	capacity := sumCapacity(remaining)
	if capacity <= 0 {
		return PriceModel{}
	}
	// Clamp available to non-negative
	pool := PriceModel{
		Credits:    maxInt(available.Credits, 0),
		Iron:       maxInt(available.Iron, 0),
		Titanium:   maxInt(available.Titanium, 0),
		Antimatter: maxInt(available.Antimatter, 0),
	}
	// Quick path: total available smaller than capacity -> take all
	totalAvail := pool.Credits + pool.Iron + pool.Titanium + pool.Antimatter
	if totalAvail == 0 {
		return PriceModel{}
	}
	if totalAvail <= capacity {
		return pool
	}
	// Randomly allocate one unit at a time up to capacity.
	rng := rand.New(rand.NewSource(NowUnixNano()))
	loot := PriceModel{}
	for taken := 0; taken < capacity; taken++ {
		currentTotal := pool.Credits + pool.Iron + pool.Titanium + pool.Antimatter
		if currentTotal == 0 {
			break
		}
		pick := rng.Intn(currentTotal)
		// Credits bucket
		if pick < pool.Credits {
			pool.Credits--
			loot.Credits++
			continue
		}
		pick -= pool.Credits
		// Iron bucket
		if pick < pool.Iron {
			pool.Iron--
			loot.Iron++
			continue
		}
		pick -= pool.Iron
		// Titanium bucket
		if pick < pool.Titanium {
			pool.Titanium--
			loot.Titanium++
			continue
		}
		// Antimatter bucket
		if pool.Antimatter > 0 {
			pool.Antimatter--
			loot.Antimatter++
		}
	}
	return loot
}

// resolveSpySkirmish performs a very simple skirmish resolution between spy units as a placeholder.
// Later, replace with the step-by-step subtraction algorithm across items.
func resolveSpySkirmish(attackers, defenders []MilitaryUnit) (atkRemaining []MilitaryUnit, defRemaining []MilitaryUnit, attackerWon bool) {
	atk := cloneUnits(attackers)
	def := cloneUnits(defenders)

	atkPower := sumAttack(atk)
	defPower := sumDefence(def)

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

// Converters from base deploy-ready items to operation units

// OperationUnitFromDeployed converts a single deploy-ready stack into an OperationUnit snapshot.
func OperationUnitFromDeployed(d DeploymentReadyItem) MilitaryUnit {
	return MilitaryUnit{
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

// OperationUnitsFromDeployed returns operation units for a list of deploy-ready stacks.
func OperationUnitsFromDeployed(items []DeploymentReadyItem) []MilitaryUnit {
	if len(items) == 0 {
		return nil
	}
	out := make([]MilitaryUnit, 0, len(items))
	for _, d := range items {
		out = append(out, OperationUnitFromDeployed(d))
	}
	return out
}
