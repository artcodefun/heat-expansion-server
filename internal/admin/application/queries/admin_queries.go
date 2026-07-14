package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/google/uuid"
)

// AdminQueries implements cqrs.AdminQueries.
type AdminQueries struct {
	admins ports.AdminReadRepository
}

func NewAdminQueries(admins ports.AdminReadRepository) *AdminQueries {
	return &AdminQueries{admins: admins}
}

func (q *AdminQueries) GetProfile(ctx context.Context, actor cqrs.Actor, adminID uuid.UUID) (*readmodels.AdminProfile, error) {
	_ = actor
	profile, err := q.admins.GetProfile(ctx, adminID)
	return profile, repoErr(err)
}

var _ cqrs.AdminQueries = (*AdminQueries)(nil)
