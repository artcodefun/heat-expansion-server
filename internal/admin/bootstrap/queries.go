package bootstrap

import "github.com/artcodefun/heat-expansion-server/internal/admin/application/queries"

// Queries aggregates all admin query handlers.
type Queries struct {
	Admin *queries.AdminQueries
}

func NewQueries(a *Adapters) *Queries {
	return &Queries{
		Admin: queries.NewAdminQueries(a.AdminRead),
	}
}
