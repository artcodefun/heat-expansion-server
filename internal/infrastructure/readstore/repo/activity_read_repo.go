package repo

import (
	"context"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/mappers"
)

type ActivityReadRepo struct{ q *gen.Queries }

func NewActivityReadRepo(q *gen.Queries) *ActivityReadRepo { return &ActivityReadRepo{q: q} }

func (r *ActivityReadRepo) ListActivities(baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListActivities(context.Background(), gen.ListActivitiesParams{BaseID: int64(baseID), Limit: int32(limit)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ActivityItem, 0, len(rows))
	for _, a := range rows {
		v := mappers.ActivityItemFromModel(a)
		out = append(out, &v)
	}
	return out, nil
}

func (r *ActivityReadRepo) ListMilitaryActivities(baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	rows, err := r.q.ListMilitaryActivities(context.Background(), gen.ListMilitaryActivitiesParams{BaseID: int64(baseID), Limit: int32(limit)})
	if err != nil {
		return nil, err
	}
	out := make([]*readmodels.ActivityItem, 0, len(rows))
	for _, a := range rows {
		v := mappers.ActivityItemFromModel(a)
		out = append(out, &v)
	}
	return out, nil
}
