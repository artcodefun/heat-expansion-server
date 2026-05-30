package readmodels

import "github.com/google/uuid"

type PurchaseOrder struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	PackageID        uuid.UUID `json:"package_id"`
	Crystals         int       `json:"crystals"`
	AmountMinorUnits int64     `json:"amount_minor_units"`
	Currency         string    `json:"currency"`
	Provider         string    `json:"provider"`
	Status           string    `json:"status"`
	ConfirmationURL  string    `json:"confirmation_url"`
	CreatedAt        int64     `json:"created_at"`
}
