package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	"github.com/google/uuid"
)

type CreateOrderResponse struct {
	OrderID         uuid.UUID `json:"order_id"`
	ConfirmationURL string    `json:"confirmation_url"`
}

type OrderStatusResponse struct {
	ID               uuid.UUID `json:"id"`
	PackageID        uuid.UUID `json:"package_id"`
	Crystals         int       `json:"crystals"`
	AmountMinorUnits int64     `json:"amount_minor_units"`
	Currency         string    `json:"currency"`
	Provider         string    `json:"provider"`
	Status           string    `json:"status"`
	ConfirmationURL  string    `json:"confirmation_url"`
	CreatedAt        int64     `json:"created_at"`
}

func OrderStatusResponseFromReadModel(o readmodels.PurchaseOrder) OrderStatusResponse {
	return OrderStatusResponse{
		ID:               o.ID,
		PackageID:        o.PackageID,
		Crystals:         o.Crystals,
		AmountMinorUnits: o.AmountMinorUnits,
		Currency:         o.Currency,
		Provider:         o.Provider,
		Status:           o.Status,
		ConfirmationURL:  o.ConfirmationURL,
		CreatedAt:        o.CreatedAt,
	}
}
