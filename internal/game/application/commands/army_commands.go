package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

type ArmyCommands struct {
	BaseRepo       ports.UserBaseRepository
	ArmyRepo       ports.ArmyPrototypeRepository
	Outbox         ports.OutboxEventRepository
	Scheduler      ports.Scheduler
	UserRepo       ports.UserRepository
	crystalService *domain.CrystalSpendingService
	TxMgr          ports.TransactionManager
	Access         *services.AccessControlService
}

func NewArmyCommands(baseRepo ports.UserBaseRepository, armyRepo ports.ArmyPrototypeRepository, userRepo ports.UserRepository, outbox ports.OutboxEventRepository, scheduler ports.Scheduler, txMgr ports.TransactionManager, access *services.AccessControlService) *ArmyCommands {
	return &ArmyCommands{BaseRepo: baseRepo, ArmyRepo: armyRepo, UserRepo: userRepo, crystalService: domain.NewCrystalSpendingService(), Outbox: outbox, Scheduler: scheduler, TxMgr: txMgr, Access: access}
}

func (c *ArmyCommands) QueueArmy(ctx context.Context, actor cqrs.Actor, baseID int, prototypeID int, count int) error {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		aRepo := c.ArmyRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(ctx, baseID)
		if err != nil {
			return repoErr(err)
		}
		proto, err := aRepo.FindPrototypeByID(ctx, prototypeID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.QueueArmy(proto, count); err != nil {
			return err
		}
		if err := bRepo.Update(ctx, base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *ArmyCommands) CancelPendingArmy(ctx context.Context, actor cqrs.Actor, baseID int, itemID uuid.UUID, count int) error {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(ctx, baseID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.CancelPendingArmyByID(itemID, count); err != nil {
			return err
		}
		if err := bRepo.Update(ctx, base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *ArmyCommands) SpeedUpArmyProductionWithCrystals(ctx context.Context, actor cqrs.Actor, baseID int, armyItemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		uRepo := c.UserRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(ctx, baseID)
		if err != nil {
			return repoErr(err)
		}
		user, err := uRepo.FindByID(ctx, actor.UserID)
		if err != nil {
			return repoErr(err)
		}
		if err := c.crystalService.SpeedUpArmyProduction(user, base, armyItemID); err != nil {
			return err
		}
		if err := uRepo.Update(ctx, user); err != nil {
			return err
		}
		if err := bRepo.Update(ctx, base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *ArmyCommands) DeletePresentArmy(ctx context.Context, actor cqrs.Actor, baseID int, itemID uuid.UUID, count int) error {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(ctx, baseID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.DeletePresentArmyByID(itemID, count); err != nil {
			return err
		}
		if err := bRepo.Update(ctx, base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *ArmyCommands) HandleProductionStartedEvent(ctx context.Context, event domain.ArmyProductionStartedEvent) error {
	cmd := ports.MoveArmyQueueJob{BaseID: event.BaseID}
	return c.Scheduler.Schedule(ctx, cmd, event.CompletionDate)
}

func (c *ArmyCommands) HandleMoveArmyQueueJob(ctx context.Context, cmd ports.MoveArmyQueueJob) error {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByID(ctx, cmd.BaseID)
		if err != nil {
			return err
		}
		base.MoveArmyQueue()
		if err := bRepo.Update(ctx, base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}
