package commands

import (
	"context"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

type TradeCommands struct {
	UserBaseRepo   ports.UserBaseRepository
	UserRepo       ports.UserRepository
	DiplomacyRepo  ports.DiplomaticRelationshipRepository
	TradeRepo      ports.TradeOperationRepository
	Scheduler      ports.Scheduler
	Outbox         ports.OutboxEventRepository
	TxMgr          ports.TransactionManager
	Access         *services.AccessControlService
	crystalService *domain.CrystalSpendingService
}

func NewTradeCommands(
	userBaseRepo ports.UserBaseRepository,
	userRepo ports.UserRepository,
	diplomacyRepo ports.DiplomaticRelationshipRepository,
	tradeRepo ports.TradeOperationRepository,
	scheduler ports.Scheduler,
	outbox ports.OutboxEventRepository,
	txMgr ports.TransactionManager,
	access *services.AccessControlService,
) *TradeCommands {
	return &TradeCommands{
		UserBaseRepo:   userBaseRepo,
		UserRepo:       userRepo,
		DiplomacyRepo:  diplomacyRepo,
		TradeRepo:      tradeRepo,
		Scheduler:      scheduler,
		Outbox:         outbox,
		TxMgr:          txMgr,
		Access:         access,
		crystalService: domain.NewCrystalSpendingService(),
	}
}

func (c *TradeCommands) CreateTradeOperation(
	ctx context.Context,
	actor cqrs.Actor,
	senderBaseID int,
	targetX, targetY int,
	offeredResources domain.PriceModel,
	offeredArmyRequests []domain.ArmyDeploymentRequest,
	offeredStorageItemIDs []uuid.UUID,
	requestedResources domain.PriceModel,
	requestedArmyRequests []domain.ArmyDeploymentRequest,
	requestedStorageItemIDs []uuid.UUID,
	transportRequests []domain.ArmyDeploymentRequest,
) (*domain.TradeOperation, error) {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, senderBaseID); err != nil {
		return nil, err
	}

	var created *domain.TradeOperation
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		bRepo := c.UserBaseRepo.Tx(tx)
		dRepo := c.DiplomacyRepo.Tx(tx)
		tRepo := c.TradeRepo.Tx(tx)

		sender, err := bRepo.FindByIDForUpdate(ctx, senderBaseID)
		if err != nil {
			return repoErr(err)
		}
		receiver, err := bRepo.FindByCoordinatesForUpdate(ctx, targetX, targetY)
		if err != nil {
			return repoErr(err)
		}

		rel, err := dRepo.FindBetweenUsers(ctx, sender.UserID, receiver.UserID)
		if err != nil {
			if !errors.Is(err, ports.ErrNotFound) {
				return repoErr(err)
			}
			rel, err = domain.NewUnknownRelationship(sender.UserID, receiver.UserID)
			if err != nil {
				return err
			}
		}
		if err := rel.CanPerformTradeOperation(); err != nil {
			return err
		}

		created, err = domain.BuildTradeOperationForCreation(
			sender,
			receiver,
			offeredResources,
			offeredArmyRequests,
			offeredStorageItemIDs,
			requestedResources,
			requestedArmyRequests,
			requestedStorageItemIDs,
			transportRequests,
		)
		if err != nil {
			return err
		}

		if err := tRepo.Create(ctx, created); err != nil {
			return repoErr(err)
		}
		created.EmitCreatedEvent()

		svc := domain.NewTradeOperationService(sender, receiver, created)
		if err := svc.CommitSenderForTradeCreation(); err != nil {
			return err
		}
		if err := bRepo.Update(ctx, sender); err != nil {
			return repoErr(err)
		}

		if err := c.Outbox.Tx(tx).Save(ctx, created.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return created, nil
}

// AcceptTradeOperation transitions a PENDING trade operation to OUTBOUND by committing
// the receiver side atomically: validating the requested payload, deducting resources,
// removing storage items, and deploying army units into the trade operation.
func (c *TradeCommands) AcceptTradeOperation(ctx context.Context, actor cqrs.Actor, operationID int) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		bRepo := c.UserBaseRepo.Tx(tx)
		tRepo := c.TradeRepo.Tx(tx)

		op, err := tRepo.FindByIDForUpdate(ctx, operationID)
		if err != nil {
			return repoErr(err)
		}

		if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, op.ReceiverBaseID); err != nil {
			return err
		}

		sender, err := bRepo.FindByID(ctx, op.SenderBaseID)
		if err != nil {
			return repoErr(err)
		}
		receiver, err := bRepo.FindByIDForUpdate(ctx, op.ReceiverBaseID)
		if err != nil {
			return repoErr(err)
		}

		svc := domain.NewTradeOperationService(sender, receiver, op)
		if err := svc.AcceptAndCommitReceiver(); err != nil {
			return err
		}

		if err := bRepo.Update(ctx, receiver); err != nil {
			return repoErr(err)
		}
		if err := tRepo.Update(ctx, op); err != nil {
			return repoErr(err)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, op.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		return nil
	})
}

// DeclineTradeOperation transitions a PENDING trade operation to COMPLETED/DECLINED.
// Sender-side restoration is driven by the TradeOperationReturnArrivedEvent emitted here
// and handled by HandleTradeOperationReturnArrivedEvent in a follow-up transaction.
func (c *TradeCommands) DeclineTradeOperation(ctx context.Context, actor cqrs.Actor, operationID int) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		tRepo := c.TradeRepo.Tx(tx)

		op, err := tRepo.FindByIDForUpdate(ctx, operationID)
		if err != nil {
			return repoErr(err)
		}

		if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, op.ReceiverBaseID); err != nil {
			return err
		}

		if err := op.Decline(); err != nil {
			return err
		}

		if err := tRepo.Update(ctx, op); err != nil {
			return repoErr(err)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, op.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		return nil
	})
}

// CancelTradeOperationByInitiator cancels a trade in PENDING or OUTBOUND phase at the
// sender's request. In OUTBOUND, receiver commitments are released in the same transaction;
// sender commitments are released when the convoy returns via TradeOperationReturnArrivedEvent.
func (c *TradeCommands) CancelTradeOperationByInitiator(ctx context.Context, actor cqrs.Actor, operationID int) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		bRepo := c.UserBaseRepo.Tx(tx)
		tRepo := c.TradeRepo.Tx(tx)

		op, err := tRepo.FindByIDForUpdate(ctx, operationID)
		if err != nil {
			return repoErr(err)
		}

		if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, op.SenderBaseID); err != nil {
			return err
		}

		sender, err := bRepo.FindByID(ctx, op.SenderBaseID)
		if err != nil {
			return repoErr(err)
		}
		receiver, err := bRepo.FindByIDForUpdate(ctx, op.ReceiverBaseID)
		if err != nil {
			return repoErr(err)
		}

		svc := domain.NewTradeOperationService(sender, receiver, op)
		if err := svc.CancelAndReleaseReceiverIfCommitted(); err != nil {
			return err
		}

		if err := bRepo.Update(ctx, receiver); err != nil {
			return repoErr(err)
		}
		if err := tRepo.Update(ctx, op); err != nil {
			return repoErr(err)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, op.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		return nil
	})
}

// HandleTradeOperationReturnArrivedEvent is the terminal release trigger for the sender.
// It is emitted for every terminal transition (decline, expire, cancel-pending, cancel-after-return,
// successful completion) and finalizes sender-side assets in one transaction.
func (c *TradeCommands) HandleTradeOperationReturnArrivedEvent(ctx context.Context, event domain.TradeOperationReturnArrivedEvent) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		tRepo := c.TradeRepo.Tx(tx)
		bRepo := c.UserBaseRepo.Tx(tx)

		op, err := tRepo.FindByID(ctx, event.OperationID)
		if err != nil {
			if errors.Is(err, ports.ErrNotFound) {
				return nil
			}
			return repoErr(err)
		}

		sender, err := bRepo.FindByIDForUpdate(ctx, op.SenderBaseID)
		if err != nil {
			return repoErr(err)
		}
		receiver, err := bRepo.FindByID(ctx, op.ReceiverBaseID)
		if err != nil {
			return repoErr(err)
		}

		svc := domain.NewTradeOperationService(sender, receiver, op)
		if err := svc.FinalizeSenderAfterReturn(); err != nil {
			return err
		}

		if err := bRepo.Update(ctx, sender); err != nil {
			return repoErr(err)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, sender.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		return nil
	})
}

func (c *TradeCommands) HandleTradeOperationCreatedEvent(ctx context.Context, event domain.TradeOperationCreatedEvent) error {
	return c.Scheduler.Schedule(ctx, ports.ExpireTradeOperationJob{OperationID: event.OperationID}, event.ExpirationAtSec)
}

func (c *TradeCommands) HandleTradeOperationOutboundEvent(ctx context.Context, event domain.TradeOperationOutboundEvent) error {
	return c.Scheduler.Schedule(ctx, ports.UpdateTradeOperationJob{OperationID: event.OperationID}, event.OutboundArriveAt)
}

func (c *TradeCommands) HandleTradeOperationArrivedEvent(ctx context.Context, event domain.TradeOperationArrivedEvent) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		tRepo := c.TradeRepo.Tx(tx)
		bRepo := c.UserBaseRepo.Tx(tx)

		op, err := tRepo.FindByIDForUpdate(ctx, event.OperationID)
		if err != nil {
			if errors.Is(err, ports.ErrNotFound) {
				return nil
			}
			return repoErr(err)
		}
		if op.Phase != domain.TradePhaseArrived {
			return nil
		}

		sender, err := bRepo.FindByIDForUpdate(ctx, op.SenderBaseID)
		if err != nil {
			return repoErr(err)
		}
		receiver, err := bRepo.FindByIDForUpdate(ctx, op.ReceiverBaseID)
		if err != nil {
			return repoErr(err)
		}

		svc := domain.NewTradeOperationService(sender, receiver, op)
		if err := svc.ProcessArrivalAndStartReturn(); err != nil {
			return err
		}

		if err := bRepo.Update(ctx, sender); err != nil {
			return repoErr(err)
		}
		if err := bRepo.Update(ctx, receiver); err != nil {
			return repoErr(err)
		}
		if err := tRepo.Update(ctx, op); err != nil {
			return repoErr(err)
		}

		if err := c.Outbox.Tx(tx).Save(ctx, sender.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, receiver.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, op.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		return nil
	})
}

func (c *TradeCommands) HandleTradeOperationReturningEvent(ctx context.Context, event domain.TradeOperationReturningEvent) error {
	return c.Scheduler.Schedule(ctx, ports.UpdateTradeOperationJob{OperationID: event.OperationID}, event.ReturnArriveAt)
}

func (c *TradeCommands) HandleUpdateTradeOperationJob(ctx context.Context, cmd ports.UpdateTradeOperationJob) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		tRepo := c.TradeRepo.Tx(tx)

		op, err := tRepo.FindByIDForUpdate(ctx, cmd.OperationID)
		if err != nil {
			if errors.Is(err, ports.ErrNotFound) {
				return nil
			}
			return repoErr(err)
		}
		op.UpdatePhaseBasedOnTime()

		if err := tRepo.Update(ctx, op); err != nil {
			return repoErr(err)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, op.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		return nil
	})
}

// SpeedUpTradeOperationWithCrystals allows the sender to spend crystals to fast-forward
// an in-flight trade operation (outbound or returning) to its arrival.
func (c *TradeCommands) SpeedUpTradeOperationWithCrystals(ctx context.Context, actor cqrs.Actor, operationID int) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		tRepo := c.TradeRepo.Tx(tx)
		uRepo := c.UserRepo.Tx(tx)

		op, err := tRepo.FindByIDForUpdate(ctx, operationID)
		if err != nil {
			return repoErr(err)
		}

		if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, op.SenderBaseID); err != nil {
			return err
		}

		user, err := uRepo.FindByIDForUpdate(ctx, actor.UserID)
		if err != nil {
			return repoErr(err)
		}

		if err := c.crystalService.SpeedUpTradeOperation(user, op); err != nil {
			return err
		}

		if err := uRepo.Update(ctx, user); err != nil {
			return repoErr(err)
		}
		if err := tRepo.Update(ctx, op); err != nil {
			return repoErr(err)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, op.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		return nil
	})
}

func (c *TradeCommands) HandleExpireTradeOperationJob(ctx context.Context, cmd ports.ExpireTradeOperationJob) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		tRepo := c.TradeRepo.Tx(tx)
		op, err := tRepo.FindByIDForUpdate(ctx, cmd.OperationID)
		if err != nil {
			if errors.Is(err, ports.ErrNotFound) {
				return nil
			}
			return repoErr(err)
		}

		if !op.ExpireIfPending() {
			return nil
		}
		if err := tRepo.Update(ctx, op); err != nil {
			return repoErr(err)
		}
		if err := c.Outbox.Tx(tx).Save(ctx, op.EventProducer.PullEvents()); err != nil {
			return repoErr(err)
		}
		return nil
	})
}
