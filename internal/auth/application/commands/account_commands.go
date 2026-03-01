package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
)

type AccountCommands struct {
	repo          ports.AccountRepository
	hasher        ports.PasswordHasher
	tokenProvider ports.TokenProvider
	outbox        ports.OutboxEventRepository
	txMgr         ports.TransactionManager
}

func NewAccountCommands(
	repo ports.AccountRepository,
	hasher ports.PasswordHasher,
	tokenProvider ports.TokenProvider,
	outbox ports.OutboxEventRepository,
	txMgr ports.TransactionManager,
) *AccountCommands {
	return &AccountCommands{
		repo:          repo,
		hasher:        hasher,
		tokenProvider: tokenProvider,
		outbox:        outbox,
		txMgr:         txMgr,
	}
}

func (c *AccountCommands) RegisterAccount(ctx context.Context, actor cqrs.Actor, name, email, password string) error {
	_ = actor

	hash, err := c.hasher.Hash(password)
	if err != nil {
		return err
	}

	acc := domain.RegisterAccount(name, email, hash)

	return c.txMgr.WithTx(ctx, func(tx ports.Transaction) error {
		repo := c.repo.Tx(tx)
		outbox := c.outbox.Tx(tx)

		// Check if email already exists
		existing, err := repo.FindByEmail(ctx, email)
		if err != nil {
			return err
		}
		if existing != nil {
			return cqrs.ErrEmailAlreadyInUse
		}

		if err := repo.Create(ctx, acc); err != nil {
			return err
		}

		return outbox.Save(ctx, acc.PullEvents())
	})
}

func (c *AccountCommands) Login(ctx context.Context, actor cqrs.Actor, email, password string) (string, error) {
	_ = actor

	acc, err := c.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if acc == nil {
		return "", cqrs.ErrInvalidCredentials
	}

	if !c.hasher.Verify(password, acc.PasswordHash) {
		return "", cqrs.ErrInvalidCredentials
	}

	return c.tokenProvider.Generate(acc.ID)
}
