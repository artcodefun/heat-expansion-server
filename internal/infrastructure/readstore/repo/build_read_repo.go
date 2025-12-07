package repo

import (
	"context"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/mappers"
)

type BuildReadRepo struct{ q *gen.Queries }

func NewBuildReadRepo(q *gen.Queries) *BuildReadRepo { return &BuildReadRepo{q: q} }

func (r *BuildReadRepo) ListNewBuildItemsByPrototypeIDs(ids []int) ([]*readmodels.BuildItemNew, error) {
	if len(ids) == 0 {
		return []*readmodels.BuildItemNew{}, nil
	}
	rows, err := r.q.ListBuildPrototypesByIDs(context.Background(), mappers.IdsToInt64(ids))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.BuildItemNew, 0, len(rows))
	for _, p := range rows {
		v := mappers.NewBuildItemFromPrototype(p)
		out = append(out, &v)
	}
	return out, nil
}

func (r *BuildReadRepo) ListPendingBuildItems(baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPending, error) {
	rs, err := r.q.ListPendingBuildItems(context.Background(), gen.ListPendingBuildItemsParams{BaseID: int64(baseID), Category: string(category)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.BuildItemPending, 0, len(rs))
	for _, r0 := range rs {
		v := mappers.BuildItemPendingFromRow(r0)
		out = append(out, &v)
	}
	return out, nil
}

func (r *BuildReadRepo) ListInProductionBuildItems(baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemInProduction, error) {
	rows, err := r.q.ListInProductionBuildItems(context.Background(), gen.ListInProductionBuildItemsParams{BaseID: int64(baseID), Category: string(category)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.BuildItemInProduction, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.BuildItemInProductionFromRow(r0)
		out = append(out, &v)
	}
	return out, nil
}

func (r *BuildReadRepo) ListPresentBuildItems(baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPresent, error) {
	rows, err := r.q.ListPresentBuildItems(context.Background(), gen.ListPresentBuildItemsParams{BaseID: int64(baseID), Category: string(category)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.BuildItemPresent, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.BuildItemPresentFromRow(r0)
		out = append(out, &v)
	}
	return out, nil
}
