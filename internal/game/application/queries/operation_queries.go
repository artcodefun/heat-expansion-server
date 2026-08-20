package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/google/uuid"
)

type OperationQueries struct {
	Repo   ports.OperationReadRepository
	Access *services.AccessControlService
}

func NewOperationQueries(repo ports.OperationReadRepository, access *services.AccessControlService) *OperationQueries {
	return &OperationQueries{Repo: repo, Access: access}
}

func (q *OperationQueries) GetOperation(ctx context.Context, _ application.Actor, operationID int) (*readmodels.MilitaryOperation, error) {
	op, err := q.Repo.GetOperation(ctx, operationID)
	return op, repoErr(err)
}

func (q *OperationQueries) GetOperationByUUID(ctx context.Context, _ application.Actor, operationUUID uuid.UUID) (*readmodels.MilitaryOperation, error) {
	op, err := q.Repo.GetOperationByUUID(ctx, operationUUID)
	return op, repoErr(err)
}

func (q *OperationQueries) ListOperationsByBase(ctx context.Context, actor application.Actor, baseID int) ([]*readmodels.MilitaryOperation, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	ops, err := q.Repo.ListOperationsByBase(ctx, baseID)
	return ops, repoErr(err)
}
func (q *OperationQueries) ListActiveOperations(ctx context.Context, actor application.Actor, baseID int) ([]*readmodels.MilitaryOperation, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	ops, err := q.Repo.ListActiveOperations(ctx, baseID)
	return ops, repoErr(err)
}
