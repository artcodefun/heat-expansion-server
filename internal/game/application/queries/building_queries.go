package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
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

func (q *BuildingQueries) ListNewBuildItems(ctx context.Context, actor cqrs.Actor, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemNew, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	// Orchestrate: load all prototypes and base, compute availability via domain
	allProtos, err := q.ProtoRepo.FindAllPrototypes(ctx)
	if err != nil {
		return nil, repoErr(err)
	}
	base, err := q.BaseRepo.FindByID(ctx, baseID)
	if err != nil {
		return nil, repoErr(err)
	}
	available := base.AvailableBuildings(allProtos)
	ids := make([]int, 0, len(available))
	for _, p := range available {
		if category == "" || string(p.Category) == string(category) {
			ids = append(ids, p.ID)
		}
	}
	items, err := q.Repo.ListNewBuildItemsByPrototypeIDs(ctx, ids)
	return items, repoErr(err)
}
func (q *BuildingQueries) ListPendingBuildItems(ctx context.Context, actor cqrs.Actor, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPending, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListPendingBuildItems(ctx, baseID, category)
	return items, repoErr(err)
}
func (q *BuildingQueries) ListInProductionBuildItems(ctx context.Context, actor cqrs.Actor, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemInProduction, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListInProductionBuildItems(ctx, baseID, category)
	return items, repoErr(err)
}
func (q *BuildingQueries) ListPresentBuildItems(ctx context.Context, actor cqrs.Actor, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPresent, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListPresentBuildItems(ctx, baseID, category)
	return items, repoErr(err)
}
