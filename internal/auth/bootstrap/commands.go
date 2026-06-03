package bootstrap

import "github.com/artcodefun/heat-expansion-server/internal/auth/application/commands"

type Commands struct {
	Account *commands.AccountCommands
}

func NewCommands(adapters *Adapters) *Commands {
	return &Commands{
		Account: commands.NewAccountCommands(
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
