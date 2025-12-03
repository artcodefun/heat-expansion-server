package commands

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
	"github.com/google/uuid"
)

type TechCommands struct {
	BaseRepo       ports.UserBaseRepository
	TechRepo       ports.TechPrototypeRepository
	EventPublisher ports.EventPublisher
	Scheduler      ports.Scheduler
	UserRepo       ports.UserRepository
	crystalService *domain.CrystalSpendingService
	TxMgr          ports.TransactionManager
	Access         *services.AccessControlService
}

func NewTechCommands(baseRepo ports.UserBaseRepository, techRepo ports.TechPrototypeRepository, userRepo ports.UserRepository, eventPublisher ports.EventPublisher, scheduler ports.Scheduler, txMgr ports.TransactionManager, access *services.AccessControlService) *TechCommands {
	return &TechCommands{BaseRepo: baseRepo, TechRepo: techRepo, UserRepo: userRepo, crystalService: domain.NewCrystalSpendingService(), EventPublisher: eventPublisher, Scheduler: scheduler, TxMgr: txMgr, Access: access}
}

func (c *TechCommands) StartTechResearch(ctx cqrs.CommandContext, baseID int, prototypeID int) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	var events []domain.DomainEvent
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
			return err
		}
		if err := bRepo.Update(base); err != nil {
			return err
		}
		events = append(events, base.EventProducer.PullEvents()...)
		return nil
	})
	if err != nil {
		return err
	}
	publishEvents(events, c.EventPublisher)
	return nil
}

func (c *TechCommands) SpeedUpTechResearchWithCrystals(ctx cqrs.CommandContext, baseID int, userID int, techItemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	var events []domain.DomainEvent
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		uRepo := c.UserRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		user, err := uRepo.FindByIDForUpdate(userID)
		if err != nil {
			return repoErr(err)
		}
		if err := c.crystalService.SpeedUpTechResearch(user, base, techItemID); err != nil {
			return err
		}
		if err := uRepo.Update(user); err != nil {
			return err
		}
		if err := bRepo.Update(base); err != nil {
			return err
		}
		events = append(events, base.EventProducer.PullEvents()...)
		return nil
	})
	if err != nil {
		return err
	}
	publishEvents(events, c.EventPublisher)
	return nil
}

func (c *TechCommands) HandleTechResearchStartedEvent(event *domain.TechResearchStartedEvent) error {
	cmd := ports.MoveTechQueueJob{BaseID: event.BaseID}
	return c.Scheduler.Schedule(cmd, event.CompletionDate)
}

func (c *TechCommands) HandleMoveTechQueueJob(cmd ports.MoveTechQueueJob) error {
	var events []domain.DomainEvent
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
		events = append(events, base.EventProducer.PullEvents()...)
		return nil
	})
	if err != nil {
		return err
	}
	publishEvents(events, c.EventPublisher)
	return nil
}
