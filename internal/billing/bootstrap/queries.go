package bootstrap

import "github.com/artcodefun/heat-expansion-server/internal/billing/application/queries"

// Queries aggregates all query handlers.
type Queries struct {
	Package *queries.PackageQueries
	Order   *queries.OrderQueries
}

// NewQueries constructs all query handlers using the provided secondary adapters.
func NewQueries(a *Adapters) *Queries {
	return &Queries{
		Package: queries.NewPackageQueries(a.PackageRead),
		Order:   queries.NewOrderQueries(a.OrderRead),
	}
}
