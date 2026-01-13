package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
)

func ActivityItemFromModel(a gen.Activity) readmodels.ActivityItem {
	item := readmodels.ActivityItem{ID: int(a.ID), Kind: readmodels.ActivityKind(a.Kind), CreatedAt: a.CreatedAt, BaseID: int(a.BaseID)}
	if a.OffenseData.Valid {
		var dto dtos.OffenseActivityDTO
		_ = json.Unmarshal(a.OffenseData.RawMessage, &dto)
		item.Offense = &readmodels.OffenseActivity{OpID: dto.OpID, Subtype: readmodels.OffenseActivitySubtype(dto.Subtype)}
	}
	if a.DefenseData.Valid {
		var dto dtos.DefenseActivityDTO
		_ = json.Unmarshal(a.DefenseData.RawMessage, &dto)
		item.Defense = &readmodels.DefenseActivity{OpID: dto.OpID, Subtype: readmodels.DefenseActivitySubtype(dto.Subtype)}
	}
	if a.ScanData.Valid {
		var dto dtos.ScanActivityDTO
		_ = json.Unmarshal(a.ScanData.RawMessage, &dto)
		item.Scan = &readmodels.ScanActivity{ReportID: dto.ReportID}
	}
	if a.RadarData.Valid {
		var dto dtos.RadarActivityDTO
		_ = json.Unmarshal(a.RadarData.RawMessage, &dto)
		item.Radar = &readmodels.RadarActivity{ThreatID: dto.ThreatID}
	}
	return item
}
