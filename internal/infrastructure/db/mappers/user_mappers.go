package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
)

// UserFromDB maps a sqlc User model to domain.User.
func UserFromDB(u gen.User) *domain.User {
	return &domain.User{
		ID:           int(u.ID),
		Name:         u.Name,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Crystals:     int(u.Crystals),
		// BaseSectorIDs populated elsewhere (aggregate loading); left empty here.
	}
}
