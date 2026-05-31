package commands

import (
	"context"

	authv1 "github.com/artcodefun/heat-expansion-server/contracts/auth/events/v1"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
)

// UserCommands maintains billing's local projection of auth accounts.
type UserCommands struct {
	UserRepo ports.UserRepository
}

func NewUserCommands(userRepo ports.UserRepository) *UserCommands {
	return &UserCommands{UserRepo: userRepo}
}

// HandleAccountRegisteredV1Event projects a newly registered auth account into
// the local billing.users table. The upsert is idempotent on the user ID, so
// redelivered events are a no-op (and a future email change would simply
// overwrite the stored address).
func (c *UserCommands) HandleAccountRegisteredV1Event(ctx context.Context, ev authv1.AccountRegisteredV1) error {
	return c.UserRepo.Upsert(ctx, &domain.User{
		ID:    ev.UserID,
		Email: ev.Email,
	})
}
