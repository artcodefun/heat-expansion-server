package bootstrap

import (
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs"
	appqueries "github.com/artcodefun/heat-expansion-server/internal/billing/application/queries"
)

// Queries aggregates all query handlers.
type Queries struct {
	Package cqrs.PackageQueries
	Order   cqrs.OrderQueries
}

// NewQueries constructs all query handlers using the provided secondary adapters.
func NewQueries(a *Adapters) *Queries {
	return &Queries{
		Package: appqueries.NewPackageQueries(a.PackageRead),
		Order:   appqueries.NewOrderQueries(a.OrderRead),
	}
}
