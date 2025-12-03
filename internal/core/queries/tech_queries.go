package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type TechQueries struct {
	Repo      ports.TechReadRepository
	ProtoRepo ports.TechPrototypeRepository
	BaseRepo  ports.UserBaseRepository
	Access    *services.AccessControlService
}

func NewTechQueries(repo ports.TechReadRepository, protoRepo ports.TechPrototypeRepository, baseRepo ports.UserBaseRepository, access *services.AccessControlService) *TechQueries {
	return &TechQueries{Repo: repo, ProtoRepo: protoRepo, BaseRepo: baseRepo, Access: access}
}

func (q *TechQueries) ListNewTechItems(ctx cqrs.QueryContext, baseID int) ([]*readmodels.TechItemNew, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	allProtos, err := q.ProtoRepo.FindAllPrototypes()
	if err != nil {
		return nil, err
	}
	base, err := q.BaseRepo.FindByID(baseID)
	if err != nil {
		return nil, err
	}
	available := base.AvailableTechnologies(allProtos)
	ids := make([]int, 0, len(available))
	for _, p := range available {
		ids = append(ids, p.ID)
	}
	return q.Repo.ListNewTechItemsByPrototypeIDs(ids)
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
