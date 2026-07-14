package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/readstore/gen"
)

// AdminProfileFromModel maps a sqlc admin profile read row to a readmodels.AdminProfile.
func AdminProfileFromModel(row gen.GetAdminProfileRow) *readmodels.AdminProfile {
	return &readmodels.AdminProfile{
		ID:        row.ID,
		Username:  row.Username,
		Active:    row.Active,
		CreatedAt: row.CreatedAt,
	}
}
