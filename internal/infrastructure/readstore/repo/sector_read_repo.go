package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/mappers"
)

type SectorReadRepo struct{ q *gen.Queries }

func NewSectorReadRepo(q *gen.Queries) *SectorReadRepo { return &SectorReadRepo{q: q} }

func (r *SectorReadRepo) GetSector(x, y int) (*readmodels.SectorModel, error) {
	row, err := r.q.GetSector(context.Background(), gen.GetSectorParams{X: int32(x), Y: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	model := mappers.SectorModelFromRow(row)
	return &model, nil
}

func (r *SectorReadRepo) GetLatestScans(baseID int) ([]*readmodels.SectorScanReport, error) {
	rows, err := r.q.GetLatestScans(context.Background(), int64(baseID))
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

func (r *SectorReadRepo) GetScansNear(baseID int, x, y, radius int) ([]*readmodels.SectorScanReport, error) {
	rows, err := r.q.GetScansNear(context.Background(), gen.GetScansNearParams{BaseID: int64(baseID), SectorX: int32(x), SectorY: int32(y), Column4: int32(radius)})
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

// Merged map-related methods
func (r *SectorReadRepo) ListOccupiedCoordinates() ([]readmodels.Vector2i, error) {
	rows, err := r.q.ListOccupiedCoordinates(context.Background())
	if err != nil {
		return nil, err
	}
	out := make([]readmodels.Vector2i, 0, len(rows))
	for _, r0 := range rows {
		out = append(out, mappers.Vector2iFromOccupiedRow(r0))
	}
	return out, nil
}

func (r *SectorReadRepo) ListSectorsInRadius(x, y, radius int) ([]*readmodels.SectorModel, error) {
	rows, err := r.q.ListSectorsInRadius(context.Background(), gen.ListSectorsInRadiusParams{X: int32(x), Y: int32(y), Column3: int32(radius)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.SectorModel, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.SectorModelFromRadiusRow(r0)
		out = append(out, &v)
	}
	return out, nil
}
