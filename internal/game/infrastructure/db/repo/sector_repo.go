package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/mappers"
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
func (r *SectorRepo) Create(sector *domain.SectorModel) error {
	params := mappers.InsertSectorParamsFromDomain(sector)
	_, err := r.q.CreateSector(context.Background(), params)
	if err != nil {
		return err
	}
	return nil
}

func (r *SectorRepo) Update(sector *domain.SectorModel) error {
	params := mappers.UpdateSectorParamsFromDomain(sector)
	_, err := r.q.UpdateSector(context.Background(), params)
	return err
}

func (r *SectorRepo) FindByCoordinates(x int, y int) (*domain.SectorModel, error) {
	row, err := r.q.GetSectorByCoordinates(context.Background(), gen.GetSectorByCoordinatesParams{X: int32(x), Y: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.SectorFromDB(row), nil
}

func (r *SectorRepo) FindAll() ([]*domain.SectorModel, error) {
	rows, err := r.q.ListSectors(context.Background())
	if err != nil {
		return nil, err
	}
	out := make([]*domain.SectorModel, 0, len(rows))
	for _, row := range rows {
		out = append(out, mappers.SectorFromDB(row))
	}
	return out, nil
}

func (r *SectorRepo) FindByCoordinatesForUpdate(x int, y int) (*domain.SectorModel, error) {
	row, err := r.q.GetSectorByCoordinatesForUpdate(context.Background(), gen.GetSectorByCoordinatesForUpdateParams{X: int32(x), Y: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.SectorFromDB(row), nil
}

func (r *SectorRepo) ListOccupiedCoordinates() ([]domain.Vector2i, error) {
	rows, err := r.q.ListOccupiedSectorCoordinates(context.Background())
	if err != nil {
		return nil, err
	}
	out := make([]domain.Vector2i, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Vector2i{X: int(row.X), Y: int(row.Y)})
	}
	return out, nil
}

func (r *SectorRepo) GetLocationTypeByCoordinates(x int, y int) (domain.LocationType, error) {
	row, err := r.q.GetLocationTypeByCoordinates(context.Background(), gen.GetLocationTypeByCoordinatesParams{X: int32(x), Y: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.LocationTypeEmpty, ports.ErrNotFound
		}
		return domain.LocationTypeEmpty, err
	}
	return domain.LocationType(row), nil
}

func (r *SectorRepo) CountLocationsInRange(x, y, radius int) (resourceful int, dangerous int, err error) {
	ctx := context.Background()
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
