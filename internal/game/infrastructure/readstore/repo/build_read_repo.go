package repo

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
)

type BuildReadRepo struct{ q *gen.Queries }

func NewBuildReadRepo(q *gen.Queries) *BuildReadRepo { return &BuildReadRepo{q: q} }

func (r *BuildReadRepo) ListNewBuildItemsByPrototypeIDs(ctx context.Context, ids []int) ([]*readmodels.BuildItemNew, error) {
	if len(ids) == 0 {
		return []*readmodels.BuildItemNew{}, nil
	}
	rows, err := r.q.ListBuildPrototypesByIDs(ctx, mappers.IdsToInt64(ids))
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

func (r *BuildReadRepo) ListPendingBuildItems(ctx context.Context, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPending, error) {
	rs, err := r.q.ListPendingBuildItems(ctx, gen.ListPendingBuildItemsParams{BaseID: int64(baseID), Category: string(category)})
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

func (r *BuildReadRepo) ListInProductionBuildItems(ctx context.Context, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemInProduction, error) {
	rows, err := r.q.ListInProductionBuildItems(ctx, gen.ListInProductionBuildItemsParams{BaseID: int64(baseID), Category: string(category)})
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

func (r *BuildReadRepo) ListPresentBuildItems(ctx context.Context, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPresent, error) {
	rows, err := r.q.ListPresentBuildItems(ctx, gen.ListPresentBuildItemsParams{BaseID: int64(baseID), Category: string(category)})
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
