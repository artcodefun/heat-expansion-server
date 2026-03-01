package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
)

type SectorReadRepo struct{ q *gen.Queries }

func NewSectorReadRepo(q *gen.Queries) *SectorReadRepo { return &SectorReadRepo{q: q} }

func (r *SectorReadRepo) GetScansNear(ctx context.Context, baseID int, x, y, radius int) ([]*readmodels.SectorScanReport, error) {
	rows, err := r.q.GetScansNear(ctx, gen.GetScansNearParams{BaseID: int64(baseID), SectorX: int32(x), SectorY: int32(y), Column4: int32(radius)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.SectorScanReport, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.SectorScanReportFromModel(r0)
		out = append(out, &v)
	}
	return out, nil
}

func (r *SectorReadRepo) GetScanReportByID(ctx context.Context, baseID, id int) (*readmodels.SectorScanReport, error) {
	row, err := r.q.GetScanReportByID(ctx, gen.GetScanReportByIDParams{ID: int64(id), BaseID: int64(baseID)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	v := mappers.SectorScanReportFromModel(row)
	return &v, nil
}

func (r *SectorReadRepo) GetLatestScanBefore(ctx context.Context, baseID, x, y int, before int64) (*readmodels.SectorScanReport, error) {
	row, err := r.q.GetLatestScanBefore(ctx, gen.GetLatestScanBeforeParams{BaseID: int64(baseID), SectorX: int32(x), SectorY: int32(y), CreatedAt: before})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	v := mappers.SectorScanReportFromModel(row)
	return &v, nil
}
