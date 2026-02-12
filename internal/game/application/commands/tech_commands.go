package commands

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-api/internal/game/domain"
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

func (c *TechCommands) StartTechResearch(ctx cqrs.CommandContext, baseID int, prototypeID int) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		tRepo := c.TechRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		proto, err := tRepo.FindPrototypeByID(prototypeID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.StartTechResearch(proto); err != nil {
			return cqrs.NewDomainError(err)
		}
		if err := bRepo.Update(base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *TechCommands) SpeedUpTechResearchWithCrystals(ctx cqrs.CommandContext, baseID int, techItemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		uRepo := c.UserRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		user, err := uRepo.FindByIDForUpdate(ctx.UserID)
		if err != nil {
			return repoErr(err)
		}
		if err := c.crystalService.SpeedUpTechResearch(user, base, techItemID); err != nil {
			return cqrs.NewDomainError(err)
		}
		if err := uRepo.Update(user); err != nil {
			return err
		}
		if err := bRepo.Update(base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *TechCommands) HandleTechResearchStartedEvent(event *domain.TechResearchStartedEvent) error {
	cmd := ports.MoveTechQueueJob{BaseID: event.BaseID}
	return c.Scheduler.Schedule(cmd, event.CompletionDate)
}

func (c *TechCommands) HandleMoveTechQueueJob(cmd ports.MoveTechQueueJob) error {
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(cmd.BaseID)
		if err != nil {
			return err
		}
		base.MoveTechQueue()
		if err := bRepo.Update(base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}
