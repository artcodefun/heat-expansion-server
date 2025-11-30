package repo

import (
	"context"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/mappers"
)

type StorageReadRepo struct{ q *gen.Queries }

func NewStorageReadRepo(q *gen.Queries) *StorageReadRepo { return &StorageReadRepo{q: q} }

func (r *StorageReadRepo) ListPresentStorageItems(baseID int) ([]*readmodels.StorageItemPresent, error) {
	rows, err := r.q.ListPresentStorageItems(context.Background(), int64(baseID))
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
