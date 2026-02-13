package bootstrap

import (
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/cqrs"
	application "github.com/artcodefun/heat-expansion-server/internal/auth/application/queries"
)

type Queries struct {
	Account cqrs.AccountQueries
}

func NewQueries(adapters *Adapters) *Queries {
	return &Queries{
		Account: application.NewAccountQueries(),
	}
}
