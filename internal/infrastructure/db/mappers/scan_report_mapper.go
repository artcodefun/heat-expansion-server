package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
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
		Type:      domain.LocationType(r.Type),
		Info:      dtos.ScanInfoFromDTO(infoDTO),
		IsCloaked: r.IsCloaked,
	}
	if r.SourceOperationID.Valid {
		sr.SourceOperationID = int(r.SourceOperationID.Int64)
	} else {
		sr.SourceOperationID = 0
	}
	if r.SourceScannerID.Valid {
		sr.SourceScannerID = &r.SourceScannerID.UUID
	}
	if r.SourceIntelItemID.Valid {
		sr.SourceIntelItemID = &r.SourceIntelItemID.UUID
	}
	return sr
}

func InsertScanReportParamsFromDomain(r *domain.SectorScanReport) gen.InsertScanReportParams {
	infoJSON, _ := json.Marshal(dtos.ScanInfoDTOFromDomain(r.Info))
	var srcOpID = toNullInt64ZeroAsNull(r.SourceOperationID)
	var srcScannerID uuid.NullUUID
	if r.SourceScannerID != nil {
		srcScannerID = uuid.NullUUID{UUID: *r.SourceScannerID, Valid: true}
	}
	var srcIntelID uuid.NullUUID
	if r.SourceIntelItemID != nil {
		srcIntelID = uuid.NullUUID{UUID: *r.SourceIntelItemID, Valid: true}
	}
	return gen.InsertScanReportParams{
		BaseID:            int64(r.BaseID),
		SectorX:           int32(r.Coordinates.X),
		SectorY:           int32(r.Coordinates.Y),
		CreatedAt:         r.CreatedAt,
		Type:              string(r.Type),
		IsCloaked:         r.IsCloaked,
		SourceOperationID: srcOpID,
		SourceScannerID:   srcScannerID,
		SourceIntelItemID: srcIntelID,
		Name:              toNullString(r.Details.Name),
		Description:       toNullString(r.Details.Description),
		ImageUrl:          toNullString(r.Details.ImageURL),
		Info:              infoJSON,
	}
}
