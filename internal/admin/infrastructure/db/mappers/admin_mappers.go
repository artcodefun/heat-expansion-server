package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/admin/domain"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/db/gen"
)

// AdminFromDB maps a sqlc Admin row to a domain.Admin.
func AdminFromDB(row gen.Admin) *domain.Admin {
	return &domain.Admin{
		ID:           row.ID,
		Username:     row.Username,
		PasswordHash: nullStringToPtr(row.PasswordHash),
		InviteToken:  nullStringToPtr(row.InviteToken),
		Active:       row.Active,
		CreatedAt:    row.CreatedAt,
		UpdatedAt:    row.UpdatedAt,
	}
}

// UpdateAdminCredentialsParamsFromDomain maps a domain.Admin to sqlc update params.
func UpdateAdminCredentialsParamsFromDomain(admin *domain.Admin) gen.UpdateAdminCredentialsParams {
	return gen.UpdateAdminCredentialsParams{
		ID:           admin.ID,
		PasswordHash: ptrToNullString(admin.PasswordHash),
		InviteToken:  ptrToNullString(admin.InviteToken),
		UpdatedAt:    admin.UpdatedAt,
	}
}
