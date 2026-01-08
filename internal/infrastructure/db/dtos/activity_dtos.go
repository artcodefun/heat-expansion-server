package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

// OffenseActivityDTO mirrors the JSON structure stored for offensive activities.
type OffenseActivityDTO struct {
	OpID    int    `json:"op_id"`
	Subtype string `json:"subtype"`
}

func OffenseActivityDTOFromDomain(o *domain.OffenseActivity) *OffenseActivityDTO {
	if o == nil {
		return nil
	}
	return &OffenseActivityDTO{OpID: o.OpID, Subtype: string(o.Subtype)}
}
func OffenseActivityFromDTO(d *OffenseActivityDTO) *domain.OffenseActivity {
	if d == nil {
		return nil
	}
	return &domain.OffenseActivity{OpID: d.OpID, Subtype: domain.OffenseActivitySubtype(d.Subtype)}
}

// DefenseActivityDTO mirrors the JSON structure stored for defensive activities.
type DefenseActivityDTO struct {
	OpID    int    `json:"op_id"`
	Subtype string `json:"subtype"`
}

func DefenseActivityDTOFromDomain(o *domain.DefenseActivity) *DefenseActivityDTO {
	if o == nil {
		return nil
	}
	return &DefenseActivityDTO{OpID: o.OpID, Subtype: string(o.Subtype)}
}
func DefenseActivityFromDTO(d *DefenseActivityDTO) *domain.DefenseActivity {
	if d == nil {
		return nil
	}
	return &domain.DefenseActivity{OpID: d.OpID, Subtype: domain.DefenseActivitySubtype(d.Subtype)}
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
		Attack   int `json:"attack"`
		Speed    int `json:"speed"`
		Stealth  int `json:"stealth"`
		Capacity int `json:"capacity"`
	} `json:"threat"`
}

func RadarActivityDTOFromDomain(r *domain.RadarActivity) *RadarActivityDTO {
	if r == nil {
		return nil
	}
	dto := &RadarActivityDTO{OpID: r.OpID, DetectedAt: r.DetectedAt, EtaAtBase: r.EtaAtBase, SourceX: r.SourceCoordinates.X, SourceY: r.SourceCoordinates.Y, TargetX: r.TargetCoordinates.X, TargetY: r.TargetCoordinates.Y}
	dto.Threat.Attack = r.Threat.Attack
	dto.Threat.Speed = r.Threat.Speed
	dto.Threat.Stealth = r.Threat.Stealth
	dto.Threat.Capacity = r.Threat.Capacity
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
		Threat: domain.Threat{
			Attack:   d.Threat.Attack,
			Speed:    d.Threat.Speed,
			Stealth:  d.Threat.Stealth,
			Capacity: d.Threat.Capacity,
		},
	}
}
