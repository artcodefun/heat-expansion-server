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

// MilitaryModifiersDTO represents multipliers for military stats at the DTO level.
type MilitaryModifiersDTO struct {
	AttackMul   float32 `json:"attack_mul"`
	DefenceMul  float32 `json:"defence_mul"`
	StealthMul  float32 `json:"stealth_mul"`
	CapacityMul float32 `json:"capacity_mul"`
	SpeedMul    float32 `json:"speed_mul"`
}

// StorageItemSnapDTO represents a snapshot of a storage item at the DTO level.
type StorageItemSnapDTO struct {
	PrototypeID      int                     `json:"prototype_id"`
	Name             string                  `json:"name"`
	ShortDescription string                  `json:"short_description"`
	ImageURL         string                  `json:"image_url"`
	Category         StorageCategory         `json:"category"`
	BuffData         *BuffStorageDataDTO     `json:"buff_data,omitempty"`
	ArtifactData     *ArtifactStorageDataDTO `json:"artifact_data,omitempty"`
}

// SpyResultDTO reports spy resolution outcomes.
type SpyResultDTO struct {
	Outcome                SpyOutcome           `json:"outcome"`
	AttackerRemaining      []MilitaryUnitDTO    `json:"attacker_remaining"`
	DefenderRemaining      []MilitaryUnitDTO    `json:"defender_remaining"`
	DefendersBefore        []MilitaryUnitDTO    `json:"defenders_before"`
	DefenderStorageSnaps   []StorageItemSnapDTO `json:"defender_storage_snaps"`
	TotalDefenderModifiers MilitaryModifiersDTO `json:"total_defender_modifiers"`
}

// AttackResultDTO reports attack outcomes.
type AttackResultDTO struct {
	Outcome                AttackOutcome          `json:"outcome"`
	AttackerRemaining      []MilitaryUnitDTO      `json:"attacker_remaining"`
	DefenderRemaining      []MilitaryUnitDTO      `json:"defender_remaining"`
	RemainingStructures    []DefenseStructureDTO  `json:"remaining_structures"`
	Loot                   PriceModelDTO          `json:"loot"`
	Trophies               []TrophyStorageItemDTO `json:"trophies"`
	DefendersBefore        []MilitaryUnitDTO      `json:"defenders_before"`
	StructuresBefore       []DefenseStructureDTO  `json:"structures_before"`
	DefenderStorageSnaps   []StorageItemSnapDTO   `json:"defender_storage_snaps"`
	TotalDefenderModifiers MilitaryModifiersDTO   `json:"total_defender_modifiers"`
}

// MilitaryOperationDTO serializes military operations for HTTP responses.
type MilitaryOperationDTO struct {
	ID                 int                  `json:"id"`
	Type               OperationType        `json:"type"`
	Phase              OperationPhase       `json:"phase"`
	Result             OperationResult      `json:"result"`
	SourceBaseID       int                  `json:"source_base_id"`
	SourceCoordinates  Vector2iDTO          `json:"source_coordinates"`
	TargetCoordinates  Vector2iDTO          `json:"target_coordinates"`
	OutboundDepartAt   int64                `json:"outbound_depart_at"`
	OutboundArriveAt   int64                `json:"outbound_arrive_at"`
	ReturnDepartAt     int64                `json:"return_depart_at"`
	ReturnArriveAt     int64                `json:"return_arrive_at"`
	CompletedAt        int64                `json:"completed_at"`
	CrystalsSkipPrice  int                  `json:"crystals_skip_price"`
	Units              []MilitaryUnitDTO    `json:"units"`
	StorageSnaps       []StorageItemSnapDTO `json:"storage_snaps"`
	TotalModifiers     MilitaryModifiersDTO `json:"total_modifiers"`
	SpyResult          *SpyResultDTO        `json:"spy_result,omitempty"`
	AttackResult       *AttackResultDTO     `json:"attack_result,omitempty"`
	ProducedScanReport *SectorDTO           `json:"produced_scan_report,omitempty"`
	PriorScanReport    *SectorDTO           `json:"prior_scan_report,omitempty"`
}

func MilitaryUnitsFromReadModel(units []readmodels.MilitaryUnitSnap) []MilitaryUnitDTO {
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

func trophiesFromReadModel(trophies []readmodels.TrophyStorageItem) []TrophyStorageItemDTO {
	out := make([]TrophyStorageItemDTO, 0, len(trophies))
	for _, t := range trophies {
		out = append(out, TrophyStorageItemDTO{
			Prototype: mapStorageItemPrototype(t.Prototype),
		})
	}
	return out
}

func SpyResultFromReadModel(res *readmodels.SpyResult) *SpyResultDTO {
	if res == nil {
		return nil
	}
	return &SpyResultDTO{
		Outcome:                SpyOutcome(res.Outcome),
		AttackerRemaining:      MilitaryUnitsFromReadModel(res.AttackerRemaining),
		DefenderRemaining:      MilitaryUnitsFromReadModel(res.DefenderRemaining),
		DefendersBefore:        MilitaryUnitsFromReadModel(res.DefendersBefore),
		DefenderStorageSnaps:   storageItemSnapsFromReadModel(res.DefenderStorageSnaps),
		TotalDefenderModifiers: MilitaryModifiersFromReadModel(res.TotalDefenderModifiers),
	}
}

func AttackResultFromReadModel(res *readmodels.AttackResult) *AttackResultDTO {
	if res == nil {
		return nil
	}
	return &AttackResultDTO{
		Outcome:                AttackOutcome(res.Outcome),
		AttackerRemaining:      MilitaryUnitsFromReadModel(res.AttackerRemaining),
		DefenderRemaining:      MilitaryUnitsFromReadModel(res.DefenderRemaining),
		RemainingStructures:    defenseStructuresFromReadModel(res.RemainingStructures),
		Loot:                   PriceModelFromReadModel(res.Loot),
		Trophies:               trophiesFromReadModel(res.Trophies),
		DefendersBefore:        MilitaryUnitsFromReadModel(res.DefendersBefore),
		StructuresBefore:       defenseStructuresFromReadModel(res.StructuresBefore),
		DefenderStorageSnaps:   storageItemSnapsFromReadModel(res.DefenderStorageSnaps),
		TotalDefenderModifiers: MilitaryModifiersFromReadModel(res.TotalDefenderModifiers),
	}
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
		Units:             MilitaryUnitsFromReadModel(m.Units),
		StorageSnaps:      storageItemSnapsFromReadModel(m.StorageSnaps),
		TotalModifiers:    MilitaryModifiersFromReadModel(m.TotalModifiers),
		SpyResult:         SpyResultFromReadModel(m.SpyResult),
		AttackResult:      AttackResultFromReadModel(m.AttackResult),
	}
	if m.ProducedScanReport != nil {
		report := SectorScanReportFromReadModel(m.ProducedScanReport)
		dto.ProducedScanReport = &report
	}
	if m.PriorScanReport != nil {
		report := SectorScanReportFromReadModel(m.PriorScanReport)
		dto.PriorScanReport = &report
	}
	return dto
}

func storageItemSnapsFromReadModel(snaps []readmodels.StorageItemSnap) []StorageItemSnapDTO {
	if len(snaps) == 0 {
		return []StorageItemSnapDTO{}
	}
	out := make([]StorageItemSnapDTO, 0, len(snaps))
	for _, s := range snaps {
		var buff *BuffStorageDataDTO
		if s.BuffData != nil {
			buff = &BuffStorageDataDTO{
				Type:            BuffType(s.BuffData.Type),
				Value:           s.BuffData.Value,
				DurationSeconds: s.BuffData.DurationSeconds,
			}
		}
		var artifact *ArtifactStorageDataDTO
		if s.ArtifactData != nil {
			artifact = &ArtifactStorageDataDTO{
				Type:  ArtifactEffectType(s.ArtifactData.Type),
				Value: s.ArtifactData.Value,
			}
		}
		out = append(out, StorageItemSnapDTO{
			PrototypeID:      s.PrototypeID,
			Name:             s.Name,
			ShortDescription: s.ShortDescription,
			ImageURL:         s.ImageURL,
			Category:         StorageCategory(s.Category),
			BuffData:         buff,
			ArtifactData:     artifact,
		})
	}
	return out
}

func MilitaryModifiersFromReadModel(m readmodels.MilitaryModifiers) MilitaryModifiersDTO {
	return MilitaryModifiersDTO{
		AttackMul:   m.AttackMul,
		DefenceMul:  m.DefenceMul,
		StealthMul:  m.StealthMul,
		CapacityMul: m.CapacityMul,
		SpeedMul:    m.SpeedMul,
	}
}

func OperationsFromReadModels(items []*readmodels.MilitaryOperation) []MilitaryOperationDTO {
	out := make([]MilitaryOperationDTO, 0, len(items))
	for _, item := range items {
		out = append(out, OperationFromReadModel(item))
	}
	return out
}
