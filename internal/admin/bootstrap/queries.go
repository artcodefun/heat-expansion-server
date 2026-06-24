package bootstrap

import (
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/queries"
)

// Queries aggregates all admin query handlers.
type Queries struct {
	Admin       *queries.AdminQueries
	Prototype   *queries.PrototypeQueries
	Translation *queries.TranslationQueries
	Package     *queries.PackageQueries
}

func NewQueries(a *Adapters) *Queries {
	return &Queries{
		Admin:       queries.NewAdminQueries(a.AdminRead),
		Prototype:   queries.NewPrototypeQueries(a.GameClient),
		Translation: queries.NewTranslationQueries(a.GameClient),
		Package:     queries.NewPackageQueries(a.BillingClient),
	}
}
