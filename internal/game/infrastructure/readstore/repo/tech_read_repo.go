package repo

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/mappers"
)

type TechReadRepo struct{ q *gen.Queries }

func NewTechReadRepo(q *gen.Queries) *TechReadRepo { return &TechReadRepo{q: q} }

func (r *TechReadRepo) ListNewTechItemsByPrototypeIDs(ctx context.Context, baseID int, ids []int) ([]*readmodels.TechItemNew, error) {
	if len(ids) == 0 {
		return []*readmodels.TechItemNew{}, nil
	}
	// Get finished items to find current levels
	doneRows, err := r.q.ListDoneTechItemsAll(ctx, int64(baseID))
	if err != nil {
		return nil, err
	}
	levels := make(map[int]int)
	for _, dr := range doneRows {
		v := mappers.TechItemDoneFromAllRow(dr)
		levels[v.Prototype.ID] = v.Level
	}

	rows, err := r.q.ListTechPrototypesByIDs(ctx, mappers.IdsToInt64(ids))
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.TechItemNew, 0, len(rows))
	for _, p := range rows {
		level := levels[int(p.ID)]
		v := mappers.NewTechItemFromPrototype(p, level)
		out = append(out, &v)
	}
	return out, nil
}

func (r *TechReadRepo) ListInResearchTechItems(ctx context.Context, baseID int, category readmodels.TechCategory) ([]*readmodels.TechItemInProgress, error) {
	// Get finished items to find current levels
	doneRows, err := r.q.ListDoneTechItemsAll(ctx, int64(baseID))
	if err != nil {
		return nil, err
	}
	levels := make(map[int]int)
	for _, dr := range doneRows {
		v := mappers.TechItemDoneFromAllRow(dr)
		levels[v.Prototype.ID] = v.Level
	}

	rows, err := r.q.ListInResearchTechItems(ctx, gen.ListInResearchTechItemsParams{
		BaseID:   int64(baseID),
		Category: string(category),
	})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.TechItemInProgress, 0, len(rows))
	for _, r0 := range rows {
		level := levels[int(r0.PrototypeID)]
		v := mappers.TechItemInProgressFromRow(r0, level)
		out = append(out, &v)
	}
	return out, nil
}

func (r *TechReadRepo) ListDoneTechItems(ctx context.Context, baseID int, category readmodels.TechCategory) ([]*readmodels.TechItemDone, error) {
	rows, err := r.q.ListDoneTechItems(ctx, gen.ListDoneTechItemsParams{
		BaseID:   int64(baseID),
		Category: string(category),
	})
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
