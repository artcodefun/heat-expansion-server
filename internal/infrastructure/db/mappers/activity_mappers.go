package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
	"github.com/sqlc-dev/pqtype"
)

// InsertActivityParamsFromDomain maps a domain.ActivityItem into sqlc params for insert.
func InsertActivityParamsFromDomain(a *domain.ActivityItem) gen.InsertActivityParams {
	offenseDTO := dtos.OffenseActivityDTOFromDomain(a.Offense)
	defenseDTO := dtos.DefenseActivityDTOFromDomain(a.Defense)
	scanDTO := dtos.ScanActivityDTOFromDomain(a.Scan)
	radarDTO := dtos.RadarActivityDTOFromDomain(a.Radar)

	return gen.InsertActivityParams{
		ID:          a.ID,
		Kind:        string(a.Kind),
		CreatedAt:   a.CreatedAt,
		BaseID:      int64(a.BaseID),
		OffenseData: toNullRawMessage(offenseDTO),
		DefenseData: toNullRawMessage(defenseDTO),
		ScanData:    toNullRawMessage(scanDTO),
		RadarData:   toNullRawMessage(radarDTO),
		TradeData:   pqtype.NullRawMessage{Valid: false},
	}
}
