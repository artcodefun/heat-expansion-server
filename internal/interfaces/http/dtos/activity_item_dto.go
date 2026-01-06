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

// OperationActivityDTO mirrors readmodels.OperationActivity, embedding the full operation readmodel.
type OperationActivityDTO struct {
	OpID              int                     `json:"opId"`
	Subtype           MilitaryActivitySubtype `json:"subtype"`
	Role              OperationRole           `json:"role"`
	Operation         *MilitaryOperationDTO   `json:"operation,omitempty"`
	PriorOpponentScan *SectorDTO              `json:"priorOpponentScan,omitempty"`
}

// ScanActivityDTO mirrors readmodels.ScanActivity, embedding the full scan report readmodel.
type ScanActivityDTO struct {
	ReportID int        `json:"reportId"`
	Report   *SectorDTO `json:"report,omitempty"`
}

// RadarActivityDTO presents a detected incoming hostility (placeholder for future wiring).
type RadarActivityDTO struct {
	OpID       int         `json:"opId"`
	DetectedAt int64       `json:"detectedAt"`
	EtaAtBase  int64       `json:"etaAtBase"`
	Source     Vector2iDTO `json:"source"`
	Target     Vector2iDTO `json:"target"`
	Threat     ThreatDTO   `json:"threat"`
}

type ThreatDTO struct {
	Attack  int `json:"attack"`
	Defence int `json:"defence"`
}

func operationActivityFromReadModel(op *readmodels.OperationActivity) *OperationActivityDTO {
	if op == nil {
		return nil
	}
	dto := &OperationActivityDTO{
		OpID:    op.OpID,
		Subtype: MilitaryActivitySubtype(op.Subtype),
		Role:    OperationRole(op.Role),
	}
	if op.Operation != nil {
		m := OperationFromReadModel(op.Operation)
		dto.Operation = &m
	}
	if op.PriorOpponentScan != nil {
		report := SectorScanReportFromReadModel(op.PriorOpponentScan)
		dto.PriorOpponentScan = &report
	}
	return dto
}

func scanActivityFromReadModel(scan *readmodels.ScanActivity) *ScanActivityDTO {
	if scan == nil {
		return nil
	}
	dto := &ScanActivityDTO{ReportID: scan.ReportID}
	if scan.Report != nil {
		report := SectorScanReportFromReadModel(scan.Report)
		dto.Report = &report
	}
	return dto
}

func radarActivityFromReadModel(radar *readmodels.RadarActivity) *RadarActivityDTO {
	if radar == nil {
		return nil
	}
	return &RadarActivityDTO{
		OpID:       radar.OpID,
		DetectedAt: radar.DetectedAt,
		EtaAtBase:  radar.EtaAtBase,
		Source:     Vector2iFromReadModel(radar.SourceCoordinates),
		Target:     Vector2iFromReadModel(radar.TargetCoordinates),
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
