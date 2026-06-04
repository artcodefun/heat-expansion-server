package cqrs

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/google/uuid"
)

// AdminQueries encapsulates admin read-side operations.
type AdminQueries interface {
	GetProfile(ctx context.Context, actor Actor, adminID uuid.UUID) (*readmodels.AdminProfile, error)
}
