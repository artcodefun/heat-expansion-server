package bootstrap

import "github.com/artcodefun/heat-expansion-server/internal/billing/application/commands"

// Commands aggregates all command handlers.
type Commands struct {
	Order *commands.OrderCommands
	User  *commands.UserCommands
}

// NewCommands constructs all command handlers using the provided secondary adapters.
func NewCommands(a *Adapters) *Commands {
	return &Commands{
		Order: commands.NewOrderCommands(a.Orders, a.Packages, a.Users, a.Gateway, a.Outbox, a.TxMgr),
		User:  commands.NewUserCommands(a.Users),
	}
}
