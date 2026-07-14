package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/mappers"
	"github.com/google/uuid"
)

type PackageReadRepo struct {
	q *dbgen.Queries
}

func NewPackageReadRepo(q *dbgen.Queries) *PackageReadRepo {
	return &PackageReadRepo{q: q}
}

func (r *PackageReadRepo) ListActive(ctx context.Context) ([]*readmodels.CrystalPackage, error) {
	rows, err := r.q.ListActivePackages(ctx)
	if err != nil {
		return nil, err
	}
	return mappers.PackageReadModelsFromRows(rows), nil
}

func (r *PackageReadRepo) ListAll(ctx context.Context) ([]*readmodels.CrystalPackage, error) {
	rows, err := r.q.ListAllPackages(ctx)
	if err != nil {
		return nil, err
	}
	return mappers.PackageAllReadModelsFromRows(rows), nil
}

func (r *PackageReadRepo) GetByID(ctx context.Context, id uuid.UUID) (*readmodels.CrystalPackage, error) {
	row, err := r.q.GetPackageByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.PackageFromGetRow(row), nil
}
