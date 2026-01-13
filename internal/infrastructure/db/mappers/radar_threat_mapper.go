package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
)

func RadarThreatFromModel(m gen.RadarThreat) *domain.RadarThreat {
	return &domain.RadarThreat{
		ID:                 m.ID,
		OperationID:        int(m.OperationID),
		OwnerBaseID:        int(m.OwnerBaseID),
		DetectedAt:         m.DetectedAt,
		SourceCoordinates:  domain.Vector2i{X: int(m.SourceX), Y: int(m.SourceY)},
		TargetCoordinates:  domain.Vector2i{X: int(m.TargetX), Y: int(m.TargetY)},
		EstimatedArrivalAt: m.EstimatedArrivalAt,
		ArrivalAt:          nullInt64ToInt64Ptr(m.ArrivalAt),
		Type:               domain.ThreatType(m.Type),
		Status:             domain.ThreatStatus(m.Status),
		Attack:             int(m.Attack),
		Speed:              int(m.Speed),
		Stealth:            int(m.Stealth),
		Capacity:           int(m.Capacity),
	}
}

func InsertRadarThreatParamsFromDomain(t *domain.RadarThreat) gen.InsertRadarThreatParams {
	return gen.InsertRadarThreatParams{
		ID:                 t.ID,
		OperationID:        int64(t.OperationID),
		OwnerBaseID:        int64(t.OwnerBaseID),
		DetectedAt:         t.DetectedAt,
		SourceX:            int32(t.SourceCoordinates.X),
		SourceY:            int32(t.SourceCoordinates.Y),
		TargetX:            int32(t.TargetCoordinates.X),
		TargetY:            int32(t.TargetCoordinates.Y),
		EstimatedArrivalAt: t.EstimatedArrivalAt,
		ArrivalAt:          int64PtrToNullInt64(t.ArrivalAt),
		Type:               string(t.Type),
		Status:             string(t.Status),
		Attack:             int32(t.Attack),
		Speed:              int32(t.Speed),
		Stealth:            int32(t.Stealth),
		Capacity:           int32(t.Capacity),
	}
}

func UpdateRadarThreatParamsFromDomain(t *domain.RadarThreat) gen.UpdateRadarThreatParams {
	return gen.UpdateRadarThreatParams{
		ID:                 t.ID,
		EstimatedArrivalAt: t.EstimatedArrivalAt,
		ArrivalAt:          int64PtrToNullInt64(t.ArrivalAt),
		Status:             string(t.Status),
	}
}
