package ports

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	"github.com/google/uuid"
)

type PackageReadRepository interface {
	ListActive(ctx context.Context) ([]*readmodels.CrystalPackage, error)
	ListAll(ctx context.Context) ([]*readmodels.CrystalPackage, error)
	GetByID(ctx context.Context, id uuid.UUID) (*readmodels.CrystalPackage, error)
}

type OrderReadRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*readmodels.PurchaseOrder, error)
}
