package commands

import (
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserCommands struct {
	UserRepo       ports.UserRepository
	PasswordHasher ports.PasswordHasher
	TokenProvider  ports.TokenProvider
	EventPublisher ports.EventPublisher
	TxMgr          ports.TransactionManager
}

func NewUserCommands(userRepo ports.UserRepository, hasher ports.PasswordHasher, tokenProvider ports.TokenProvider, eventPublisher ports.EventPublisher, txMgr ports.TransactionManager) *UserCommands {
	return &UserCommands{UserRepo: userRepo, PasswordHasher: hasher, TokenProvider: tokenProvider, EventPublisher: eventPublisher, TxMgr: txMgr}
}

func (c *UserCommands) Authenticate(ctx cqrs.CommandContext, email, password string) (string, error) {
	user, err := c.UserRepo.FindByEmail(email)
	if err != nil {
		return "", repoErr(err)
	}
	if !c.PasswordHasher.Verify(password, user.PasswordHash) {
		return "", ErrInvalidCredentials
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
	var events []domain.DomainEvent
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
		events = append(events, user.EventProducer.PullEvents()...)
		return nil
	})
	if err != nil {
		return err
	}
	publishEvents(events, c.EventPublisher)
	return nil
}
