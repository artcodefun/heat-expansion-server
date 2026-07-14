package ports

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/google/uuid"
)

// BillingPrivateClient abstracts outbound calls to the billing module's private gRPC API.
type BillingPrivateClient interface {
	ListCrystalPackages(ctx context.Context) ([]*readmodels.CrystalPackage, error)
	GetCrystalPackage(ctx context.Context, id uuid.UUID) (*readmodels.CrystalPackage, error)
	CreateCrystalPackage(ctx context.Context, p *readmodels.CrystalPackage) (*readmodels.CrystalPackage, error)
	UpdateCrystalPackage(ctx context.Context, p *readmodels.CrystalPackage) (*readmodels.CrystalPackage, error)
}
