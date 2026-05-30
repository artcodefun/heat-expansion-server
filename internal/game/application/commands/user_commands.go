package commands

import (
	"context"
	"errors"
	"time"

	authv1 "github.com/artcodefun/heat-expansion-server/contracts/auth/events/v1"
	billingv1 "github.com/artcodefun/heat-expansion-server/contracts/billing/events/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type UserCommands struct {
	UserRepo       ports.UserRepository
	CrystalCredits ports.CrystalCreditsRepository
	Outbox         ports.OutboxEventRepository
	TxMgr          ports.TransactionManager
}

func NewUserCommands(userRepo ports.UserRepository, crystalCredits ports.CrystalCreditsRepository, outbox ports.OutboxEventRepository, txMgr ports.TransactionManager) *UserCommands {
	return &UserCommands{UserRepo: userRepo, CrystalCredits: crystalCredits, Outbox: outbox, TxMgr: txMgr}
}

func (c *UserCommands) HandleCrystalsPurchasedV1Event(ctx context.Context, ev billingv1.CrystalsPurchasedV1) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		creditsRepo := c.CrystalCredits.Tx(tx)

		exists, err := creditsRepo.Exists(ctx, ev.OrderID)
		if err != nil {
			return err
		}
		if exists {
			return nil // already credited, idempotent no-op
		}

		user, err := c.UserRepo.Tx(tx).FindByIDForUpdate(ctx, ev.UserID)
		if err != nil {
			return err
		}

		if err := user.AddCrystals(ev.Crystals); err != nil {
			return err
		}

		if err := c.UserRepo.Tx(tx).Update(ctx, user); err != nil {
			return err
		}

		if err := c.Outbox.Tx(tx).Save(ctx, user.EventProducer.PullEvents()); err != nil {
			return err
		}

		return creditsRepo.Insert(ctx, ev.OrderID, ev.UserID, ev.Crystals, time.Now().Unix())
	})
}

func (c *UserCommands) HandleAccountRegisteredV1Event(ctx context.Context, ev authv1.AccountRegisteredV1) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		uRepo := c.UserRepo.Tx(tx)

		// Check for idempotency: if user already exists, skip creation
		existing, err := uRepo.FindByID(ctx, ev.UserID)
		if err != nil && !errors.Is(err, ports.ErrNotFound) {
			return err
		}
		if existing != nil {
			return nil // User already created, event processed
		}

		user := domain.NewUser(ev.UserID, ev.Name)

		if err := uRepo.Create(ctx, user); err != nil {
			return err
		}

		if err := c.Outbox.Tx(tx).Save(ctx, user.EventProducer.PullEvents()); err != nil {
			return err
		}

		return nil
	})
}
