package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/readmodels"
	"github.com/google/uuid"
)

// AdminQueries implements application.AdminQueries.
type AdminQueries struct {
	admins ports.AdminReadRepository
}

func NewAdminQueries(admins ports.AdminReadRepository) *AdminQueries {
	return &AdminQueries{admins: admins}
}

func (q *AdminQueries) GetProfile(ctx context.Context, actor application.Actor, adminID uuid.UUID) (*readmodels.AdminProfile, error) {
	_ = actor
	profile, err := q.admins.GetProfile(ctx, adminID)
	return profile, repoErr(err)
}

var _ application.AdminQueries = (*AdminQueries)(nil)
