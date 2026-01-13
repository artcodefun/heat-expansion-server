package dtos

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
)

// ActivityKind mirrors readmodels.ActivityKind for HTTP.
type ActivityKind string

const (
	ActivityKindOffense ActivityKind = "OFFENSE"
	ActivityKindDefense ActivityKind = "DEFENSE"
	ActivityKindScan    ActivityKind = "SCAN"
	ActivityKindRadar   ActivityKind = "RADAR"
	ActivityKindTrade   ActivityKind = "TRADE"
)

// OffenseActivitySubtype mirrors readmodels.OffenseActivitySubtype.
type OffenseActivitySubtype string

const (
	OffenseActivitySubtypeAttack OffenseActivitySubtype = "ATTACK"
	OffenseActivitySubtypeSpy    OffenseActivitySubtype = "SPY"
)

// DefenseActivitySubtype mirrors readmodels.DefenseActivitySubtype.
type DefenseActivitySubtype string

const (
	DefenseActivitySubtypeAttack DefenseActivitySubtype = "ATTACK"
	DefenseActivitySubtypeSpy    DefenseActivitySubtype = "SPY"
)

// ActivityItemDTO is a unified envelope for different activity kinds.
type ActivityItemDTO struct {
	ID        int          `json:"id"`
	Kind      ActivityKind `json:"kind"`
	CreatedAt int64        `json:"createdAt"`
	BaseID    int          `json:"baseId"`

	Offense *OffenseActivityDTO `json:"offense,omitempty"`
	Defense *DefenseActivityDTO `json:"defense,omitempty"`
	Scan    *ScanActivityDTO    `json:"scan,omitempty"`
	Radar   *RadarActivityDTO   `json:"radar,omitempty"`
}

// OffenseActivityDTO mirrors readmodels.OffenseActivity.
type OffenseActivityDTO struct {
	OpID              int                    `json:"opId"`
	Subtype           OffenseActivitySubtype `json:"subtype"`
	Operation         *MilitaryOperationDTO  `json:"operation,omitempty"`
	PriorOpponentScan *SectorDTO             `json:"priorOpponentScan,omitempty"`
}

// DefenseActivityDTO mirrors readmodels.DefenseActivity.
type DefenseActivityDTO struct {
	OpID              int                    `json:"opId"`
	Subtype           DefenseActivitySubtype `json:"subtype"`
	Operation         *MilitaryOperationDTO  `json:"operation,omitempty"`
	PriorOpponentScan *SectorDTO             `json:"priorOpponentScan,omitempty"`
}

// ScanActivityDTO mirrors readmodels.ScanActivity, embedding the full scan report readmodel.
type ScanActivityDTO struct {
	ReportID int        `json:"reportId"`
	Report   *SectorDTO `json:"report,omitempty"`
}

// RadarActivityDTO presents a link to a stateful radar threat.
type RadarActivityDTO struct {
	ThreatID string          `json:"threatId"`
	Threat   *RadarThreatDTO `json:"threat,omitempty"`
}

func offenseActivityFromReadModel(offense *readmodels.OffenseActivity) *OffenseActivityDTO {
	if offense == nil {
		return nil
	}
	dto := &OffenseActivityDTO{
		OpID:    offense.OpID,
		Subtype: OffenseActivitySubtype(offense.Subtype),
	}
	if offense.Operation != nil {
		m := OperationFromReadModel(offense.Operation)
		dto.Operation = &m
	}
	if offense.PriorOpponentScan != nil {
		report := SectorScanReportFromReadModel(offense.PriorOpponentScan)
		dto.PriorOpponentScan = &report
	}
	return dto
}

func defenseActivityFromReadModel(defense *readmodels.DefenseActivity) *DefenseActivityDTO {
	if defense == nil {
		return nil
	}
	dto := &DefenseActivityDTO{
		OpID:    defense.OpID,
		Subtype: DefenseActivitySubtype(defense.Subtype),
	}
	if defense.Operation != nil {
		m := OperationFromReadModel(defense.Operation)
		dto.Operation = &m
	}
	if defense.PriorOpponentScan != nil {
		report := SectorScanReportFromReadModel(defense.PriorOpponentScan)
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
	dto := &RadarActivityDTO{
		ThreatID: radar.ThreatID.String(),
	}
	if radar.Threat != nil {
		t := RadarThreatFromReadModel(radar.Threat)
		dto.Threat = &t
	}
	return dto
}

func ActivityItemDTOFromReadModel(a *readmodels.ActivityItem) ActivityItemDTO {
	dto := ActivityItemDTO{
		ID:        a.ID,
		Kind:      ActivityKind(a.Kind),
		CreatedAt: a.CreatedAt,
		BaseID:    a.BaseID,
		Offense:   offenseActivityFromReadModel(a.Offense),
		Defense:   defenseActivityFromReadModel(a.Defense),
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
