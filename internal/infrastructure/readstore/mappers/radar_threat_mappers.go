package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
)

func RadarThreatFromModel(m gen.RadarThreat) *readmodels.RadarThreat {
	return &readmodels.RadarThreat{
		ID:                 m.ID,
		OperationID:        int(m.OperationID),
		OwnerBaseID:        int(m.OwnerBaseID),
		DetectedAt:         m.DetectedAt,
		SourceCoordinates:  readmodels.Vector2i{X: int(m.SourceX), Y: int(m.SourceY)},
		TargetCoordinates:  readmodels.Vector2i{X: int(m.TargetX), Y: int(m.TargetY)},
		EstimatedArrivalAt: m.EstimatedArrivalAt,
		ArrivalAt:          nullInt64ToInt64Ptr(m.ArrivalAt),
		Type:               readmodels.ThreatType(m.Type),
		Status:             readmodels.ThreatStatus(m.Status),
		Attack:             int(m.Attack),
		Speed:              int(m.Speed),
		Stealth:            int(m.Stealth),
		Capacity:           int(m.Capacity),
	}
}
