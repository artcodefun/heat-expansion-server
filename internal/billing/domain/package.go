package domain

import "github.com/google/uuid"

// CrystalPackage is a purchasable tier of crystals.
type CrystalPackage struct {
	ID              uuid.UUID
	Name            string
	Crystals        int
	PriceMinorUnits int64  // price in minor currency units (e.g. kopecks)
	Currency        string // e.g. "RUB"
	ImageURL        string // URL of the package art displayed in the client store
	IsActive        bool
}
