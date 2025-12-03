package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type BaseQueries struct {
	Repo   ports.BaseReadRepository
	Access *services.AccessControlService
}

func NewBaseQueries(repo ports.BaseReadRepository, access *services.AccessControlService) *BaseQueries {
	return &BaseQueries{Repo: repo, Access: access}
}

func (q *BaseQueries) GetBaseStats(ctx cqrs.QueryContext, baseID int) (*readmodels.UserBaseStats, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	stats, err := q.Repo.GetBaseStats(baseID)
	return stats, repoErr(err)
}

// ListUserBases returns basic info for bases owned by the authenticated user.
func (q *BaseQueries) ListUserBases(ctx cqrs.QueryContext) ([]*readmodels.UserBaseModel, error) {
	// Only allow requesting own bases for now; later add roles/tenant etc.
	bases, err := q.Repo.ListUserBases(ctx.UserID)
	return bases, repoErr(err)
}
