package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
)

type BaseQueries struct {
	Repo   ports.BaseReadRepository
	Access *services.AccessControlService
}

func NewBaseQueries(repo ports.BaseReadRepository, access *services.AccessControlService) *BaseQueries {
	return &BaseQueries{Repo: repo, Access: access}
}

func (q *BaseQueries) GetBaseStats(ctx context.Context, actor cqrs.Actor, baseID int) (*readmodels.UserBaseStats, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	stats, err := q.Repo.GetBaseStats(ctx, baseID)
	return stats, repoErr(err)
}

// ListUserBases returns basic info for bases owned by the authenticated user.
func (q *BaseQueries) ListUserBases(ctx context.Context, actor cqrs.Actor) ([]*readmodels.UserBaseModel, error) {
	// Only allow requesting own bases for now; later add roles/tenant etc.
	bases, err := q.Repo.ListUserBases(ctx, actor.UserID)
	return bases, repoErr(err)
}
