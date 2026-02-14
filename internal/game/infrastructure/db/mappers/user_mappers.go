package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

// UserFromDB maps a sqlc User model to domain.User.
func UserFromDB(u gen.User) *domain.User {
	return &domain.User{
		ID:       u.ID,
		Name:     u.Name,
		Crystals: int(u.Crystals),
	}
}
