package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type ActivityQueries struct {
	Repo   ports.ActivityReadRepository
	Access *services.AccessControlService
}

func NewActivityQueries(repo ports.ActivityReadRepository, access *services.AccessControlService) *ActivityQueries {
	return &ActivityQueries{Repo: repo, Access: access}
}

func (q *ActivityQueries) ListActivities(ctx cqrs.QueryContext, baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListActivities(baseID, limit)
}
func (q *ActivityQueries) ListMilitaryActivities(ctx cqrs.QueryContext, baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListMilitaryActivities(baseID, limit)
}
