package cqrs

import "github.com/google/uuid"

// CommandContext carries caller identity & auth scope for mutations.
type CommandContext struct {
	AccountID uuid.UUID
}

// AccountCommands defines the available command actions for the account aggregate.
type AccountCommands interface {
	RegisterAccount(ctx CommandContext, name, email, password string) error
	Login(ctx CommandContext, email, password string) (string, error)
}
