package cqrs

import (
	"context"
)

// AccountCommands defines the available command actions for the account aggregate.
type AccountCommands interface {
	RegisterAccount(ctx context.Context, actor Actor, name, email, password string) error
	Login(ctx context.Context, actor Actor, email, password string) (string, error)
	RequestPasswordReset(ctx context.Context, actor Actor, email string) error
	ResetPassword(ctx context.Context, actor Actor, email, rawToken, newPassword string) error
}
