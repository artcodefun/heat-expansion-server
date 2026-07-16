package ports

import (
	"context"

	"github.com/google/uuid"
)

// BillingPrivateClient abstracts outbound calls to the billing module's private gRPC API.
type BillingPrivateClient interface {
	ListCrystalPackages(ctx context.Context) ([]*CrystalPackage, error)
	GetCrystalPackage(ctx context.Context, id uuid.UUID) (*CrystalPackage, error)
	CreateCrystalPackage(ctx context.Context, p *CrystalPackage) (*CrystalPackage, error)
	UpdateCrystalPackage(ctx context.Context, p *CrystalPackage) (*CrystalPackage, error)
}
