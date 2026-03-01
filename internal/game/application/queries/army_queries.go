package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
)

type ArmyQueries struct {
	Repo      ports.ArmyReadRepository
	ProtoRepo ports.ArmyPrototypeRepository
	BaseRepo  ports.UserBaseRepository
	Access    *services.AccessControlService
}

func NewArmyQueries(repo ports.ArmyReadRepository, protoRepo ports.ArmyPrototypeRepository, baseRepo ports.UserBaseRepository, access *services.AccessControlService) *ArmyQueries {
	return &ArmyQueries{Repo: repo, ProtoRepo: protoRepo, BaseRepo: baseRepo, Access: access}
}

func (q *ArmyQueries) ListNewArmyItems(ctx context.Context, actor cqrs.Actor, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemNew, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	// Load all army prototypes and user base aggregate
	allProtos, err := q.ProtoRepo.FindAllPrototypes(ctx)
	if err != nil {
		return nil, repoErr(err)
	}
	base, err := q.BaseRepo.FindByID(ctx, baseID)
	if err != nil {
		return nil, repoErr(err)
	}

	// Compute available prototypes using domain logic
	available := base.AvailableArmies(allProtos)
	// Filter by category if provided
	ids := make([]int, 0, len(available))
	for _, p := range available {
		if category == "" || string(p.Category) == string(category) {
			ids = append(ids, p.ID)
		}
	}
	items, err := q.Repo.ListNewArmyItemsByPrototypeIDs(ctx, ids)
	return items, repoErr(err)
}
func (q *ArmyQueries) ListPendingArmyItems(ctx context.Context, actor cqrs.Actor, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPending, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListPendingArmyItems(ctx, baseID, category)
	return items, repoErr(err)
}
func (q *ArmyQueries) ListInProductionArmyItems(ctx context.Context, actor cqrs.Actor, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemInProduction, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListInProductionArmyItems(ctx, baseID, category)
	return items, repoErr(err)
}
func (q *ArmyQueries) ListPresentArmyItems(ctx context.Context, actor cqrs.Actor, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPresent, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.Repo.ListPresentArmyItems(ctx, baseID, category)
	return items, repoErr(err)
}
