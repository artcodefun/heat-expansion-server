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
	return q.Repo.GetOperation(operationID)
}
func (q *OperationQueries) ListOperationsByBase(ctx cqrs.QueryContext, baseID int) ([]*readmodels.MilitaryOperation, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListOperationsByBase(baseID)
}
func (q *OperationQueries) ListActiveOperations(ctx cqrs.QueryContext, baseID int) ([]*readmodels.MilitaryOperation, error) {
	if err := q.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return nil, err
	}
	return q.Repo.ListActiveOperations(baseID)
}
