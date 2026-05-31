package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/gen"
)

func UserFromRow(row gen.User) *domain.User {
	return &domain.User{
		ID:    row.ID,
		Email: row.Email,
	}
}
