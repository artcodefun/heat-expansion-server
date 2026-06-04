package cqrs

import "context"

// AdminCommands encapsulates admin authentication mutations.
type AdminCommands interface {
	// Register completes first-time setup using the invite token, sets a password,
	// and issues a session immediately. Returns the raw session token.
	Register(ctx context.Context, actor Actor, username, inviteToken, newPassword string) (string, error)
	// Login verifies credentials and issues a new session. Returns the raw session token.
	Login(ctx context.Context, actor Actor, username, password string) (string, error)
	// Logout revokes the session identified by the given bearer token.
	Logout(ctx context.Context, actor Actor, token string) error
}
