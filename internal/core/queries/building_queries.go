package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type BuildingQueries struct {
	Repo      ports.BuildingReadRepository
	ProtoRepo ports.BuildPrototypeRepository
	BaseRepo  ports.UserBaseRepository
	Access    *services.AccessControlService
}

func NewBuildingQueries(repo ports.BuildingReadRepository, protoRepo ports.BuildPrototypeRepository, baseRepo ports.UserBaseRepository, access *services.AccessControlService) *BuildingQueries {
	return &BuildingQueries{Repo: repo, ProtoRepo: protoRepo, BaseRepo: baseRepo, Access: access}
}

func (q *BuildingQueries) ListNewBuildItems(ctx cqrs.QueryContext, baseID int, category string) ([]*readmodels.BuildItemNew, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	// Orchestrate: load all prototypes and base, compute availability via domain
	allProtos, err := q.ProtoRepo.FindAllPrototypes()
	if err != nil {
		return nil, err
	}
	base, err := q.BaseRepo.FindByID(baseID)
	if err != nil {
		return nil, err
	}
	available := base.AvailableBuildings(allProtos)
	ids := make([]int, 0, len(available))
	for _, p := range available {
		if category == "" || string(p.Category) == category {
			ids = append(ids, p.ID)
		}
	}
	return q.Repo.ListNewBuildItemsByPrototypeIDs(ids)
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
