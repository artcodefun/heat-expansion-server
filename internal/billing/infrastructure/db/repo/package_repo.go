package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/gen"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/mappers"
	"github.com/google/uuid"
)

type PackageRepo struct {
	q *gen.Queries
}

func NewPackageRepo(q *gen.Queries) *PackageRepo {
	return &PackageRepo{q: q}
}

func (r *PackageRepo) FindByID(ctx context.Context, id uuid.UUID) (*domain.CrystalPackage, error) {
	row, err := r.q.GetPackageByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.PackageFromRow(row), nil
}
