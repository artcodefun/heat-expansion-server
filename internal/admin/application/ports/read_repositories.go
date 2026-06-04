package ports

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/google/uuid"
)

// AdminReadRepository provides read-only projections of admin data.
type AdminReadRepository interface {
	GetProfile(ctx context.Context, adminID uuid.UUID) (*readmodels.AdminProfile, error)
}
