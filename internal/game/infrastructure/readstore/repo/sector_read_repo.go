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

type SectorReadRepo struct {
	q     *gen.Queries
	bases ports.BaseReadRepository
}

func NewSectorReadRepo(q *gen.Queries, bases ports.BaseReadRepository) *SectorReadRepo {
	return &SectorReadRepo{q: q, bases: bases}
}

func (r *SectorReadRepo) GetScansNear(ctx context.Context, baseID int, x, y, radius int) ([]*readmodels.SectorScanReport, error) {
	rows, err := r.q.GetScansNear(ctx, gen.GetScansNearParams{BaseID: int64(baseID), SectorX: int32(x), SectorY: int32(y), Column4: int32(radius)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.SectorScanReport, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.SectorScanReportFromModel(r0)
		if err := r.enrichOwnerUserID(ctx, &v); err != nil {
			return nil, err
		}
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
	if err := r.enrichOwnerUserID(ctx, &v); err != nil {
		return nil, err
	}
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
	if err := r.enrichOwnerUserID(ctx, &v); err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *SectorReadRepo) enrichOwnerUserID(ctx context.Context, report *readmodels.SectorScanReport) error {
	if report == nil || report.Type != readmodels.LocationTypeUserBase {
		return nil
	}
	owner, err := r.bases.GetBaseOwnerByCoordinates(ctx, report.Coordinates.X, report.Coordinates.Y)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return nil
		}
		return err
	}
	report.Owner = owner
	return nil
}
