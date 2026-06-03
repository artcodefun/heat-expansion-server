package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	"github.com/google/uuid"
)

type CrystalPackageResponse struct {
	ID              uuid.UUID `json:"id"`
	Crystals        int       `json:"crystals"`
	PriceMinorUnits int64     `json:"price_minor_units"`
	Currency        string    `json:"currency"`
	ImageURL        string    `json:"image_url"`
}

func CrystalPackageResponseFromReadModel(p *readmodels.CrystalPackage) CrystalPackageResponse {
	return CrystalPackageResponse{
		ID:              p.ID,
		Crystals:        p.Crystals,
		PriceMinorUnits: p.PriceMinorUnits,
		Currency:        p.Currency,
		ImageURL:        p.ImageURL,
	}
}

func CrystalPackageResponsesFromReadModels(pkgs []*readmodels.CrystalPackage) []CrystalPackageResponse {
	out := make([]CrystalPackageResponse, len(pkgs))
	for i, p := range pkgs {
		out[i] = CrystalPackageResponseFromReadModel(p)
	}
	return out
}
