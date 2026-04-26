package queries

import (
	"context"
	"errors"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/google/uuid"
)

type TradeQueries struct {
	ArmyRepo      ports.ArmyReadRepository
	StorageRepo   ports.StorageReadRepository
	BaseReadRepo  ports.BaseReadRepository
	BaseRepo      ports.UserBaseRepository
	DiplomacyRepo ports.DiplomacyReadRepository
	TradeRepo     ports.TradeOperationReadRepository
	Access        *services.AccessControlService
}

func NewTradeQueries(
	armyRepo ports.ArmyReadRepository,
	storageRepo ports.StorageReadRepository,
	baseReadRepo ports.BaseReadRepository,
	baseRepo ports.UserBaseRepository,
	diplomacyRepo ports.DiplomacyReadRepository,
	tradeRepo ports.TradeOperationReadRepository,
	access *services.AccessControlService,
) *TradeQueries {
	return &TradeQueries{
		ArmyRepo:      armyRepo,
		StorageRepo:   storageRepo,
		BaseReadRepo:  baseReadRepo,
		BaseRepo:      baseRepo,
		DiplomacyRepo: diplomacyRepo,
		TradeRepo:     tradeRepo,
		Access:        access,
	}
}

func (q *TradeQueries) GetTradeInfo(ctx context.Context, actor cqrs.Actor, targetX, targetY int) (*readmodels.TradeInfo, error) {
	targetBase, err := q.BaseRepo.FindByCoordinates(ctx, targetX, targetY)
	if err != nil {
		return nil, repoErr(err)
	}

	if err := q.ensureTradeInventoryReadAccess(ctx, actor.UserID, targetBase.ID); err != nil {
		return nil, err
	}

	army, err := q.ArmyRepo.ListPresentArmyItemsAll(ctx, targetBase.ID)
	if err != nil {
		return nil, repoErr(err)
	}

	storage, err := q.StorageRepo.ListTradeableStorageItems(ctx, targetBase.ID)
	if err != nil {
		return nil, repoErr(err)
	}

	stats, err := q.BaseReadRepo.GetBaseStats(ctx, targetBase.ID)
	if err != nil {
		return nil, repoErr(err)
	}
	resources := stats.CurrentResources(time.Now().Unix())

	return &readmodels.TradeInfo{
		Resources: resources,
		Army:      army,
		Storage:   storage,
	}, nil
}

func (q *TradeQueries) GetTradeOperation(ctx context.Context, actor cqrs.Actor, baseID int, operationID int) (*readmodels.TradeOperation, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	op, err := q.TradeRepo.GetTradeOperation(ctx, operationID)
	if err != nil {
		return nil, repoErr(err)
	}
	if op.SenderBaseID != baseID && op.ReceiverBaseID != baseID {
		return nil, cqrs.ErrForbidden
	}
	return op, nil
}

func (q *TradeQueries) ListActiveTradeOperations(ctx context.Context, actor cqrs.Actor, baseID int) ([]*readmodels.TradeOperation, error) {
	if err := q.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return nil, err
	}
	items, err := q.TradeRepo.ListActiveTradeOperations(ctx, baseID)
	return items, repoErr(err)
}

func (q *TradeQueries) ensureTradeInventoryReadAccess(ctx context.Context, actorUserID uuid.UUID, baseID int) error {
	if actorUserID == uuid.Nil {
		return cqrs.ErrForbidden
	}

	ownerUserID, err := q.BaseRepo.GetOwnerID(ctx, baseID)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return cqrs.ErrNotFound
		}
		return err
	}

	if ownerUserID == actorUserID {
		return nil
	}

	rel, err := q.DiplomacyRepo.GetRelationship(ctx, actorUserID, ownerUserID)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return cqrs.ErrForbidden
		}
		return err
	}

	if rel.Status != readmodels.DiplomaticStatusAllied {
		return cqrs.ErrForbidden
	}

	return nil
}
