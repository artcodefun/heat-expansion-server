package commands

import (
	"time"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
	"github.com/google/uuid"
)

type BuildingCommands struct {
	BaseRepo       ports.UserBaseRepository
	BuildingRepo   ports.BuildPrototypeRepository
	EventPublisher ports.EventPublisher
	Scheduler      ports.Scheduler
	UserRepo       ports.UserRepository
	crystalService *domain.CrystalSpendingService
	TxMgr          ports.TransactionManager
	Access         *services.AccessControlService
}

func NewBuildingCommands(baseRepo ports.UserBaseRepository, buildingRepo ports.BuildPrototypeRepository, userRepo ports.UserRepository, eventPublisher ports.EventPublisher, scheduler ports.Scheduler, txMgr ports.TransactionManager, access *services.AccessControlService) *BuildingCommands {
	return &BuildingCommands{BaseRepo: baseRepo, BuildingRepo: buildingRepo, UserRepo: userRepo, crystalService: domain.NewCrystalSpendingService(), EventPublisher: eventPublisher, Scheduler: scheduler, TxMgr: txMgr, Access: access}
}

func (c *BuildingCommands) QueueBuilding(ctx cqrs.CommandContext, baseID int, prototypeID int) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	var events []domain.DomainEvent
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
		events = append(events, base.EventProducer.PullEvents()...)
		return nil
	})
	if err != nil {
		return err
	}
	publishEvents(events, c.EventPublisher)
	return nil
}

func (c *BuildingCommands) CancelPendingBuilding(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID) error {
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
		if err := base.CancelPendingBuildingByID(itemID); err != nil {
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

func (c *BuildingCommands) SpeedUpProductionWithCrystals(ctx cqrs.CommandContext, baseID int, userID int, buildingItemID uuid.UUID) error {
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
		if err := c.crystalService.SpeedUpBuildingProduction(user, base, buildingItemID); err != nil {
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

func (c *BuildingCommands) DeletePresentBuilding(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID) error {
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
		if err := base.DeletePresentBuildingByID(itemID); err != nil {
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

func (c *BuildingCommands) HandleProductionStartedEvent(event *domain.BuildingProductionStartedEvent) error {
	cmd := ports.MoveBuildQueueJob{BaseID: event.BaseID}
	return c.Scheduler.Schedule(cmd, event.CompletionDate)
}

func (c *BuildingCommands) HandleProductionFinishedEvent(event *domain.BuildingProductionFinishedEvent) error {
	base, err := c.BaseRepo.FindByID(event.BaseID)
	if err != nil {
		return nil
	}
	for _, b := range base.BuildingsPresent {
		if b.ID == event.PresentItemID {
			if b.Prototype.IntelligenceData != nil && b.Prototype.IntelligenceData.Subtype == domain.IntelligenceSubtypeRadar {
				cooldown := b.Prototype.IntelligenceData.ScanCooldown
				if cooldown <= 0 {
					cooldown = 3600
				}
				firstAt := time.Now().Unix() + cooldown
				_ = c.Scheduler.Schedule(ports.RadarScanJob{BaseID: event.BaseID, BuildingID: b.ID}, firstAt)
			}
			break
		}
	}
	return nil
}

func (c *BuildingCommands) HandleMoveBuildQueueJob(cmd ports.MoveBuildQueueJob) error {
	var events []domain.DomainEvent
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
		events = append(events, base.EventProducer.PullEvents()...)
		return nil
	})
	if err != nil {
		return err
	}
	publishEvents(events, c.EventPublisher)
	return nil
}
