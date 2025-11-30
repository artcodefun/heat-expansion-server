package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
)

func ActivityItemFromModel(a gen.Activity) readmodels.ActivityItem {
	item := readmodels.ActivityItem{ID: int(a.ID), Kind: readmodels.ActivityKind(a.Kind), CreatedAt: a.CreatedAt, BaseID: int(a.BaseID)}
	if a.OperationData.Valid {
		var dto dtos.OperationActivityDTO
		_ = json.Unmarshal(a.OperationData.RawMessage, &dto)
		item.Operation = &readmodels.OperationActivity{OpID: dto.OpID, Subtype: readmodels.MilitaryActivitySubtype(dto.Subtype), Role: readmodels.OperationRole(dto.Role)}
	}
	if a.ScanData.Valid {
		var dto dtos.ScanActivityDTO
		_ = json.Unmarshal(a.ScanData.RawMessage, &dto)
		item.Scan = &readmodels.ScanActivity{ReportID: dto.ReportID}
	}
	if a.RadarData.Valid {
		var dto dtos.RadarActivityDTO
		_ = json.Unmarshal(a.RadarData.RawMessage, &dto)
		item.Radar = &readmodels.RadarActivity{OpID: dto.OpID, DetectedAt: dto.DetectedAt, EtaAtBase: dto.EtaAtBase, SourceCoordinates: readmodels.Vector2i{X: dto.SourceX, Y: dto.SourceY}, TargetCoordinates: readmodels.Vector2i{X: dto.TargetX, Y: dto.TargetY}, Threat: readmodels.Threat{Attack: dto.Threat.Attack, Defence: dto.Threat.Defence}}
	}
	return item
}
