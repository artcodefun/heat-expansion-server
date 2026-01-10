package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"

// OperationType represents the type of a military operation at the DTO level.
type OperationType string

const (
	OperationTypeAttack     OperationType = "ATTACK"
	OperationTypeSpy        OperationType = "SPY"
	OperationTypeOccupation OperationType = "OCCUPATION"
)

// OperationPhase represents the lifecycle stage of an operation at the DTO level.
type OperationPhase string

const (
	OperationPhasePending   OperationPhase = "PENDING"
	OperationPhaseOutbound  OperationPhase = "OUTBOUND"
	OperationPhaseAtTarget  OperationPhase = "AT_TARGET"
	OperationPhaseResolving OperationPhase = "RESOLVING"
	OperationPhaseReturning OperationPhase = "RETURNING"
	OperationPhaseCompleted OperationPhase = "COMPLETED"
)

// OperationResult represents the outcome of an operation at the DTO level.
type OperationResult string

const (
	OperationResultUnknown  OperationResult = "UNKNOWN"
	OperationResultSuccess  OperationResult = "SUCCESS"
	OperationResultFailure  OperationResult = "FAILURE"
	OperationResultCanceled OperationResult = "CANCELED"
)

// SpyOutcome represents possible results of a spy operation at the DTO level.
type SpyOutcome string

const (
	SpyOutcomeBlockedByCloaking SpyOutcome = "BLOCKED_BY_CLOAKING_EMPTY_REPORT"
	SpyOutcomeDefeatedBySpies   SpyOutcome = "DEFEATED_BY_DEFENDING_SPIES"
	SpyOutcomeReportProduced    SpyOutcome = "REPORT_PRODUCED"
)

// AttackOutcome represents possible results of an attack operation at the DTO level.
type AttackOutcome string

const (
	AttackOutcomeAttackerWon  AttackOutcome = "ATTACKER_WON"
	AttackOutcomeDefenderHeld AttackOutcome = "DEFENDER_HELD"
)

// MilitaryUnitDTO serializes military unit snapshots.
type MilitaryUnitDTO struct {
	PrototypeID int          `json:"prototype_id"`
	Name        string       `json:"name"`
	ImageURL    string       `json:"image_url"`
	Category    ArmyCategory `json:"category"`
	Attack      int          `json:"attack"`
	Defence     int          `json:"defence"`
	Capacity    int          `json:"capacity"`
	Stealth     int          `json:"stealth"`
	Speed       int          `json:"speed"`
	Count       int          `json:"count"`
}

// DefenseStructureDTO represents defensive structure snapshots.
type DefenseStructureDTO struct {
	PrototypeID int    `json:"prototype_id"`
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Defence     int    `json:"defence"`
	Count       int    `json:"count"`
}

// TrophyStorageItemDTO represents trophy item snapshots.
type TrophyStorageItemDTO struct {
	Prototype StorageItemPrototypeDTO `json:"prototype"`
}

// SpyResultDTO reports spy resolution outcomes.
type SpyResultDTO struct {
	Outcome           SpyOutcome        `json:"outcome"`
	AttackerRemaining []MilitaryUnitDTO `json:"attacker_remaining"`
	DefenderRemaining []MilitaryUnitDTO `json:"defender_remaining"`
	DefendersBefore   []MilitaryUnitDTO `json:"defenders_before"`
}

// AttackResultDTO reports attack outcomes.
type AttackResultDTO struct {
	Outcome             AttackOutcome          `json:"outcome"`
	AttackerRemaining   []MilitaryUnitDTO      `json:"attacker_remaining"`
	DefenderRemaining   []MilitaryUnitDTO      `json:"defender_remaining"`
	RemainingStructures []DefenseStructureDTO  `json:"remaining_structures"`
	Loot                PriceModelDTO          `json:"loot"`
	Trophies            []TrophyStorageItemDTO `json:"trophies"`
	DefendersBefore     []MilitaryUnitDTO      `json:"defenders_before"`
	StructuresBefore    []DefenseStructureDTO  `json:"structures_before"`
}

// MilitaryOperationDTO serializes military operations for HTTP responses.
type MilitaryOperationDTO struct {
	ID                 int               `json:"id"`
	Type               OperationType     `json:"type"`
	Phase              OperationPhase    `json:"phase"`
	Result             OperationResult   `json:"result"`
	SourceBaseID       int               `json:"source_base_id"`
	SourceCoordinates  Vector2iDTO       `json:"source_coordinates"`
	TargetCoordinates  Vector2iDTO       `json:"target_coordinates"`
	OutboundDepartAt   int64             `json:"outbound_depart_at"`
	OutboundArriveAt   int64             `json:"outbound_arrive_at"`
	ReturnDepartAt     int64             `json:"return_depart_at"`
	ReturnArriveAt     int64             `json:"return_arrive_at"`
	CompletedAt        int64             `json:"completed_at"`
	CrystalsSkipPrice  int               `json:"crystals_skip_price"`
	Units              []MilitaryUnitDTO `json:"units"`
	SpyResult          *SpyResultDTO     `json:"spy_result,omitempty"`
	AttackResult       *AttackResultDTO  `json:"attack_result,omitempty"`
	ProducedScanReport *SectorDTO        `json:"produced_scan_report,omitempty"`
}

func militaryUnitsFromReadModel(units []readmodels.MilitaryUnitSnap) []MilitaryUnitDTO {
	out := make([]MilitaryUnitDTO, 0, len(units))
	for _, unit := range units {
		out = append(out, MilitaryUnitDTO{
			PrototypeID: unit.PrototypeID,
			Name:        unit.Name,
			ImageURL:    unit.ImageURL,
			Category:    ArmyCategory(unit.Category),
			Attack:      unit.Attack,
			Defence:     unit.Defence,
			Capacity:    unit.Capacity,
			Stealth:     unit.Stealth,
			Speed:       unit.Speed,
			Count:       unit.Count,
		})
	}
	return out
}

func defenseStructuresFromReadModel(structs []readmodels.DefenseStructureSnap) []DefenseStructureDTO {
	out := make([]DefenseStructureDTO, 0, len(structs))
	for _, s := range structs {
		out = append(out, DefenseStructureDTO{
			PrototypeID: s.PrototypeID,
			Name:        s.Name,
			ImageURL:    s.ImageURL,
			Defence:     s.Defence,
			Count:       s.Count,
		})
	}
	return out
}

func spyResultFromReadModel(res *readmodels.SpyResult) *SpyResultDTO {
	if res == nil {
		return nil
	}
	return &SpyResultDTO{
		Outcome:           SpyOutcome(res.Outcome),
		AttackerRemaining: militaryUnitsFromReadModel(res.AttackerRemaining),
		DefenderRemaining: militaryUnitsFromReadModel(res.DefenderRemaining),
		DefendersBefore:   militaryUnitsFromReadModel(res.DefendersBefore),
	}
}

func attackResultFromReadModel(res *readmodels.AttackResult) *AttackResultDTO {
	if res == nil {
		return nil
	}
	dto := &AttackResultDTO{
		Outcome:             AttackOutcome(res.Outcome),
		AttackerRemaining:   militaryUnitsFromReadModel(res.AttackerRemaining),
		DefenderRemaining:   militaryUnitsFromReadModel(res.DefenderRemaining),
		RemainingStructures: defenseStructuresFromReadModel(res.RemainingStructures),
		Loot:                PriceModelFromReadModel(res.Loot),
		DefendersBefore:     militaryUnitsFromReadModel(res.DefendersBefore),
		StructuresBefore:    defenseStructuresFromReadModel(res.StructuresBefore),
	}
	if len(res.Trophies) > 0 {
		dto.Trophies = make([]TrophyStorageItemDTO, 0, len(res.Trophies))
		for _, t := range res.Trophies {
			dto.Trophies = append(dto.Trophies, TrophyStorageItemDTO{
				Prototype: mapStorageItemPrototype(t.Prototype),
			})
		}
	}
	return dto
}

func OperationFromReadModel(m *readmodels.MilitaryOperation) MilitaryOperationDTO {
	dto := MilitaryOperationDTO{
		ID:                m.ID,
		Type:              OperationType(m.Type),
		Phase:             OperationPhase(m.Phase),
		Result:            OperationResult(m.Result),
		SourceBaseID:      m.SourceBaseID,
		SourceCoordinates: Vector2iDTO{X: m.SourceCoordinates.X, Y: m.SourceCoordinates.Y},
		TargetCoordinates: Vector2iDTO{X: m.TargetCoordinates.X, Y: m.TargetCoordinates.Y},
		OutboundDepartAt:  m.OutboundDepartAt,
		OutboundArriveAt:  m.OutboundArriveAt,
		ReturnDepartAt:    m.ReturnDepartAt,
		ReturnArriveAt:    m.ReturnArriveAt,
		CompletedAt:       m.CompletedAt,
		CrystalsSkipPrice: m.CrystalsSkipPrice,
		Units:             militaryUnitsFromReadModel(m.Units),
		SpyResult:         spyResultFromReadModel(m.SpyResult),
		AttackResult:      attackResultFromReadModel(m.AttackResult),
	}
	if m.ProducedScanReport != nil {
		report := SectorScanReportFromReadModel(m.ProducedScanReport)
		dto.ProducedScanReport = &report
	}
	return dto
}
