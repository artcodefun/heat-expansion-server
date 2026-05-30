package ports

import (
	"context"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

type PurchaseOrderRepository interface {
	Save(ctx context.Context, order *domain.PurchaseOrder) error
	Update(ctx context.Context, order *domain.PurchaseOrder) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.PurchaseOrder, error)
	FindByProviderOrderID(ctx context.Context, providerOrderID string) (*domain.PurchaseOrder, error)
	FindByProviderOrderIDForUpdate(ctx context.Context, providerOrderID string) (*domain.PurchaseOrder, error)
	Tx(tx Transaction) PurchaseOrderRepository
}

type CrystalPackageRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (*domain.CrystalPackage, error)
}
