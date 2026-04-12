package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
	"github.com/google/uuid"
)

type BaseReadRepo struct{ q *gen.Queries }

func NewBaseReadRepo(q *gen.Queries) *BaseReadRepo { return &BaseReadRepo{q: q} }

func (r *BaseReadRepo) GetBase(ctx context.Context, baseID int) (*readmodels.UserBaseModel, error) {
	row, err := r.q.GetBase(ctx, int64(baseID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	v := mappers.UserBaseFromGetRow(row)
	return &v, nil
}

func (r *BaseReadRepo) GetBaseOwnerByCoordinates(ctx context.Context, x, y int) (*readmodels.SectorOwner, error) {
	row, err := r.q.GetBaseOwnerByCoordinates(ctx, gen.GetBaseOwnerByCoordinatesParams{SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	v := readmodels.SectorOwner{ID: row.ID, Name: row.Name}
	return &v, nil
}

func (r *BaseReadRepo) GetBaseStats(ctx context.Context, baseID int) (*readmodels.UserBaseStats, error) {
	row, err := r.q.GetBaseStats(ctx, int64(baseID))
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

func (r *BaseReadRepo) ListUserBases(ctx context.Context, userID uuid.UUID) ([]*readmodels.UserBaseModel, error) {
	rows, err := r.q.ListUserBases(ctx, userID)
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
