package bootstrap

import "github.com/artcodefun/heat-expansion-server/internal/auth/application/queries"

type Queries struct {
	Account *queries.AccountQueries
}

func NewQueries(adapters *Adapters) *Queries {
	return &Queries{
		Account: queries.NewAccountQueries(),
	}
}
