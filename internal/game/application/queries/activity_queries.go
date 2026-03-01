package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
)

type ActivityQueries struct {
	Repo   ports.ActivityReadRepository
	Access *services.AccessControlService
}

func NewActivityQueries(repo ports.ActivityReadRepository, access *services.AccessControlService) *ActivityQueries {
	return &ActivityQueries{Repo: repo, Access: access}
}

func (q *ActivityQueries) ListOffenseActivities(ctx context.Context, actor cqrs.Actor, baseID int, subtype readmodels.OffenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListOffenseActivities(ctx, baseID, subtype, limit)
	return items, repoErr(err)
}

func (q *ActivityQueries) ListDefenseActivities(ctx context.Context, actor cqrs.Actor, baseID int, subtype readmodels.DefenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListDefenseActivities(ctx, baseID, subtype, limit)
	return items, repoErr(err)
}

func (q *ActivityQueries) ListScanActivities(ctx context.Context, actor cqrs.Actor, baseID int, subtype readmodels.ScanActivitySubtype, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListScanActivities(ctx, baseID, subtype, limit)
	return items, repoErr(err)
}

func (q *ActivityQueries) ListRadarActivities(ctx context.Context, actor cqrs.Actor, baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListRadarActivities(ctx, baseID, limit)
	return items, repoErr(err)
}

func (q *ActivityQueries) ListTradeActivities(ctx context.Context, actor cqrs.Actor, baseID int, limit int) ([]*readmodels.ActivityItem, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListTradeActivities(ctx, baseID, limit)
	return items, repoErr(err)
}
