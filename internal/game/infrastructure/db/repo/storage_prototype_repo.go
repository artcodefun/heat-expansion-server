package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/mappers"
)

type StoragePrototypeRepo struct {
	q *gen.Queries
}

func NewStoragePrototypeRepo(q *gen.Queries) *StoragePrototypeRepo {
	return &StoragePrototypeRepo{q: q}
}

func (r *StoragePrototypeRepo) Tx(tx ports.Transaction) ports.StoragePrototypeRepository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &StoragePrototypeRepo{q: r.q.WithTx(sqlTx)}
	}
	return r
}

func (r *StoragePrototypeRepo) CreatePrototype(proto *domain.StorageItemPrototype) error {
	return errors.New("CreatePrototype not implemented for storage prototypes (read-only in this service)")
}

func (r *StoragePrototypeRepo) FindPrototypeByID(id int) (*domain.StorageItemPrototype, error) {
	p, err := r.q.GetStoragePrototypeByID(context.Background(), int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.StoragePrototypeFromDB(p), nil
}

func (r *StoragePrototypeRepo) FindAllPrototypes() ([]*domain.StorageItemPrototype, error) {
	rows, err := r.q.ListStoragePrototypes(context.Background())
	if err != nil {
		return nil, err
	}
	return mappers.StoragePrototypesFromDB(rows), nil
}

func (r *StoragePrototypeRepo) UpdatePrototype(proto *domain.StorageItemPrototype) error {
	return errors.New("UpdatePrototype not implemented for storage prototypes (read-only in this service)")
}

func (r *StoragePrototypeRepo) DeletePrototype(id int) error {
	return errors.New("DeletePrototype not implemented for storage prototypes (read-only in this service)")
}
