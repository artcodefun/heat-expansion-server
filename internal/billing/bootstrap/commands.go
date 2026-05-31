package bootstrap

import (
	appcommands "github.com/artcodefun/heat-expansion-server/internal/billing/application/commands"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs"
)

// Commands aggregates all command handlers.
type Commands struct {
	Order cqrs.OrderCommands
	User  cqrs.UserCommands
}

// NewCommands constructs all command handlers using the provided secondary adapters.
func NewCommands(a *Adapters) *Commands {
	return &Commands{
		Order: appcommands.NewOrderCommands(a.Orders, a.Packages, a.Users, a.Gateway, a.Outbox, a.TxMgr),
		User:  appcommands.NewUserCommands(a.Users),
	}
}
