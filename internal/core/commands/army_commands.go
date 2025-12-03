package commands

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
	"github.com/google/uuid"
)

type ArmyCommands struct {
	BaseRepo       ports.UserBaseRepository
	ArmyRepo       ports.ArmyPrototypeRepository
	EventPublisher ports.EventPublisher
	Scheduler      ports.Scheduler
	UserRepo       ports.UserRepository
	crystalService *domain.CrystalSpendingService
	TxMgr          ports.TransactionManager
	Access         *services.AccessControlService
}

func NewArmyCommands(baseRepo ports.UserBaseRepository, armyRepo ports.ArmyPrototypeRepository, userRepo ports.UserRepository, eventPublisher ports.EventPublisher, scheduler ports.Scheduler, txMgr ports.TransactionManager, access *services.AccessControlService) *ArmyCommands {
	return &ArmyCommands{BaseRepo: baseRepo, ArmyRepo: armyRepo, UserRepo: userRepo, crystalService: domain.NewCrystalSpendingService(), EventPublisher: eventPublisher, Scheduler: scheduler, TxMgr: txMgr, Access: access}
}

func (c *ArmyCommands) QueueArmy(ctx cqrs.CommandContext, baseID int, prototypeID int, count int) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	var events []domain.DomainEvent
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		aRepo := c.ArmyRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		proto, err := aRepo.FindPrototypeByID(prototypeID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.QueueArmy(proto, count); err != nil {
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

func (c *ArmyCommands) CancelPendingArmy(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID, count int) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	var events []domain.DomainEvent
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.CancelPendingArmyByID(itemID, count); err != nil {
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

func (c *ArmyCommands) SpeedUpArmyProductionWithCrystals(ctx cqrs.CommandContext, baseID int, userID int, armyItemID uuid.UUID) error {
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
		user, err := uRepo.FindByID(userID)
		if err != nil {
			return repoErr(err)
		}
		if err := c.crystalService.SpeedUpArmyProduction(user, base, armyItemID); err != nil {
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

func (c *ArmyCommands) DeletePresentArmy(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID, count int) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	var events []domain.DomainEvent
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.DeletePresentArmyByID(itemID, count); err != nil {
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

func (c *ArmyCommands) HandleProductionStartedEvent(event *domain.ArmyProductionStartedEvent) error {
	cmd := ports.MoveArmyQueueJob{BaseID: event.BaseID}
	return c.Scheduler.Schedule(cmd, event.CompletionDate)
}

func (c *ArmyCommands) HandleMoveArmyQueueJob(cmd ports.MoveArmyQueueJob) error {
	var events []domain.DomainEvent
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByID(cmd.BaseID)
		if err != nil {
			return err
		}
		base.MoveArmyQueue()
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
