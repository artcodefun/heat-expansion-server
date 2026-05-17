package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/services"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type BlackMarketCommands struct {
	UserRepo    ports.UserRepository
	BaseRepo    ports.UserBaseRepository
	TxMgr       ports.TransactionManager
	Access      *services.AccessControlService
	BlackMarket *domain.BlackMarketService
}

func NewBlackMarketCommands(userRepo ports.UserRepository, baseRepo ports.UserBaseRepository, txMgr ports.TransactionManager, access *services.AccessControlService) *BlackMarketCommands {
	return &BlackMarketCommands{UserRepo: userRepo, BaseRepo: baseRepo, TxMgr: txMgr, Access: access, BlackMarket: domain.NewBlackMarketService()}
}

func (c *BlackMarketCommands) PurchaseResources(ctx context.Context, actor cqrs.Actor, baseID int, resourceType domain.ResourceType, crystals int) error {
	if err := c.Access.EnsureBaseOwnership(ctx, actor.UserID, baseID); err != nil {
		return err
	}

	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		uRepo := c.UserRepo.Tx(tx)
		bRepo := c.BaseRepo.Tx(tx)

		user, err := uRepo.FindByIDForUpdate(ctx, actor.UserID)
		if err != nil {
			return repoErr(err)
		}
		base, err := bRepo.FindByIDForUpdate(ctx, baseID)
		if err != nil {
			return repoErr(err)
		}

		if err := c.BlackMarket.PurchaseResources(user, base, resourceType, crystals); err != nil {
			return err
		}
		if err := uRepo.Update(ctx, user); err != nil {
			return repoErr(err)
		}
		if err := bRepo.Update(ctx, base); err != nil {
			return repoErr(err)
		}
		return nil
	})
}
