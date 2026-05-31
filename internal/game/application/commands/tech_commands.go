package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

type TechCommands struct {
	BaseRepo       ports.UserBaseRepository
	TechRepo       ports.TechPrototypeRepository
	Outbox         ports.OutboxEventRepository
	Scheduler      ports.Scheduler
	UserRepo       ports.UserRepository
	crystalService *domain.CrystalSpendingService
	TxMgr          ports.TransactionManager
	Access         *services.AccessControlService
}

func NewTechCommands(baseRepo ports.UserBaseRepository, techRepo ports.TechPrototypeRepository, userRepo ports.UserRepository, outbox ports.OutboxEventRepository, scheduler ports.Scheduler, txMgr ports.TransactionManager, access *services.AccessControlService) *TechCommands {
	return &TechCommands{BaseRepo: baseRepo, TechRepo: techRepo, UserRepo: userRepo, crystalService: domain.NewCrystalSpendingService(), Outbox: outbox, Scheduler: scheduler, TxMgr: txMgr, Access: access}
}

func (c *TechCommands) StartTechResearch(ctx context.Context, actor cqrs.Actor, baseID int, prototypeID int) error {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		tRepo := c.TechRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(ctx, baseID)
		if err != nil {
			return repoErr(err)
		}
		proto, err := tRepo.FindPrototypeByID(ctx, prototypeID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.StartTechResearch(proto); err != nil {
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

func (c *TechCommands) SpeedUpTechResearchWithCrystals(ctx context.Context, actor cqrs.Actor, baseID int, techItemID uuid.UUID) error {
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
		user, err := uRepo.FindByIDForUpdate(ctx, actor.UserID)
		if err != nil {
			return repoErr(err)
		}
		if err := c.crystalService.SpeedUpTechResearch(user, base, techItemID); err != nil {
			return err
		}
		if err := uRepo.Update(ctx, user); err != nil {
			return err
		}
		if err := bRepo.Update(ctx, base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, user.EventProducer.PullEvents()); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(ctx, base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *TechCommands) HandleTechResearchStartedEvent(ctx context.Context, event domain.TechResearchStartedEvent) error {
	cmd := ports.MoveTechQueueJob{BaseID: event.BaseID}
	return c.Scheduler.Schedule(ctx, cmd, event.CompletionDate)
}

func (c *TechCommands) HandleMoveTechQueueJob(ctx context.Context, cmd ports.MoveTechQueueJob) error {
	err := c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(ctx, cmd.BaseID)
		if err != nil {
			return err
		}
		base.MoveTechQueue()
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
