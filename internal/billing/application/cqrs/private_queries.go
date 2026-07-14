package cqrs

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	"github.com/google/uuid"
)

// CrystalPackageQueries exposes full crystal package catalog reads. No Actor; pre-authorized by caller.
type CrystalPackageQueries interface {
	ListAllCrystalPackages(ctx context.Context) ([]*readmodels.CrystalPackage, error)
	GetCrystalPackage(ctx context.Context, id uuid.UUID) (*readmodels.CrystalPackage, error)
}
