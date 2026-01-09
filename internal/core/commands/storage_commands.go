package commands

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
	"github.com/google/uuid"
)

type StorageCommands struct {
	BaseRepo  ports.UserBaseRepository
	Outbox    ports.OutboxEventRepository
	Scheduler ports.Scheduler
	TxMgr     ports.TransactionManager
	Access    *services.AccessControlService
}

func NewStorageCommands(baseRepo ports.UserBaseRepository, outbox ports.OutboxEventRepository, scheduler ports.Scheduler, txMgr ports.TransactionManager, access *services.AccessControlService) *StorageCommands {
	return &StorageCommands{BaseRepo: baseRepo, Outbox: outbox, Scheduler: scheduler, TxMgr: txMgr, Access: access}
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

func (c *StorageCommands) HandleBuffActivatedEvent(event *domain.BuffActivatedEvent) error {
	base, err := c.BaseRepo.FindByID(event.BaseID)
	if err != nil {
		return err
	}
	var activatedBuff *domain.StorageItemPresent
	for _, item := range base.StorageItemsPresent {
		if item.ID == event.ItemID && item.Prototype.BuffData != nil && item.ActivatedAt != nil {
			activatedBuff = &item
			break
		}
	}
	if activatedBuff == nil {
		return nil
	}
	expireAt := *activatedBuff.ActivatedAt + activatedBuff.Prototype.BuffData.DurationSeconds
	cmd := ports.DeleteExpiredBuffJob{BaseID: event.BaseID, ItemID: event.ItemID}
	return c.Scheduler.Schedule(cmd, expireAt)
}

func (c *StorageCommands) HandleDeleteExpiredBuffJob(baseID int) (int, error) {
	var count int
	err := c.TxMgr.WithTx(func(tx ports.Transaction) error {
		bRepo := c.BaseRepo.Tx(tx)
		base, err := bRepo.FindByIDForUpdate(baseID)
		if err != nil {
			return err
		}
		count = base.DeleteExpiredBuffs()
		if err := bRepo.Update(base); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(base.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return count, err
	}
	return count, nil
}
