package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type TechQueries struct {
	Repo   ports.TechReadRepository
	Access *services.AccessControlService
}

func NewTechQueries(repo ports.TechReadRepository, access *services.AccessControlService) *TechQueries {
	return &TechQueries{Repo: repo, Access: access}
}

func (q *TechQueries) ListNewTechItems(ctx cqrs.QueryContext, baseID int) ([]*readmodels.TechItemNew, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListNewTechItems(baseID)
}
func (q *TechQueries) ListInResearchTechItems(ctx cqrs.QueryContext, baseID int) ([]*readmodels.TechItemInProgress, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListInResearchTechItems(baseID)
}
func (q *TechQueries) ListDoneTechItems(ctx cqrs.QueryContext, baseID int) ([]*readmodels.TechItemDone, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListDoneTechItems(baseID)
}
