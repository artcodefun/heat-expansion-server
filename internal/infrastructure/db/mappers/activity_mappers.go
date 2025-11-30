package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
	"github.com/sqlc-dev/pqtype"
)

// InsertActivityParamsFromDomain maps a domain.ActivityItem into sqlc params for insert.
func InsertActivityParamsFromDomain(a *domain.ActivityItem) gen.InsertActivityParams {
	var op any
	if a.Operation != nil {
		dto := dtos.OperationActivityDTOFromDomain(*a.Operation)
		op = dto
	}

	var scan any
	if a.Scan != nil {
		dto := dtos.ScanActivityDTOFromDomain(*a.Scan)
		scan = dto
	}

	var radar any
	if a.Radar != nil {
		dto := dtos.RadarActivityDTOFromDomain(*a.Radar)
		radar = dto
	}

	return gen.InsertActivityParams{
		Kind:          string(a.Kind),
		CreatedAt:     a.CreatedAt,
		BaseID:        int64(a.BaseID),
		OperationData: toNullRawMessage(op),
		ScanData:      toNullRawMessage(scan),
		RadarData:     toNullRawMessage(radar),
		TradeData:     pqtype.NullRawMessage{Valid: false},
	}
}
