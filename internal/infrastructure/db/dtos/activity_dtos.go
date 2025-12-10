package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

// OperationActivityDTO mirrors the JSON structure stored for operation activities.
type OperationActivityDTO struct {
	OpID    int    `json:"op_id"`
	Subtype string `json:"subtype"`
	Role    string `json:"role"`
}

func OperationActivityDTOFromDomain(o *domain.OperationActivity) *OperationActivityDTO {
	if o == nil {
		return nil
	}
	return &OperationActivityDTO{OpID: o.OpID, Subtype: string(o.Subtype), Role: string(o.Role)}
}
func OperationActivityFromDTO(d *OperationActivityDTO) *domain.OperationActivity {
	if d == nil {
		return nil
	}
	return &domain.OperationActivity{OpID: d.OpID, Subtype: domain.MilitaryActivitySubtype(d.Subtype), Role: domain.OperationRole(d.Role)}
}

// ScanActivityDTO mirrors the JSON structure stored for scan activities.
type ScanActivityDTO struct {
	ReportID int `json:"report_id"`
}

func ScanActivityDTOFromDomain(s *domain.ScanActivity) *ScanActivityDTO {
	if s == nil {
		return nil
	}
	return &ScanActivityDTO{ReportID: s.ReportID}
}
func ScanActivityFromDTO(d *ScanActivityDTO) *domain.ScanActivity {
	if d == nil {
		return nil
	}
	return &domain.ScanActivity{ReportID: d.ReportID}
}

// RadarActivityDTO mirrors the JSON structure stored for radar activities.
type RadarActivityDTO struct {
	OpID       int   `json:"op_id"`
	DetectedAt int64 `json:"detected_at"`
	EtaAtBase  int64 `json:"eta_at_base"`
	SourceX    int   `json:"source_x"`
	SourceY    int   `json:"source_y"`
	TargetX    int   `json:"target_x"`
	TargetY    int   `json:"target_y"`
	Threat     struct {
		Attack  int `json:"attack"`
		Defence int `json:"defence"`
	} `json:"threat"`
}

func RadarActivityDTOFromDomain(r *domain.RadarActivity) *RadarActivityDTO {
	if r == nil {
		return nil
	}
	dto := &RadarActivityDTO{OpID: r.OpID, DetectedAt: r.DetectedAt, EtaAtBase: r.EtaAtBase, SourceX: r.SourceCoordinates.X, SourceY: r.SourceCoordinates.Y, TargetX: r.TargetCoordinates.X, TargetY: r.TargetCoordinates.Y}
	dto.Threat.Attack = r.Threat.Attack
	dto.Threat.Defence = r.Threat.Defence
	return dto
}
func RadarActivityFromDTO(d *RadarActivityDTO) *domain.RadarActivity {
	if d == nil {
		return nil
	}
	return &domain.RadarActivity{
		OpID:              d.OpID,
		DetectedAt:        d.DetectedAt,
		EtaAtBase:         d.EtaAtBase,
		SourceCoordinates: domain.Vector2i{X: d.SourceX, Y: d.SourceY},
		TargetCoordinates: domain.Vector2i{X: d.TargetX, Y: d.TargetY},
		Threat:            domain.Threat{Attack: d.Threat.Attack, Defence: d.Threat.Defence},
	}
}
