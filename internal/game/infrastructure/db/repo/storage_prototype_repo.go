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

func (r *StoragePrototypeRepo) CreatePrototype(ctx context.Context, proto *domain.StorageItemPrototype) error {
	return r.q.CreateStoragePrototype(ctx, mappers.StoragePrototypeToCreateParams(proto))
}

func (r *StoragePrototypeRepo) FindPrototypeByID(ctx context.Context, id int) (*domain.StorageItemPrototype, error) {
	p, err := r.q.GetStoragePrototypeByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.StoragePrototypeFromDB(p), nil
}

func (r *StoragePrototypeRepo) FindAllPrototypes(ctx context.Context) ([]*domain.StorageItemPrototype, error) {
	rows, err := r.q.ListStoragePrototypes(ctx)
	if err != nil {
		return nil, err
	}
	return mappers.StoragePrototypesFromDB(rows), nil
}

func (r *StoragePrototypeRepo) UpdatePrototype(ctx context.Context, proto *domain.StorageItemPrototype) error {
	_, err := r.q.UpdateStoragePrototype(ctx, mappers.StoragePrototypeToUpdateParams(proto))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ports.ErrNotFound
		}
		return err
	}
	return nil
}

func (r *StoragePrototypeRepo) DeletePrototype(_ context.Context, id int) error {
	return errors.New("DeletePrototype not implemented for storage prototypes (read-only in this service)")
}
