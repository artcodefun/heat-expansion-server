package dtos

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/google/uuid"
)

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
	ThreatID uuid.UUID `json:"threat_id"`
}

func RadarActivityDTOFromDomain(r *domain.RadarActivity) *RadarActivityDTO {
	if r == nil {
		return nil
	}
	return &RadarActivityDTO{ThreatID: r.ThreatID}
}
func RadarActivityFromDTO(d *RadarActivityDTO) *domain.RadarActivity {
	if d == nil {
		return nil
	}
	return &domain.RadarActivity{
		ThreatID: d.ThreatID,
	}
}
