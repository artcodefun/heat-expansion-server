package repo

import (
	"context"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/mappers"
)

type ArmyReadRepo struct{ q *gen.Queries }

func NewArmyReadRepo(q *gen.Queries) *ArmyReadRepo { return &ArmyReadRepo{q: q} }

func (r *ArmyReadRepo) ListNewArmyItemsByPrototypeIDs(ids []int) ([]*readmodels.ArmyItemNew, error) {
	if len(ids) == 0 {
		return []*readmodels.ArmyItemNew{}, nil
	}
	rows, err := r.q.ListArmyPrototypesByIDs(context.Background(), mappers.IdsToInt64(ids))
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

func (r *ArmyReadRepo) ListPendingArmyItems(baseID int, category string) ([]*readmodels.ArmyItemPending, error) {
	rs, err := r.q.ListPendingArmyItems(context.Background(), gen.ListPendingArmyItemsParams{BaseID: int64(baseID), Category: category})
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

func (r *ArmyReadRepo) ListInProductionArmyItems(baseID int, category string) ([]*readmodels.ArmyItemInProduction, error) {
	rows, err := r.q.ListInProductionArmyItems(context.Background(), gen.ListInProductionArmyItemsParams{BaseID: int64(baseID), Category: category})
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

func (r *ArmyReadRepo) ListPresentArmyItems(baseID int, category string) ([]*readmodels.ArmyItemPresent, error) {
	rows, err := r.q.ListPresentArmyItems(context.Background(), gen.ListPresentArmyItemsParams{BaseID: int64(baseID), Category: category})
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
