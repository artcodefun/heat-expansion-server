package commands

import (
	"context"
	"errors"

	v1 "github.com/artcodefun/heat-expansion-server/contracts/auth/v1"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type UserCommands struct {
	UserRepo ports.UserRepository
	Outbox   ports.OutboxEventRepository
	TxMgr    ports.TransactionManager
}

func NewUserCommands(userRepo ports.UserRepository, outbox ports.OutboxEventRepository, txMgr ports.TransactionManager) *UserCommands {
	return &UserCommands{UserRepo: userRepo, Outbox: outbox, TxMgr: txMgr}
}

func (c *UserCommands) HandleAccountRegisteredV1Event(ctx context.Context, ev v1.AccountRegisteredV1) error {
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
