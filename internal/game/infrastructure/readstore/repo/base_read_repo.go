package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/readstore/mappers"
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

func (r *BaseReadRepo) ListUserBases(userID int) ([]*readmodels.UserBaseModel, error) {
	rows, err := r.q.ListUserBases(context.Background(), int64(userID))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.UserBaseModel, 0, len(rows))
	for _, row := range rows {
		v := mappers.UserBaseFromBasicRow(row)
		out = append(out, &v)
	}
	return out, nil
}
