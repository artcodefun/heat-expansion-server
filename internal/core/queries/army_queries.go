package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type ArmyQueries struct {
	Repo   ports.ArmyReadRepository
	Access *services.AccessControlService
}

func NewArmyQueries(repo ports.ArmyReadRepository, access *services.AccessControlService) *ArmyQueries {
	return &ArmyQueries{Repo: repo, Access: access}
}

func (q *ArmyQueries) ListNewArmyItems(ctx cqrs.QueryContext, baseID int, category string) ([]*readmodels.ArmyItemNew, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListNewArmyItems(baseID, category)
}
func (q *ArmyQueries) ListPendingArmyItems(ctx cqrs.QueryContext, baseID int, category string) ([]*readmodels.ArmyItemPending, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListPendingArmyItems(baseID, category)
}
func (q *ArmyQueries) ListInProductionArmyItems(ctx cqrs.QueryContext, baseID int, category string) ([]*readmodels.ArmyItemInProduction, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListInProductionArmyItems(baseID, category)
}
func (q *ArmyQueries) ListPresentArmyItems(ctx cqrs.QueryContext, baseID int, category string) ([]*readmodels.ArmyItemPresent, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListPresentArmyItems(baseID, category)
}
