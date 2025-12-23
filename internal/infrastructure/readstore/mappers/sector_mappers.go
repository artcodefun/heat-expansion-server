package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
)

// SectorScanReportFromModel converts a generic scan report row returned by scan-report queries.
func SectorScanReportFromModel(r gen.ScanReport) readmodels.SectorScanReport {
	info := parseScanInfo(r.Info)
	return readmodels.SectorScanReport{ID: int(r.ID), BaseID: int(r.BaseID), Coordinates: readmodels.Vector2i{X: int(r.SectorX), Y: int(r.SectorY)}, CreatedAt: r.CreatedAt, Details: readmodels.LocationDetails{Name: nullString(r.Name), Description: nullString(r.Description), ImageURL: nullString(r.ImageUrl)}, Type: readmodels.LocationType(r.Type), Info: info, IsCloaked: r.IsCloaked, SourceOperationID: int(nullInt64(r.SourceOperationID))}
}

// Helpers local to sector mapping
func parseScanInfo(b []byte) readmodels.ScanInfo {
	if len(b) == 0 {
		return readmodels.ScanInfo{}
	}
	var dto dtos.ScanInfoDTO
	if err := json.Unmarshal(b, &dto); err != nil {
		return readmodels.ScanInfo{}
	}
	return readmodels.ScanInfo{Credits: dto.Credits, Iron: dto.Iron, Titanium: dto.Titanium, Antimatter: dto.Antimatter, Defence: dto.Defence, Attack: dto.Attack, Space: dto.Space}
}
