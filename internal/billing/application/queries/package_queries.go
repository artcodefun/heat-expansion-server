package queries

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/google/uuid"
)

type PackageQueries struct {
	Packages ports.PackageReadRepository
}

func NewPackageQueries(packages ports.PackageReadRepository) *PackageQueries {
	return &PackageQueries{Packages: packages}
}

func (q *PackageQueries) ListPackages(ctx context.Context) ([]*readmodels.CrystalPackage, error) {
	return q.Packages.ListActive(ctx)
}

func (q *PackageQueries) ListAllCrystalPackages(ctx context.Context) ([]*readmodels.CrystalPackage, error) {
	return q.Packages.ListAll(ctx)
}

func (q *PackageQueries) GetCrystalPackage(ctx context.Context, id uuid.UUID) (*readmodels.CrystalPackage, error) {
	return q.Packages.GetByID(ctx, id)
}
