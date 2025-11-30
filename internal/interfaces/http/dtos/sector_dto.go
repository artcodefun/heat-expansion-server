package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"

type SectorType string

// SectorType enum for DTOs
const (
	HomeDTO        SectorType = "HOME"
	UnknownDTO     SectorType = "UNKNOWN"
	SignalDTO      SectorType = "SIGNAL"
	UserBaseDTO    SectorType = "BASE"
	EmptyDTO       SectorType = "EMPTY"
	ResourcefulDTO SectorType = "RESOURCEFUL"
	DangerousDTO   SectorType = "DANGEROUS"
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
	if info == (readmodels.ScanInfo{}) {
		return nil
	}
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
		return UserBaseDTO
	case readmodels.LocationTypeResourceful:
		return ResourcefulDTO
	case readmodels.LocationTypeDangerous:
		return DangerousDTO
	case readmodels.LocationTypeEmpty:
		return EmptyDTO
	default:
		return UnknownDTO
	}
}

func SectorFromReadModel(m *readmodels.SectorModel) SectorDTO {
	return SectorDTO{
		Coordinates: Vector2iFromReadModel(m.Coordinates),
		Type:        UnknownDTO,
		Name:        m.Details.Name,
		Description: m.Details.Description,
		ImageURL:    m.Details.ImageURL,
	}
}

func SectorScanReportFromReadModel(r *readmodels.SectorScanReport) SectorDTO {
	return SectorDTO{
		Coordinates:  Vector2iFromReadModel(r.Coordinates),
		Type:         sectorTypeFromLocation(r.Type),
		Name:         r.Details.Name,
		Description:  r.Details.Description,
		ImageURL:     r.Details.ImageURL,
		ScanDate:     int(r.CreatedAt),
		ScanReportID: r.ID,
		ScanInfo:     scanInfoFromReadModel(r.Info),
	}
}

func SectorScanReportsFromReadModels(reports []*readmodels.SectorScanReport) []SectorDTO {
	out := make([]SectorDTO, 0, len(reports))
	for _, r := range reports {
		out = append(out, SectorScanReportFromReadModel(r))
	}
	return out
}

func SectorModelsFromReadModels(models []*readmodels.SectorModel) []SectorDTO {
	out := make([]SectorDTO, 0, len(models))
	for _, m := range models {
		out = append(out, SectorFromReadModel(m))
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
