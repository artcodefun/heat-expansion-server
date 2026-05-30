package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/gen"
)

func OrderFromRow(row gen.PurchaseOrder) *domain.PurchaseOrder {
	return &domain.PurchaseOrder{
		ID:               row.ID,
		UserID:           row.UserID,
		PackageID:        row.PackageID,
		Crystals:         int(row.Crystals),
		AmountMinorUnits: row.AmountMinorUnits,
		Currency:         row.Currency,
		Provider:         domain.PaymentProvider(row.Provider),
		Status:           domain.OrderStatus(row.Status),
		ProviderOrderID:  row.ProviderOrderID,
		ConfirmationURL:  row.ConfirmationUrl,
		CreatedAt:        row.CreatedAt,
		UpdatedAt:        row.UpdatedAt,
	}
}

func InsertOrderParamsFromDomain(order *domain.PurchaseOrder) gen.InsertOrderParams {
	return gen.InsertOrderParams{
		ID:               order.ID,
		UserID:           order.UserID,
		PackageID:        order.PackageID,
		Crystals:         int32(order.Crystals),
		AmountMinorUnits: order.AmountMinorUnits,
		Currency:         order.Currency,
		Provider:         string(order.Provider),
		Status:           string(order.Status),
		ProviderOrderID:  order.ProviderOrderID,
		ConfirmationUrl:  order.ConfirmationURL,
		CreatedAt:        order.CreatedAt,
		UpdatedAt:        order.UpdatedAt,
	}
}

func UpdateOrderParamsFromDomain(order *domain.PurchaseOrder) gen.UpdateOrderParams {
	return gen.UpdateOrderParams{
		ID:              order.ID,
		Status:          string(order.Status),
		ProviderOrderID: order.ProviderOrderID,
		ConfirmationUrl: order.ConfirmationURL,
		UpdatedAt:       order.UpdatedAt,
	}
}
