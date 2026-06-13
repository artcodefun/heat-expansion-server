package readmodels

import "github.com/google/uuid"

type CrystalPackage struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Crystals        int       `json:"crystals"`
	PriceMinorUnits int64     `json:"price_minor_units"`
	Currency        string    `json:"currency"`
	ImageURL        string    `json:"image_url"`
	IsActive        bool      `json:"is_active"`
}
