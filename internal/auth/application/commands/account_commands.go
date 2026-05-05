package commands

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

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
	resetRepo     ports.PasswordResetRepository
	emailSender   ports.EmailSender
}

func NewAccountCommands(
	repo ports.AccountRepository,
	hasher ports.PasswordHasher,
	tokenProvider ports.TokenProvider,
	outbox ports.OutboxEventRepository,
	txMgr ports.TransactionManager,
	resetRepo ports.PasswordResetRepository,
	emailSender ports.EmailSender,
) *AccountCommands {
	return &AccountCommands{
		repo:          repo,
		hasher:        hasher,
		tokenProvider: tokenProvider,
		outbox:        outbox,
		txMgr:         txMgr,
		resetRepo:     resetRepo,
		emailSender:   emailSender,
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

func (c *AccountCommands) RequestPasswordReset(ctx context.Context, actor cqrs.Actor, email string) error {
	_ = actor

	acc, err := c.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if acc == nil {
		return cqrs.ErrAccountNotFound
	}

	now := time.Now().Unix()
	if err := c.resetRepo.InvalidateByAccount(ctx, acc.ID, now); err != nil {
		return err
	}

	rawToken, tokenHash, err := generateResetToken()
	if err != nil {
		return err
	}

	resetToken := domain.NewPasswordResetToken(acc.ID, tokenHash)
	if err := c.resetRepo.Create(ctx, resetToken); err != nil {
		return err
	}

	if err := c.emailSender.SendPasswordReset(ctx, acc.Email, acc.Name, rawToken); err != nil {
		slog.ErrorContext(ctx, "failed to send password reset email", "account_id", acc.ID.String(), "error", err)
		return err
	}

	return nil
}

func (c *AccountCommands) ResetPassword(ctx context.Context, actor cqrs.Actor, email, rawToken, newPassword string) error {
	_ = actor

	acc, err := c.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if acc == nil {
		return cqrs.ErrAccountNotFound
	}

	tokenHash := hashToken(rawToken)

	resetToken, err := c.resetRepo.FindByAccountAndTokenHash(ctx, acc.ID, tokenHash)
	if err != nil {
		return err
	}
	if resetToken == nil || resetToken.IsExpired() || resetToken.IsUsed() {
		return cqrs.ErrInvalidResetToken
	}

	newHash, err := c.hasher.Hash(newPassword)
	if err != nil {
		return err
	}
	acc.ChangePassword(newHash)

	now := time.Now().Unix()
	return c.txMgr.WithTx(ctx, func(tx ports.Transaction) error {
		if err := c.repo.Tx(tx).UpdatePassword(ctx, acc.ID, acc.PasswordHash); err != nil {
			return err
		}
		return c.resetRepo.Tx(tx).MarkUsed(ctx, resetToken.ID, now)
	})
}

// generateResetToken returns an 8-digit numeric code and its SHA-256 hash.
// Short enough to type in a game client; 1-hour TTL limits brute-force exposure.
func generateResetToken() (rawToken, tokenHash string, err error) {
	b := make([]byte, 4)
	if _, err = rand.Read(b); err != nil {
		return
	}
	n := binary.BigEndian.Uint32(b) % 100_000_000
	rawToken = fmt.Sprintf("%08d", n)
	tokenHash = hashToken(rawToken)
	return
}

func hashToken(rawToken string) string {
	h := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(h[:])
}
