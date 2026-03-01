package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
)

type OperationQueries struct {
	Repo   ports.OperationReadRepository
	Access *services.AccessControlService
}

func NewOperationQueries(repo ports.OperationReadRepository, access *services.AccessControlService) *OperationQueries {
	return &OperationQueries{Repo: repo, Access: access}
}

func (q *OperationQueries) GetOperation(ctx context.Context, _ cqrs.Actor, operationID int) (*readmodels.MilitaryOperation, error) {
	op, err := q.Repo.GetOperation(ctx, operationID)
	return op, repoErr(err)
}
func (q *OperationQueries) ListOperationsByBase(ctx context.Context, actor cqrs.Actor, baseID int) ([]*readmodels.MilitaryOperation, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	ops, err := q.Repo.ListOperationsByBase(ctx, baseID)
	return ops, repoErr(err)
}
func (q *OperationQueries) ListActiveOperations(ctx context.Context, actor cqrs.Actor, baseID int) ([]*readmodels.MilitaryOperation, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	ops, err := q.Repo.ListActiveOperations(ctx, baseID)
	return ops, repoErr(err)
}
