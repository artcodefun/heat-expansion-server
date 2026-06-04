package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/admin/infrastructure/readstore/mappers"
	"github.com/google/uuid"
)

type AdminReadRepo struct {
	q *gen.Queries
}

func NewAdminReadRepo(q *gen.Queries) *AdminReadRepo {
	return &AdminReadRepo{q: q}
}

func (r *AdminReadRepo) GetProfile(ctx context.Context, adminID uuid.UUID) (*readmodels.AdminProfile, error) {
	row, err := r.q.GetAdminProfile(ctx, adminID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.AdminProfileFromModel(row), nil
}

var _ ports.AdminReadRepository = (*AdminReadRepo)(nil)
