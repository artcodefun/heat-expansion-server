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

func (r *PackageRepo) Create(ctx context.Context, pkg *domain.CrystalPackage) (*domain.CrystalPackage, error) {
	row, err := r.q.CreatePackage(ctx, gen.CreatePackageParams{
		ID:              pkg.ID,
		Name:            pkg.Name,
		Crystals:        int32(pkg.Crystals),
		PriceMinorUnits: pkg.PriceMinorUnits,
		Currency:        pkg.Currency,
		ImageUrl:        pkg.ImageURL,
		IsActive:        pkg.IsActive,
		CreatedAt:       domain.NowUnix(),
		UpdatedAt:       domain.NowUnix(),
	})
	if err != nil {
		return nil, err
	}
	return mappers.PackageFromRow(row), nil
}

func (r *PackageRepo) Update(ctx context.Context, pkg *domain.CrystalPackage) (*domain.CrystalPackage, error) {
	row, err := r.q.UpdatePackage(ctx, gen.UpdatePackageParams{
		ID:              pkg.ID,
		Name:            pkg.Name,
		Crystals:        int32(pkg.Crystals),
		PriceMinorUnits: pkg.PriceMinorUnits,
		Currency:        pkg.Currency,
		ImageUrl:        pkg.ImageURL,
		IsActive:        pkg.IsActive,
		UpdatedAt:       domain.NowUnix(),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ports.ErrNotFound
		}
		return nil, err
	}
	return mappers.PackageFromRow(row), nil
}
