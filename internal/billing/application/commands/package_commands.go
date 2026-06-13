package commands

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
)

// CrystalPackageCommands handles privileged writes to the crystal package catalog.
type CrystalPackageCommands struct {
	Packages ports.CrystalPackageRepository
}

func NewCrystalPackageCommands(packages ports.CrystalPackageRepository) *CrystalPackageCommands {
	return &CrystalPackageCommands{Packages: packages}
}

func (c *CrystalPackageCommands) CreateCrystalPackage(ctx context.Context, pkg *domain.CrystalPackage) (*domain.CrystalPackage, error) {
	return c.Packages.Create(ctx, pkg)
}

func (c *CrystalPackageCommands) UpdateCrystalPackage(ctx context.Context, pkg *domain.CrystalPackage) (*domain.CrystalPackage, error) {
	return c.Packages.Update(ctx, pkg)
}
