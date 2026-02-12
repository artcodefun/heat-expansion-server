package commands

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-api/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-api/internal/game/domain"
	"github.com/google/uuid"
)

type StorageCommands struct {
	BaseRepo           ports.UserBaseRepository
	UserRepo           ports.UserRepository
	SectorRepo         ports.SectorRepository
	StoragePrototypes  ports.StoragePrototypeRepository
	ArmyPrototypes     ports.ArmyPrototypeRepository
	ResourceLocations  ports.ResourceLocationRepository
	DangerousLocations ports.DangerousLocationRepository
	ScanReports        ports.ScanReportRepository
	Outbox             ports.OutboxEventRepository
	Scheduler          ports.Scheduler
	TxMgr              ports.TransactionManager
	Access             *services.AccessControlService
	rewards            *domain.ConsumableRewardService
}

func NewStorageCommands(
	baseRepo ports.UserBaseRepository,
	userRepo ports.UserRepository,
	sectorRepo ports.SectorRepository,
	storageProtos ports.StoragePrototypeRepository,
	armyProtos ports.ArmyPrototypeRepository,
	resourceLocs ports.ResourceLocationRepository,
	dangerousLocs ports.DangerousLocationRepository,
	scanReports ports.ScanReportRepository,
	outbox ports.OutboxEventRepository,
	scheduler ports.Scheduler,
	txMgr ports.TransactionManager,
	access *services.AccessControlService,
) *StorageCommands {
	return &StorageCommands{
		BaseRepo:           baseRepo,
		UserRepo:           userRepo,
		SectorRepo:         sectorRepo,
		StoragePrototypes:  storageProtos,
		ArmyPrototypes:     armyProtos,
		ResourceLocations:  resourceLocs,
		DangerousLocations: dangerousLocs,
		ScanReports:        scanReports,
		Outbox:             outbox,
		Scheduler:          scheduler,
		TxMgr:              txMgr,
		Access:             access,
		rewards:            domain.NewConsumableRewardService(),
	}
}

func (c *StorageCommands) DeletePresentStorageItem(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		if err := base.DeletePresentStorageItemByID(itemID); err != nil {
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

func (c *StorageCommands) ActivateBuff(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		if err = base.ActivateBuffByID(itemID); err != nil {
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

func (c *StorageCommands) StartIntelDecryption(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		if err = base.StartIntelDecryptionByID(itemID); err != nil {
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

func (c *StorageCommands) StartDamagedItemRestoration(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	armyProtos, err := c.ArmyPrototypes.FindAllPrototypes()
	if err != nil {
		return err
	}
	err = c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		if err = base.StartDamagedItemRestorationByID(itemID, armyProtos); err != nil {
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

func (c *StorageCommands) ActivateArtifact(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		if err = base.ActivateArtifactByID(itemID); err != nil {
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

func (c *StorageCommands) DeactivateArtifact(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return repoErr(err)
		}
		if err = base.DeactivateArtifactByID(itemID); err != nil {
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

func (c *StorageCommands) OpenConsumableBox(ctx cqrs.CommandContext, baseID int, itemID uuid.UUID) error {
	if err := c.Access.EnsureBaseOwnership(ctx.UserID, baseID); err != nil {
		return err
	}

	// Load all prototypes for the reward service
	allProtos, err := c.StoragePrototypes.FindAllPrototypes()
	if err != nil {
		return err
	}

	// Filter them by category
	var buffProtos, intelProtos, damagedProtos, artifactProtos []domain.StorageItemPrototype
	for _, p := range allProtos {
		switch p.Category {
		case domain.StorageCategoryBuff:
			buffProtos = append(buffProtos, *p)
		case domain.StorageCategoryIntel:
			intelProtos = append(intelProtos, *p)
		case domain.StorageCategoryDamaged:
			damagedProtos = append(damagedProtos, *p)
		case domain.StorageCategoryArtifact:
			artifactProtos = append(artifactProtos, *p)
		}
	}

	err = c.TxMgr.WithTx(func(tx ports.Transaction) error {
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

		_, err = c.rewards.OpenBox(base, user, itemID, buffProtos, intelProtos, damagedProtos, artifactProtos)
		if err != nil {
			return cqrs.NewDomainError(err)
		}

		if err := bRepo.Update(base); err != nil {
			return err
		}
		if err := uRepo.Update(user); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (c *StorageCommands) HandleBuffActivatedEvent(event *domain.BuffActivatedEvent) error {
	base, err := c.BaseRepo.FindByID(event.BaseID)
	if err != nil {
		return err
	}
	var activatedBuff *domain.StorageItemPresent
	for _, item := range base.StorageItemsPresent {
		if item.ID == event.ItemID && item.Prototype.BuffData != nil && item.ExpiresAt != nil {
			activatedBuff = &item
			break
		}
	}
	if activatedBuff == nil {
		return nil
	}
	cmd := ports.DeleteExpiredBuffJob{BaseID: event.BaseID, ItemID: event.ItemID}
	return c.Scheduler.Schedule(cmd, *activatedBuff.ExpiresAt)
}

func (c *StorageCommands) HandleIntelDecryptionStartedEvent(event *domain.IntelDecryptionStartedEvent) error {
	base, err := c.BaseRepo.FindByID(event.BaseID)
	if err != nil {
		return err
	}
	var intelItem *domain.StorageItemPresent
	for _, item := range base.StorageItemsPresent {
		if item.ID == event.ItemID && item.Prototype.IntelData != nil && item.ExpiresAt != nil {
			intelItem = &item
			break
		}
	}
	if intelItem == nil {
		return nil
	}
	cmd := ports.DecryptIntelItemJob{BaseID: event.BaseID, ItemID: event.ItemID}
	return c.Scheduler.Schedule(cmd, *intelItem.ExpiresAt)
}

func (c *StorageCommands) HandleDamagedItemRestorationStartedEvent(event *domain.DamagedItemRestorationStartedEvent) error {
	base, err := c.BaseRepo.FindByID(event.BaseID)
	if err != nil {
		return err
	}
	var damagedItem *domain.StorageItemPresent
	for _, item := range base.StorageItemsPresent {
		if item.ID == event.ItemID && item.Prototype.DamagedData != nil && item.ExpiresAt != nil {
			damagedItem = &item
			break
		}
	}
	if damagedItem == nil {
		return nil
	}
	cmd := ports.RestoreDamagedItemJob{BaseID: event.BaseID, ItemID: event.ItemID}
	return c.Scheduler.Schedule(cmd, *damagedItem.ExpiresAt)
}

func (c *StorageCommands) HandleDeleteExpiredBuffJob(cmd ports.DeleteExpiredBuffJob) error {
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(cmd.BaseID)
		if err != nil {
			return err
		}
		base.DeleteExpiredBuffs()
		if err := bRepo.Update(base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
}

func (c *StorageCommands) HandleDecryptIntelItemJob(cmd ports.DecryptIntelItemJob) error {
	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(cmd.BaseID)
		if err != nil {
			return err
		}
		intelType, err := base.DecryptIntelItemByID(cmd.ItemID)
		if err != nil {
			return err
		}

		// Find closest location of the revealed type
		var report *domain.SectorScanReport
		switch intelType {
		case domain.HiddenLocationTypeResourceful:
			loc, err := c.ResourceLocations.Tx(tx).FindClosest(base.Coordinates.X, base.Coordinates.Y)
			if err == nil {
				sector, _ := c.SectorRepo.Tx(tx).FindByCoordinates(loc.Coordinates.X, loc.Coordinates.Y)
				if sector != nil {
					report = domain.NewSectorScanReportFromResourceLocation(base.ID, sector, loc)
				}
			}
		case domain.HiddenLocationTypeDangerous:
			loc, err := c.DangerousLocations.Tx(tx).FindClosest(base.Coordinates.X, base.Coordinates.Y)
			if err == nil {
				sector, _ := c.SectorRepo.Tx(tx).FindByCoordinates(loc.Coordinates.X, loc.Coordinates.Y)
				if sector != nil {
					report = domain.NewSectorScanReportFromDangerousLocation(base.ID, sector, loc)
				}
			}
		case domain.HiddenLocationTypeUserBase:
			target, err := c.BaseRepo.Tx(tx).FindClosest(base.Coordinates.X, base.Coordinates.Y)
			if err == nil {
				sector, _ := c.SectorRepo.Tx(tx).FindByCoordinates(target.Coordinates.X, target.Coordinates.Y)
				if sector != nil {
					report = domain.NewSectorScanReportFromUserBase(base.ID, sector, target)
				}
			}
		}

		if report != nil {
			report.SourceIntelItemID = &cmd.ItemID
			if err := c.ScanReports.Tx(tx).Create(report); err != nil {
				return err
			}
			report.EmitCreated()
			if err := c.Outbox.Tx(tx).Save(report.PullEvents()); err != nil {
				return err
			}
		}

		if err := bRepo.Update(base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
}

func (c *StorageCommands) HandleRestoreDamagedItemJob(cmd ports.RestoreDamagedItemJob) error {
	armyProtos, err := c.ArmyPrototypes.FindAllPrototypes()
	if err != nil {
		return err
	}

	return c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(cmd.BaseID)
		if err != nil {
			return err
		}
		if err := base.RestoreDamagedItemByID(cmd.ItemID, armyProtos); err != nil {
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
}
