package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
)

type SectorType string

// SectorType enum values
const (
	Home        SectorType = "HOME"
	Unknown     SectorType = "UNKNOWN"
	UserBase    SectorType = "BASE"
	Empty       SectorType = "EMPTY"
	Resourceful SectorType = "RESOURCEFUL"
	Dangerous   SectorType = "DANGEROUS"
)

type ScanSource string

const (
	ScanSourceUnknown   ScanSource = "UNKNOWN"
	ScanSourceOperation ScanSource = "OPERATION"
	ScanSourceScanner   ScanSource = "SCANNER"
	ScanSourceIntel     ScanSource = "INTEL"
)

type SectorDTO struct {
	Coordinates  Vector2iDTO  `json:"coordinates"`
	Type         SectorType   `json:"type"`
	Name         string       `json:"name"`
	Description  string       `json:"description"`
	ImageURL     string       `json:"image_url"`
	ScanDate     int          `json:"scan_date"`
	ScanReportID int          `json:"scan_report_id"`
	ScanInfo     *ScanInfoDTO `json:"scan_info"`
	Source       ScanSource   `json:"source"`
}

type ScanInfoDTO struct {
	Credits    int `json:"credits"`
	Iron       int `json:"iron"`
	Titanium   int `json:"titanium"`
	Antimatter int `json:"antimatter"`
	Defence    int `json:"defence"`
	Attack     int `json:"attack"`
	Space      int `json:"space"`
}

func scanInfoFromReadModel(info readmodels.ScanInfo) *ScanInfoDTO {
	return &ScanInfoDTO{
		Credits:    info.Credits,
		Iron:       info.Iron,
		Titanium:   info.Titanium,
		Antimatter: info.Antimatter,
		Defence:    info.Defence,
		Attack:     info.Attack,
		Space:      info.Space,
	}
}

func sectorTypeFromLocation(loc readmodels.LocationType) SectorType {
	switch loc {
	case readmodels.LocationTypeUserBase:
		return UserBase
	case readmodels.LocationTypeResourceful:
		return Resourceful
	case readmodels.LocationTypeDangerous:
		return Dangerous
	case readmodels.LocationTypeEmpty:
		return Empty
	default:
		return Unknown
	}
}

func SectorScanReportFromReadModel(r *readmodels.SectorScanReport, tr ports.Translator, locale string) SectorDTO {
	source := ScanSourceUnknown
	if r.SourceIntelItemID != nil {
		source = ScanSourceIntel
	} else if r.SourceScannerID != nil {
		source = ScanSourceScanner
	} else if r.SourceOperationID > 0 {
		source = ScanSourceOperation
	}

	return SectorDTO{
		Coordinates:  Vector2iFromReadModel(r.Coordinates),
		Type:         sectorTypeFromLocation(r.Type),
		Name:         tr.T(locale, r.Details.Name, nil),
		Description:  tr.T(locale, r.Details.Description, nil),
		ImageURL:     r.Details.ImageURL,
		ScanDate:     int(r.CreatedAt),
		ScanReportID: r.ID,
		ScanInfo:     scanInfoFromReadModel(r.Info),
		Source:       source,
	}
}

func SectorScanReportsFromReadModels(reports []*readmodels.SectorScanReport, tr ports.Translator, locale string) []SectorDTO {
	out := make([]SectorDTO, 0, len(reports))
	for _, r := range reports {
		out = append(out, SectorScanReportFromReadModel(r, tr, locale))
	}
	return out
}

func Vector2iListFromReadModels(coords []readmodels.Vector2i) []Vector2iDTO {
	out := make([]Vector2iDTO, 0, len(coords))
	for _, coord := range coords {
		out = append(out, Vector2iFromReadModel(coord))
	}
	return out
}
