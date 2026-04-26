package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
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
	Subtype   string            `json:"subtype"`
	ReportID  *int              `json:"report_id,omitempty"`
	Intercept *ScanInterceptDTO `json:"intercept,omitempty"`
}

type ScanInterceptDTO struct {
	ScannedCoordinates     Vector2iDTO  `json:"scanned_coordinates"`
	ScanPenetratedCloaking bool         `json:"scan_penetrated_cloaking"`
	PossibleSource         *Vector2iDTO `json:"possible_source,omitempty"`
	UncertaintyRadius      int          `json:"uncertainty_radius"`
}

func ScanActivityDTOFromDomain(s *domain.ScanActivity) *ScanActivityDTO {
	if s == nil {
		return nil
	}
	dto := &ScanActivityDTO{
		Subtype:  string(s.Subtype),
		ReportID: s.ReportID,
	}
	if s.Intercept != nil {
		dto.Intercept = &ScanInterceptDTO{
			ScannedCoordinates:     Vector2iDTOFromDomain(s.Intercept.ScannedCoordinates),
			ScanPenetratedCloaking: s.Intercept.ScanPenetratedCloaking,
			UncertaintyRadius:      s.Intercept.UncertaintyRadius,
		}
		if s.Intercept.PossibleSource != nil {
			src := Vector2iDTOFromDomain(*s.Intercept.PossibleSource)
			dto.Intercept.PossibleSource = &src
		}
	}
	return dto
}

func ScanActivityFromDTO(d *ScanActivityDTO) *domain.ScanActivity {
	if d == nil {
		return nil
	}
	res := &domain.ScanActivity{
		Subtype:  domain.ScanActivitySubtype(d.Subtype),
		ReportID: d.ReportID,
	}
	if d.Intercept != nil {
		res.Intercept = &domain.ScanInterceptInfo{
			ScannedCoordinates:     d.Intercept.ScannedCoordinates.ToDomain(),
			ScanPenetratedCloaking: d.Intercept.ScanPenetratedCloaking,
			PossibleSource: func() *domain.Vector2i {
				if d.Intercept.PossibleSource == nil {
					return nil
				}
				v := d.Intercept.PossibleSource.ToDomain()
				return &v
			}(),
			UncertaintyRadius: d.Intercept.UncertaintyRadius,
		}
	}
	return res
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

// TradeActivityDTO mirrors the JSON structure stored for trade activities.
type TradeActivityDTO struct {
	OpID int `json:"op_id"`
}

func TradeActivityDTOFromDomain(t *domain.TradeActivity) *TradeActivityDTO {
	if t == nil {
		return nil
	}
	return &TradeActivityDTO{OpID: t.OpID}
}

func TradeActivityFromDTO(d *TradeActivityDTO) *domain.TradeActivity {
	if d == nil {
		return nil
	}
	return &domain.TradeActivity{OpID: d.OpID}
}
