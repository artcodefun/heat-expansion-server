package cqrs

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
)

// CrystalPackageCommands encapsulates privileged writes to the crystal package catalog. No Actor; pre-authorized by caller.
type CrystalPackageCommands interface {
	CreateCrystalPackage(ctx context.Context, pkg *domain.CrystalPackage) (*domain.CrystalPackage, error)
	UpdateCrystalPackage(ctx context.Context, pkg *domain.CrystalPackage) (*domain.CrystalPackage, error)
}
