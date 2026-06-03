package repo

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/gen"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/mappers"
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
