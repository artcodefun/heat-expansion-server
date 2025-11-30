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

type ResourceLocationRepo struct {
	q *gen.Queries
}

func NewResourceLocationRepo(q *gen.Queries) *ResourceLocationRepo {
	return &ResourceLocationRepo{q: q}
}

func (r *ResourceLocationRepo) Tx(tx ports.Transaction) ports.ResourceLocationRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &ResourceLocationRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *ResourceLocationRepo) Create(loc *domain.ResourceLocationModel) error {
	id, err := r.q.InsertResourceLocation(context.Background(), mappers.InsertResourceLocationParamsFromDomain(loc))
	if err != nil {
		return err
	}
	loc.ID = int(id)
	return nil
}

func (r *ResourceLocationRepo) FindByID(id int) (*domain.ResourceLocationModel, error) {
	row, err := r.q.GetResourceLocationByID(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	loc := mappers.ResourceLocationFromDB(row)
	return loc, nil
}

func (r *ResourceLocationRepo) FindByCoordinates(x, y int) (*domain.ResourceLocationModel, error) {
	row, err := r.q.GetResourceLocationBySector(context.Background(), gen.GetResourceLocationBySectorParams{SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	loc := mappers.ResourceLocationFromDB(row)
	return loc, nil
}

func (r *ResourceLocationRepo) FindByCoordinatesForUpdate(x, y int) (*domain.ResourceLocationModel, error) {
	row, err := r.q.GetResourceLocationBySectorForUpdate(context.Background(), gen.GetResourceLocationBySectorForUpdateParams{SectorX: int32(x), SectorY: int32(y)})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	loc := mappers.ResourceLocationFromDB(row)
	return loc, nil
}

func (r *ResourceLocationRepo) Update(loc *domain.ResourceLocationModel) error {
	return r.q.UpdateResourceLocation(context.Background(), mappers.UpdateResourceLocationParamsFromDomain(loc))
}

func (r *ResourceLocationRepo) Delete(id int) error {
	return r.q.DeleteResourceLocation(context.Background(), int64(id))
}
