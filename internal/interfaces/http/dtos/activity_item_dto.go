package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"

// ActivityKind mirrors readmodels.ActivityKind for HTTP.
type ActivityKind string

const (
	ActivityKindMilitary ActivityKind = "MILITARY"
	ActivityKindScan     ActivityKind = "SCAN"
	ActivityKindRadar    ActivityKind = "RADAR"
	ActivityKindTrade    ActivityKind = "TRADE"
)

// MilitaryActivitySubtype mirrors readmodels.MilitaryActivitySubtype.
type MilitaryActivitySubtype string

const (
	MilitaryActivitySubtypeAttack  MilitaryActivitySubtype = "ATTACK"
	MilitaryActivitySubtypeSpy     MilitaryActivitySubtype = "SPY"
	MilitaryActivitySubtypeDefense MilitaryActivitySubtype = "DEFENSE"
)

// OperationRole mirrors readmodels.OperationRole.
type OperationRole string

const (
	OperationRoleAttacker OperationRole = "ATTACKER"
	OperationRoleDefender OperationRole = "DEFENDER"
)

// ActivityItemDTO is a unified envelope for different activity kinds.
type ActivityItemDTO struct {
	ID        int          `json:"id"`
	Kind      ActivityKind `json:"kind"`
	CreatedAt int64        `json:"createdAt"`
	BaseID    int          `json:"baseId"`

	// New: UI routing helpers
	// Category removed; Kind + Subtype are sufficient for UI routing.

	Operation *OperationActivityDTO `json:"operation,omitempty"`
	Scan      *ScanActivityDTO      `json:"scan,omitempty"`
	Radar     *RadarActivityDTO     `json:"radar,omitempty"`
}

// OperationActivityDTO presents an operation in the activities list.
type OperationActivityDTO struct {
	OpID       int                     `json:"opId"`
	OpType     string                  `json:"opType"`
	Subtype    MilitaryActivitySubtype `json:"subtype"`
	Role       OperationRole           `json:"role"`
	Phase      string                  `json:"phase"`
	Result     string                  `json:"result"`
	Source     SectorRef               `json:"source"`
	Target     SectorRef               `json:"target"`
	Timestamps OperationTimestamps     `json:"timestamps"`
	Outcome    *OperationOutcomeDTO    `json:"outcome,omitempty"`
}

type SectorRef struct {
	Coordinates Vector2iDTO `json:"coordinates"`
}

type OperationTimestamps struct {
	OutboundDepartAt int64 `json:"outboundDepartAt"`
	OutboundArriveAt int64 `json:"outboundArriveAt"`
	ReturnDepartAt   int64 `json:"returnDepartAt"`
	ReturnArriveAt   int64 `json:"returnArriveAt"`
	CompletedAt      int64 `json:"completedAt"`
}

// OperationOutcomeDTO summarizes post-resolution results useful for the UI.
type OperationOutcomeDTO struct {
	Loot          PriceModelDTO `json:"loot"`
	UnitsSent     int           `json:"unitsSent"`
	UnitsSurvived int           `json:"unitsSurvived"`
	// Defender-side details when viewing defence outcome
	ResourcesLost  PriceModelDTO      `json:"resourcesLost"`
	UnitsLost      []UnitLossDTO      `json:"unitsLost"`
	StructuresLost []StructureLossDTO `json:"structuresLost"`
}

// ScanActivityDTO presents a scan report.
type ScanActivityDTO struct {
	ReportID          int         `json:"reportId"`
	Coordinates       Vector2iDTO `json:"coordinates"`
	SectorType        string      `json:"sectorType"`
	IsCloaked         bool        `json:"isCloaked"`
	Info              ScanInfoDTO `json:"info"`
	SourceOperationID int         `json:"sourceOperationId"`
}

// RadarActivityDTO presents a detected incoming hostility (placeholder for future wiring).
type RadarActivityDTO struct {
	OpID       int       `json:"opId"`
	DetectedAt int64     `json:"detectedAt"`
	EtaAtBase  int64     `json:"etaAtBase"`
	Source     SectorRef `json:"source"`
	Target     SectorRef `json:"target"`
	Threat     ThreatDTO `json:"threat"`
}

type ThreatDTO struct {
	Attack  int `json:"attack"`
	Defence int `json:"defence"`
}

// Note: defence reports are now embedded into OperationActivityDTO.Outcome

type UnitLossDTO struct {
	PrototypeID int    `json:"prototypeId"`
	Name        string `json:"name"`
	Lost        int    `json:"lost"`
	Remaining   int    `json:"remaining"`
}

type StructureLossDTO struct {
	PrototypeID int    `json:"prototypeId"`
	Name        string `json:"name"`
	Destroyed   int    `json:"destroyed"`
	Remaining   int    `json:"remaining"`
}

func operationActivityFromReadModel(op *readmodels.OperationActivity) *OperationActivityDTO {
	if op == nil {
		return nil
	}
	return &OperationActivityDTO{
		OpID:    op.OpID,
		Subtype: MilitaryActivitySubtype(op.Subtype),
		Role:    OperationRole(op.Role),
	}
}

func scanActivityFromReadModel(scan *readmodels.ScanActivity) *ScanActivityDTO {
	if scan == nil {
		return nil
	}
	return &ScanActivityDTO{
		ReportID: scan.ReportID,
	}
}

func radarActivityFromReadModel(radar *readmodels.RadarActivity) *RadarActivityDTO {
	if radar == nil {
		return nil
	}
	return &RadarActivityDTO{
		OpID:       radar.OpID,
		DetectedAt: radar.DetectedAt,
		EtaAtBase:  radar.EtaAtBase,
		Source:     SectorRef{Coordinates: Vector2iFromReadModel(radar.SourceCoordinates)},
		Target:     SectorRef{Coordinates: Vector2iFromReadModel(radar.TargetCoordinates)},
		Threat:     ThreatDTO{Attack: radar.Threat.Attack, Defence: radar.Threat.Defence},
	}
}

func ActivityItemDTOFromReadModel(a *readmodels.ActivityItem) ActivityItemDTO {
	dto := ActivityItemDTO{
		ID:        a.ID,
		Kind:      ActivityKind(a.Kind),
		CreatedAt: a.CreatedAt,
		BaseID:    a.BaseID,
		Operation: operationActivityFromReadModel(a.Operation),
		Scan:      scanActivityFromReadModel(a.Scan),
		Radar:     radarActivityFromReadModel(a.Radar),
	}
	return dto
}

func ActivityItemsFromReadModels(items []*readmodels.ActivityItem) []ActivityItemDTO {
	out := make([]ActivityItemDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ActivityItemDTOFromReadModel(item))
	}
	return out
}
