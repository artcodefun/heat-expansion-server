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
	items, err := q.Repo.ListActivities(baseID, limit)
	return items, repoErr(err)
}

func (q *ActivityQueries) ListActivitiesByKind(ctx cqrs.QueryContext, baseID int, kind readmodels.ActivityKind, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListActivitiesByKind(baseID, kind, limit)
	return items, repoErr(err)
}
