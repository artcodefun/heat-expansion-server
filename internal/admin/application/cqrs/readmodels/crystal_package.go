package readmodels

import "github.com/google/uuid"

// CrystalPackage is the admin read model for a billing crystal package.
type CrystalPackage struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Crystals        int32     `json:"crystals"`
	PriceMinorUnits int64     `json:"price_minor_units"`
	Currency        string    `json:"currency"`
	ImageURL        string    `json:"image_url"`
	IsActive        bool      `json:"is_active"`
}
