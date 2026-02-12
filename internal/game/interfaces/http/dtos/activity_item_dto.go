package dtos

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/application/cqrs/readmodels"
	"github.com/google/uuid"
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
	ID        uuid.UUID    `json:"id"`
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
	OpID      int                    `json:"opId"`
	Subtype   OffenseActivitySubtype `json:"subtype"`
	Operation *MilitaryOperationDTO  `json:"operation,omitempty"`
}

// OffenderInfoDTO provides a restricted view of an attacking operation for the defender.
type OffenderInfoDTO struct {
	Type              OperationType        `json:"type"`
	SourceCoordinates Vector2iDTO          `json:"sourceCoordinates"`
	TargetCoordinates Vector2iDTO          `json:"targetCoordinates"`
	ContactDate       int64                `json:"contactDate"`
	Result            OperationResult      `json:"result"`
	Units             []MilitaryUnitDTO    `json:"units"`
	StorageSnaps      []StorageItemSnapDTO `json:"storageSnaps"`
	TotalModifiers    MilitaryModifiersDTO `json:"totalModifiers"`
	SpyResult         *SpyResultDTO        `json:"spyResult,omitempty"`
	AttackResult      *AttackResultDTO     `json:"attackResult,omitempty"`
}

// DefenseActivityDTO mirrors readmodels.DefenseActivity.
type DefenseActivityDTO struct {
	OpID              int                    `json:"opId"`
	Subtype           DefenseActivitySubtype `json:"subtype"`
	Offender          *OffenderInfoDTO       `json:"offender,omitempty"`
	PriorOpponentScan *SectorDTO             `json:"priorOpponentScan,omitempty"`
}

// ScanActivitySubtype mirrors readmodels.ScanActivitySubtype.
type ScanActivitySubtype string

const (
	ScanActivitySubtypeReportProduced       ScanActivitySubtype = "REPORT_PRODUCED"
	ScanActivitySubtypeExternalScanDetected ScanActivitySubtype = "EXTERNAL_SCAN_DETECTED"
)

// ScanInterceptInfoDTO mirrors readmodels.ScanInterceptInfo.
type ScanInterceptInfoDTO struct {
	ScannedCoordinates     Vector2iDTO  `json:"scannedCoordinates"`
	ScanPenetratedCloaking bool         `json:"scanPenetratedCloaking"`
	PossibleSource         *Vector2iDTO `json:"possibleSource,omitempty"`
	UncertaintyRadius      int          `json:"uncertaintyRadius"`
}

// ScanActivityDTO mirrors readmodels.ScanActivity, embedding the full scan report readmodel.
type ScanActivityDTO struct {
	Subtype   ScanActivitySubtype   `json:"subtype"`
	ReportID  *int                  `json:"reportId,omitempty"`
	Intercept *ScanInterceptInfoDTO `json:"intercept,omitempty"`
	Report    *SectorDTO            `json:"report,omitempty"`
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
	return dto
}

func offenderInfoFromReadModel(o *readmodels.OffenderInfo) *OffenderInfoDTO {
	if o == nil {
		return nil
	}
	return &OffenderInfoDTO{
		Type:              OperationType(o.Type),
		SourceCoordinates: Vector2iDTO{X: o.SourceCoordinates.X, Y: o.SourceCoordinates.Y},
		TargetCoordinates: Vector2iDTO{X: o.TargetCoordinates.X, Y: o.TargetCoordinates.Y},
		ContactDate:       o.ContactDate,
		Result:            OperationResult(o.Result),
		Units:             MilitaryUnitsFromReadModel(o.Units),
		StorageSnaps:      storageItemSnapsFromReadModel(o.StorageSnaps),
		TotalModifiers:    MilitaryModifiersFromReadModel(o.TotalModifiers),
		SpyResult:         SpyResultFromReadModel(o.SpyResult),
		AttackResult:      AttackResultFromReadModel(o.AttackResult),
	}
}

func defenseActivityFromReadModel(defense *readmodels.DefenseActivity) *DefenseActivityDTO {
	if defense == nil {
		return nil
	}
	dto := &DefenseActivityDTO{
		OpID:    defense.OpID,
		Subtype: DefenseActivitySubtype(defense.Subtype),
	}
	if defense.Offender != nil {
		o := offenderInfoFromReadModel(defense.Offender)
		dto.Offender = o
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
	dto := &ScanActivityDTO{
		Subtype:  ScanActivitySubtype(scan.Subtype),
		ReportID: scan.ReportID,
	}
	if scan.Intercept != nil {
		dto.Intercept = &ScanInterceptInfoDTO{
			ScannedCoordinates:     Vector2iFromReadModel(scan.Intercept.ScannedCoordinates),
			ScanPenetratedCloaking: scan.Intercept.ScanPenetratedCloaking,
			UncertaintyRadius:      scan.Intercept.UncertaintyRadius,
		}
		if scan.Intercept.PossibleSource != nil {
			src := Vector2iFromReadModel(*scan.Intercept.PossibleSource)
			dto.Intercept.PossibleSource = &src
		}
	}
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
