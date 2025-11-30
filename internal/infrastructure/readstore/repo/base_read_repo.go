package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/mappers"
)

type BaseReadRepo struct{ q *gen.Queries }

func NewBaseReadRepo(q *gen.Queries) *BaseReadRepo { return &BaseReadRepo{q: q} }

func (r *BaseReadRepo) GetBaseStats(baseID int) (*readmodels.UserBaseStats, error) {
	row, err := r.q.GetBaseStats(context.Background(), int64(baseID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	var dto dtos.BaseStatsDTO
	if err := json.Unmarshal(row.Stats, &dto); err != nil {
		return nil, err
	}
	domainStats := mappers.UserBaseStatsFromDTO(dto, row.StatsCalcTimestamp)
	return &domainStats, nil
}
