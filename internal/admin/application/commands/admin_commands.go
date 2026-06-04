package commands

import (
	"context"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/domain"
)

// AdminCommands implements cqrs.AdminCommands.
type AdminCommands struct {
	admins   ports.AdminRepository
	sessions ports.SessionRepository
	hasher   ports.PasswordHasher
	tokenGen ports.SessionTokenGenerator
	txMgr    ports.TransactionManager
}

func NewAdminCommands(
	admins ports.AdminRepository,
	sessions ports.SessionRepository,
	hasher ports.PasswordHasher,
	tokenGen ports.SessionTokenGenerator,
	txMgr ports.TransactionManager,
) *AdminCommands {
	return &AdminCommands{
		admins:   admins,
		sessions: sessions,
		hasher:   hasher,
		tokenGen: tokenGen,
		txMgr:    txMgr,
	}
}

// Register completes first-time setup for an unregistered admin and issues a session.
func (c *AdminCommands) Register(ctx context.Context, actor cqrs.Actor, username, inviteToken, newPassword string) (string, error) {
	_ = actor

	admin, err := c.admins.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			// Do not reveal whether the username exists.
			return "", cqrs.NewAppError(cqrs.KindForbidden, "error.application.admin.invalid_invite_token")
		}
		return "", err
	}

	passwordHash, err := c.hasher.Hash(newPassword)
	if err != nil {
		return "", err
	}
	if err := admin.Register(inviteToken, passwordHash); err != nil {
		return "", err
	}

	sessionToken, err := c.tokenGen.Generate()
	if err != nil {
		return "", err
	}
	session := domain.NewSession(admin.ID, sessionToken)

	err = c.txMgr.WithTx(ctx, func(tx ports.Transaction) error {
		if err := c.admins.Tx(tx).Save(ctx, admin); err != nil {
			return err
		}
		if err := c.sessions.Tx(tx).Create(ctx, session); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	return session.Token, nil
}

// Login verifies admin credentials and issues a new session.
func (c *AdminCommands) Login(ctx context.Context, actor cqrs.Actor, username, password string) (string, error) {
	_ = actor

	admin, err := c.admins.FindByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return "", cqrs.NewAppError(cqrs.KindInvalidInput, "error.application.admin.invalid_credentials")
		}
		return "", err
	}

	if err := admin.CanAuthenticate(); err != nil {
		return "", err
	}

	if !c.hasher.Verify(password, *admin.PasswordHash) {
		return "", cqrs.NewAppError(cqrs.KindInvalidInput, "error.application.admin.invalid_credentials")
	}

	sessionToken, err := c.tokenGen.Generate()
	if err != nil {
		return "", err
	}
	session := domain.NewSession(admin.ID, sessionToken)

	if err := c.sessions.Create(ctx, session); err != nil {
		return "", err
	}
	return session.Token, nil
}

// Logout revokes the session identified by the bearer token.
func (c *AdminCommands) Logout(ctx context.Context, actor cqrs.Actor, token string) error {
	_ = actor

	if err := c.sessions.Delete(ctx, token); err != nil && !errors.Is(err, ports.ErrNotFound) {
		return err
	}
	return nil
}

var _ cqrs.AdminCommands = (*AdminCommands)(nil)
