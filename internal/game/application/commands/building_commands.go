package commands

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

type BuildingCommands struct {
	BaseRepo       ports.UserBaseRepository
	BuildingRepo   ports.BuildPrototypeRepository
	Outbox         ports.OutboxEventRepository
	Scheduler      ports.Scheduler
	UserRepo       ports.UserRepository
	crystalService *domain.CrystalSpendingService
	TxMgr          ports.TransactionManager
	Access         *services.AccessControlService
}

func NewBuildingCommands(baseRepo ports.UserBaseRepository, buildingRepo ports.BuildPrototypeRepository, userRepo ports.UserRepository, outbox ports.OutboxEventRepository, scheduler ports.Scheduler, txMgr ports.TransactionManager, access *services.AccessControlService) *BuildingCommands {
	return &BuildingCommands{BaseRepo: baseRepo, BuildingRepo: buildingRepo, UserRepo: userRepo, crystalService: domain.NewCrystalSpendingService(), Outbox: outbox, Scheduler: scheduler, TxMgr: txMgr, Access: access}
}

func (c *BuildingCommands) QueueBuilding(ctx cqrs.CommandContext, baseID int, prototypeID int) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		biRepo := c.BuildingRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		proto, err := biRepo.FindPrototypeByID(prototypeID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.AddToBuildQueue(proto); err != nil {
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

func (c *BuildingCommands) CancelPendingBuilding(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.CancelPendingBuildingByID(itemID); err != nil {
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

func (c *BuildingCommands) SpeedUpProductionWithCrystals(ctx cqrs.CommandContext, baseID int, buildingItemID uuid.UUID) error {
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
		if err := c.crystalService.SpeedUpBuildingProduction(user, base, buildingItemID); err != nil {
			return err
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

func (c *BuildingCommands) DeletePresentBuilding(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.DeletePresentBuildingByID(itemID); err != nil {
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

func (c *BuildingCommands) HandleProductionStartedEvent(event *domain.BuildingProductionStartedEvent) error {
	cmd := ports.MoveBuildQueueJob{BaseID: event.BaseID}
	return c.Scheduler.Schedule(cmd, event.CompletionDate)
}

func (c *BuildingCommands) HandleMoveBuildQueueJob(cmd ports.MoveBuildQueueJob) error {
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(cmd.BaseID)
		if err != nil {
			return err
		}
		base.MoveBuildQueue()
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
