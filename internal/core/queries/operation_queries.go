package queries

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

type OperationQueries struct {
	Repo   ports.OperationReadRepository
	Access *services.AccessControlService
}

func NewOperationQueries(repo ports.OperationReadRepository, access *services.AccessControlService) *OperationQueries {
	return &OperationQueries{Repo: repo, Access: access}
}

func (q *OperationQueries) GetOperation(_ cqrs.QueryContext, operationID int) (*readmodels.MilitaryOperation, error) {
	op, err := q.Repo.GetOperation(operationID)
	return op, repoErr(err)
}
func (q *OperationQueries) ListOperationsByBase(ctx cqrs.QueryContext, baseID int) ([]*readmodels.MilitaryOperation, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	ops, err := q.Repo.ListOperationsByBase(baseID)
	return ops, repoErr(err)
}
func (q *OperationQueries) ListActiveOperations(ctx cqrs.QueryContext, baseID int) ([]*readmodels.MilitaryOperation, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	ops, err := q.Repo.ListActiveOperations(baseID)
	return ops, repoErr(err)
}
