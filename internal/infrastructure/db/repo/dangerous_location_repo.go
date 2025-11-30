package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/mappers"
)

type DangerousLocationRepo struct {
	q *gen.Queries
}

func NewDangerousLocationRepo(q *gen.Queries) *DangerousLocationRepo {
	return &DangerousLocationRepo{q: q}
}

func (r *DangerousLocationRepo) Tx(tx ports.Transaction) ports.DangerousLocationRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &DangerousLocationRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *DangerousLocationRepo) Create(loc *domain.DangerousLocationModel) error {
	id, err := r.q.InsertDangerousLocation(context.Background(), mappers.InsertDangerousLocationParamsFromDomain(loc))
	if err != nil {
		return err
	}
	loc.ID = int(id)
	return nil
}

func (r *DangerousLocationRepo) FindByID(id int) (*domain.DangerousLocationModel, error) {
	row, err := r.q.GetDangerousLocationByID(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	loc := mappers.DangerousLocationFromDB(row)
	return loc, nil
}

func (r *DangerousLocationRepo) FindByCoordinates(x, y int) (*domain.DangerousLocationModel, error) {
	row, err := r.q.GetDangerousLocationBySector(context.Background(), gen.GetDangerousLocationBySectorParams{SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	loc := mappers.DangerousLocationFromDB(row)
	return loc, nil
}

func (r *DangerousLocationRepo) FindByCoordinatesForUpdate(x, y int) (*domain.DangerousLocationModel, error) {
	row, err := r.q.GetDangerousLocationBySectorForUpdate(context.Background(), gen.GetDangerousLocationBySectorForUpdateParams{SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	loc := mappers.DangerousLocationFromDB(row)
	return loc, nil
}

func (r *DangerousLocationRepo) Update(loc *domain.DangerousLocationModel) error {
	return r.q.UpdateDangerousLocation(context.Background(), mappers.UpdateDangerousLocationParamsFromDomain(loc))
}

func (r *DangerousLocationRepo) Delete(id int) error {
	return r.q.DeleteDangerousLocation(context.Background(), int64(id))
}
