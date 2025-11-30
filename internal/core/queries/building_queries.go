package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type BuildingQueries struct {
	Repo   ports.BuildingReadRepository
	Access *services.AccessControlService
}

func NewBuildingQueries(repo ports.BuildingReadRepository, access *services.AccessControlService) *BuildingQueries {
	return &BuildingQueries{Repo: repo, Access: access}
}

func (q *BuildingQueries) ListNewBuildItems(ctx cqrs.QueryContext, baseID int, category string) ([]*readmodels.BuildItemNew, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListNewBuildItems(baseID, category)
}
func (q *BuildingQueries) ListPendingBuildItems(ctx cqrs.QueryContext, baseID int, category string) ([]*readmodels.BuildItemPending, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListPendingBuildItems(baseID, category)
}
func (q *BuildingQueries) ListInProductionBuildItems(ctx cqrs.QueryContext, baseID int, category string) ([]*readmodels.BuildItemInProduction, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListInProductionBuildItems(baseID, category)
}
func (q *BuildingQueries) ListPresentBuildItems(ctx cqrs.QueryContext, baseID int, category string) ([]*readmodels.BuildItemPresent, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListPresentBuildItems(baseID, category)
}
