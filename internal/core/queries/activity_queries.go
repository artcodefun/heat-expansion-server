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

func (q *ActivityQueries) ListOffenseActivities(ctx cqrs.QueryContext, baseID int, subtype readmodels.OffenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListOffenseActivities(baseID, subtype, limit)
	return items, repoErr(err)
}

func (q *ActivityQueries) ListDefenseActivities(ctx cqrs.QueryContext, baseID int, subtype readmodels.DefenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListDefenseActivities(baseID, subtype, limit)
	return items, repoErr(err)
}

func (q *ActivityQueries) ListScanActivities(ctx cqrs.QueryContext, baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListScanActivities(baseID, limit)
	return items, repoErr(err)
}

func (q *ActivityQueries) ListRadarActivities(ctx cqrs.QueryContext, baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListRadarActivities(baseID, limit)
	return items, repoErr(err)
}

func (q *ActivityQueries) ListTradeActivities(ctx cqrs.QueryContext, baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListTradeActivities(baseID, limit)
	return items, repoErr(err)
}
