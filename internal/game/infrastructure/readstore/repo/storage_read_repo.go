package repo

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
)

type StorageReadRepo struct{ q *gen.Queries }

func NewStorageReadRepo(q *gen.Queries) *StorageReadRepo { return &StorageReadRepo{q: q} }

func (r *StorageReadRepo) ListPresentStorageItems(ctx context.Context, baseID int, category readmodels.StorageCategory) ([]*readmodels.StorageItemPresent, error) {
	rows, err := r.q.ListPresentStorageItems(ctx, gen.ListPresentStorageItemsParams{
		BaseID:   int64(baseID),
		Category: string(category),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.StorageItemPresent, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.StorageItemPresentFromRow(r0)
		out = append(out, &v)
	}
	return out, nil
}

func (r *StorageReadRepo) ListTradeableStorageItems(ctx context.Context, baseID int) ([]*readmodels.StorageItemPresent, error) {
	rows, err := r.q.ListTradeableStorageItems(ctx, int64(baseID))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.StorageItemPresent, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.StorageItemPresentFromTradeableRow(r0)
		out = append(out, &v)
	}
	return out, nil
}
