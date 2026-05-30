package cqrs

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	"github.com/google/uuid"
)

type PackageQueries interface {
	ListPackages(ctx context.Context) ([]readmodels.CrystalPackage, error)
}

type OrderQueries interface {
	GetOrder(ctx context.Context, actor Actor, orderID uuid.UUID) (*readmodels.PurchaseOrder, error)
}
