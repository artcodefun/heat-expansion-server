package commands

import (
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/core/ports"
)

type UserCommands struct {
	UserRepo       ports.UserRepository
	PasswordHasher ports.PasswordHasher
	TokenProvider  ports.TokenProvider
	Outbox         ports.OutboxEventRepository
	TxMgr          ports.TransactionManager
}

func NewUserCommands(userRepo ports.UserRepository, hasher ports.PasswordHasher, tokenProvider ports.TokenProvider, outbox ports.OutboxEventRepository, txMgr ports.TransactionManager) *UserCommands {
	return &UserCommands{UserRepo: userRepo, PasswordHasher: hasher, TokenProvider: tokenProvider, Outbox: outbox, TxMgr: txMgr}
}

func (c *UserCommands) Authenticate(ctx cqrs.CommandContext, email, password string) (string, error) {
	user, err := c.UserRepo.FindByEmail(email)
	if err != nil {
		return "", repoErr(err)
	}
	if !c.PasswordHasher.Verify(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}
	token, err := c.TokenProvider.Generate(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (c *UserCommands) Create(ctx cqrs.CommandContext, name, email, password string) error {
	hashed, err := c.PasswordHasher.Hash(password)
	if err != nil {
		return err
	}
	err = c.TxMgr.WithTx(func(tx ports.Transaction) error {
		uRepo := c.UserRepo.Tx(tx)
		user := &domain.User{Name: name, Email: email, PasswordHash: hashed}
		if err := uRepo.Create(user); err != nil {
			return err
		}
		user.Initialize()
		if err := uRepo.Update(user); err != nil {
			return err
		}
		if err := c.Outbox.Tx(tx).Save(user.EventProducer.PullEvents()); err != nil {
			return err
		}
		return nil
	})
	return err
}
