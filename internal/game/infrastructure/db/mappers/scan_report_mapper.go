package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/google/uuid"
)

func ScanReportFromDB(r gen.ScanReport) *domain.SectorScanReport {
	var infoDTO dtos.ScanInfoDTO
	_ = json.Unmarshal(r.Info, &infoDTO)

	sr := &domain.SectorScanReport{
		ID:          int(r.ID),
		BaseID:      int(r.BaseID),
		Coordinates: domain.Vector2i{X: int(r.SectorX), Y: int(r.SectorY)},
		CreatedAt:   r.CreatedAt,
		Details: domain.LocationDetails{
			Name:        nullStringToString(&r.Name.String, r.Name.Valid),
			Description: nullStringToString(&r.Description.String, r.Description.Valid),
			ImageURL:    nullStringToString(&r.ImageUrl.String, r.ImageUrl.Valid),
		},
		Type:       domain.LocationType(r.Type),
		Info:       dtos.ScanInfoFromDTO(infoDTO),
		IsCloaked:  r.IsCloaked,
		SourceType: domain.ScanReportSourceType(r.SourceType),
	}
	if r.SourceID.Valid {
		sr.SourceID = &r.SourceID.UUID
	}
	return sr
}

func InsertScanReportParamsFromDomain(r *domain.SectorScanReport) gen.InsertScanReportParams {
	infoJSON, _ := json.Marshal(dtos.ScanInfoDTOFromDomain(r.Info))
	var srcID uuid.NullUUID
	if r.SourceID != nil {
		srcID = uuid.NullUUID{UUID: *r.SourceID, Valid: true}
	}
	return gen.InsertScanReportParams{
		BaseID:      int64(r.BaseID),
		SectorX:     int32(r.Coordinates.X),
		SectorY:     int32(r.Coordinates.Y),
		CreatedAt:   r.CreatedAt,
		Type:        string(r.Type),
		IsCloaked:   r.IsCloaked,
		SourceType:  string(r.SourceType),
		SourceID:    srcID,
		Name:        toNullString(r.Details.Name),
		Description: toNullString(r.Details.Description),
		ImageUrl:    toNullString(r.Details.ImageURL),
		Info:        infoJSON,
	}
}
