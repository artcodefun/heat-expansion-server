package bootstrap

import (
	application "github.com/artcodefun/heat-expansion-server/internal/auth/application/commands"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/cqrs"
)

type Commands struct {
	Account cqrs.AccountCommands
}

func NewCommands(adapters *Adapters) *Commands {
	return &Commands{
		Account: application.NewAccountCommands(
			adapters.Repo,
			adapters.Hasher,
			adapters.TokenProvider,
			adapters.Outbox,
			adapters.TxMgr,
			adapters.ResetRepo,
			adapters.EmailSender,
		),
	}
}
