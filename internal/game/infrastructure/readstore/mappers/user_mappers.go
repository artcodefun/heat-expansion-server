package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
)

// UserFromModel maps a sqlc User model into a readmodels.User
func UserFromModel(u gen.User) readmodels.User {
	return readmodels.User{
		ID:       u.ID,
		Name:     u.Name,
		Crystals: int(u.Crystals),
	}
}
