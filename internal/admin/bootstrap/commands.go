package bootstrap

import (
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/commands"
)

// Commands aggregates all admin command handlers.
type Commands struct {
	Admin       *commands.AdminCommands
	Prototype   *commands.PrototypeCommands
	Translation *commands.TranslationCommands
	Package     *commands.PackageCommands
}

func NewCommands(a *Adapters) *Commands {
	return &Commands{
		Admin:       commands.NewAdminCommands(a.Admins, a.Sessions, a.Hasher, a.TokenGen, a.TxMgr),
		Prototype:   commands.NewPrototypeCommands(a.GameClient),
		Translation: commands.NewTranslationCommands(a.GameClient),
		Package:     commands.NewPackageCommands(a.BillingClient),
	}
}
