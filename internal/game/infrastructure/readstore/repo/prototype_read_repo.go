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

type PrototypeReadRepo struct{ q *gen.Queries }

func NewPrototypeReadRepo(q *gen.Queries) *PrototypeReadRepo { return &PrototypeReadRepo{q: q} }

func (r *PrototypeReadRepo) ListArmyPrototypes(ctx context.Context) ([]*readmodels.ArmyItemPrototype, error) {
	rows, err := r.q.ListArmyPrototypes(ctx)
	if err != nil {
		return nil, err
	}
	return mappers.ArmyPrototypesFromModels(rows), nil
}

func (r *PrototypeReadRepo) GetArmyPrototype(ctx context.Context, id int) (*readmodels.ArmyItemPrototype, error) {
	row, err := r.q.GetArmyPrototypeByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	v := mappers.ArmyPrototypeFromModel(row)
	return &v, nil
}

func (r *PrototypeReadRepo) ListBuildPrototypes(ctx context.Context) ([]*readmodels.BuildItemPrototype, error) {
	rows, err := r.q.ListBuildPrototypes(ctx)
	if err != nil {
		return nil, err
	}
	return mappers.BuildPrototypesFromModels(rows), nil
}

func (r *PrototypeReadRepo) GetBuildPrototype(ctx context.Context, id int) (*readmodels.BuildItemPrototype, error) {
	row, err := r.q.GetBuildPrototypeByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	v := mappers.BuildPrototypeFromModel(row)
	return &v, nil
}

func (r *PrototypeReadRepo) ListStoragePrototypes(ctx context.Context) ([]*readmodels.StorageItemPrototype, error) {
	rows, err := r.q.ListStoragePrototypes(ctx)
	if err != nil {
		return nil, err
	}
	return mappers.StoragePrototypesFromModels(rows), nil
}

func (r *PrototypeReadRepo) GetStoragePrototype(ctx context.Context, id int) (*readmodels.StorageItemPrototype, error) {
	row, err := r.q.GetStoragePrototypeByID(ctx, int64(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	v := mappers.StoragePrototypeFromModel(row)
	return &v, nil
}
