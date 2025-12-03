package repo

import (
	"context"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/mappers"
)

type TechReadRepo struct{ q *gen.Queries }

func NewTechReadRepo(q *gen.Queries) *TechReadRepo { return &TechReadRepo{q: q} }

func (r *TechReadRepo) ListNewTechItemsByPrototypeIDs(ids []int) ([]*readmodels.TechItemNew, error) {
	if len(ids) == 0 {
		return []*readmodels.TechItemNew{}, nil
	}
	rows, err := r.q.ListTechPrototypesByIDs(context.Background(), mappers.IdsToInt64(ids))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.TechItemNew, 0, len(rows))
	for _, p := range rows {
		v := mappers.NewTechItemFromPrototype(p)
		out = append(out, &v)
	}
	return out, nil
}

func (r *TechReadRepo) ListInResearchTechItems(baseID int) ([]*readmodels.TechItemInProgress, error) {
	rows, err := r.q.ListInResearchTechItems(context.Background(), int64(baseID))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.TechItemInProgress, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.TechItemInProgressFromRow(r0)
		out = append(out, &v)
	}
	return out, nil
}

func (r *TechReadRepo) ListDoneTechItems(baseID int) ([]*readmodels.TechItemDone, error) {
	rows, err := r.q.ListDoneTechItems(context.Background(), int64(baseID))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.TechItemDone, 0, len(rows))
	for _, r0 := range rows {
		v := mappers.TechItemDoneFromRow(r0)
		out = append(out, &v)
	}
	return out, nil
}
