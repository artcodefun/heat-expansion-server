package repo

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
)

type ArmyReadRepo struct{ q *gen.Queries }

func NewArmyReadRepo(q *gen.Queries) *ArmyReadRepo { return &ArmyReadRepo{q: q} }

func (r *ArmyReadRepo) ListNewArmyItemsByPrototypeIDs(ctx context.Context, ids []int) ([]*readmodels.ArmyItemNew, error) {
	if len(ids) == 0 {
		return []*readmodels.ArmyItemNew{}, nil
	}
	rows, err := r.q.ListArmyPrototypesByIDs(ctx, mappers.IdsToInt64(ids))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ArmyItemNew, 0, len(rows))
	for _, p := range rows {
		v := mappers.NewArmyItemFromPrototype(p)
		out = append(out, &v)
	}
	return out, nil
}

func (r *ArmyReadRepo) ListPendingArmyItems(ctx context.Context, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPending, error) {
	rs, err := r.q.ListPendingArmyItems(ctx, gen.ListPendingArmyItemsParams{BaseID: int64(baseID), Category: string(category)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ArmyItemPending, 0, len(rs))
	for _, r0 := range rs {
		v := mappers.ArmyItemPendingFromRow(r0)
		out = append(out, &v)
	}
	return out, nil
}

func (r *ArmyReadRepo) ListInProductionArmyItems(ctx context.Context, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemInProduction, error) {
	rows, err := r.q.ListInProductionArmyItems(ctx, gen.ListInProductionArmyItemsParams{BaseID: int64(baseID), Category: string(category)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ArmyItemInProduction, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.ArmyItemInProductionFromRow(r0)
		out = append(out, &v)
	}
	return out, nil
}

func (r *ArmyReadRepo) ListPresentArmyItems(ctx context.Context, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPresent, error) {
	rows, err := r.q.ListPresentArmyItems(ctx, gen.ListPresentArmyItemsParams{BaseID: int64(baseID), Category: string(category)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ArmyItemPresent, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.ArmyItemPresentFromRow(r0)
		out = append(out, &v)
	}
	return out, nil
}

func (r *ArmyReadRepo) ListPresentArmyItemsAll(ctx context.Context, baseID int) ([]*readmodels.ArmyItemPresent, error) {
	rows, err := r.q.ListPresentArmyItemsAll(ctx, int64(baseID))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ArmyItemPresent, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.ArmyItemPresentFromAllRow(r0)
		out = append(out, &v)
	}
	return out, nil
}
