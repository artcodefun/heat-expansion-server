package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/mappers"
)

type SectorRepo struct {
	q *gen.Queries
}

func NewSectorRepo(q *gen.Queries) *SectorRepo { return &SectorRepo{q: q} }

func (r *SectorRepo) Tx(tx ports.Transaction) ports.SectorRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &SectorRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

// Create persists a sector. Location details are stored in sector row (name/description/image_url).
func (r *SectorRepo) Create(ctx context.Context, sector *domain.SectorModel) error {
	params := mappers.InsertSectorParamsFromDomain(sector)
	_, err := r.q.CreateSector(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (r *SectorRepo) Update(ctx context.Context, sector *domain.SectorModel) error {
	params := mappers.UpdateSectorParamsFromDomain(sector)
	_, err := r.q.UpdateSector(ctx, params)
	return err
}

func (r *SectorRepo) FindByCoordinates(ctx context.Context, x int, y int) (*domain.SectorModel, error) {
	row, err := r.q.GetSectorByCoordinates(ctx, gen.GetSectorByCoordinatesParams{X: int32(x), Y: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.SectorFromDB(row), nil
}

func (r *SectorRepo) FindAll(ctx context.Context) ([]*domain.SectorModel, error) {
	rows, err := r.q.ListSectors(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*domain.SectorModel, 0, len(rows))
	for _, row := range rows {
		out = append(out, mappers.SectorFromDB(row))
	}
	return out, nil
}

func (r *SectorRepo) FindByCoordinatesForUpdate(ctx context.Context, x int, y int) (*domain.SectorModel, error) {
	row, err := r.q.GetSectorByCoordinatesForUpdate(ctx, gen.GetSectorByCoordinatesForUpdateParams{X: int32(x), Y: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.SectorFromDB(row), nil
}

func (r *SectorRepo) ListOccupiedCoordinates(ctx context.Context) ([]domain.Vector2i, error) {
	rows, err := r.q.ListOccupiedSectorCoordinates(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Vector2i, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Vector2i{X: int(row.X), Y: int(row.Y)})
	}
	return out, nil
}

func (r *SectorRepo) GetLocationTypeByCoordinates(ctx context.Context, x int, y int) (domain.LocationType, error) {
	row, err := r.q.GetLocationTypeByCoordinates(ctx, gen.GetLocationTypeByCoordinatesParams{X: int32(x), Y: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.LocationTypeEmpty, ports.ErrNotFound
		}
		return domain.LocationTypeEmpty, err
	}
	return domain.LocationType(row), nil
}

func (r *SectorRepo) CountLocationsInRange(ctx context.Context, x, y, radius int) (resourceful int, dangerous int, err error) {
	params := gen.CountResourcefulLocationsInRangeParams{
		CenterX: int32(x),
		CenterY: int32(y),
		Radius:  int32(radius),
	}
	rCount, err := r.q.CountResourcefulLocationsInRange(ctx, params)
	if err != nil {
		return 0, 0, err
	}
	dCount, err := r.q.CountDangerousLocationsInRange(ctx, gen.CountDangerousLocationsInRangeParams(params))
	if err != nil {
		return 0, 0, err
	}
	return int(rCount), int(dCount), nil
}
