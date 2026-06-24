package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/google/uuid"
)

// GetPackageURI binds the :id parameter for package endpoints.
type GetPackageURI struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

// CrystalPackageResponse is the JSON representation of a crystal package.
type CrystalPackageResponse struct {
	ID              uuid.UUID `json:"id"`
	Name            string    `json:"name"`
	Crystals        int32     `json:"crystals"`
	PriceMinorUnits int64     `json:"price_minor_units"`
	Currency        string    `json:"currency"`
	ImageURL        string    `json:"image_url"`
	IsActive        bool      `json:"is_active"`
}

// CreateCrystalPackageRequest is the body for POST /api/v1/billing/packages.
type CreateCrystalPackageRequest struct {
	Name            string `json:"name"             binding:"required"`
	Crystals        int32  `json:"crystals"`
	PriceMinorUnits int64  `json:"price_minor_units"`
	Currency        string `json:"currency"         binding:"required"`
	ImageURL        string `json:"image_url"`
	IsActive        bool   `json:"is_active"`
}

// ToModel converts the create request to a readmodel for the command layer.
func (r CreateCrystalPackageRequest) ToModel() *readmodels.CrystalPackage {
	return &readmodels.CrystalPackage{
		Name:            r.Name,
		Crystals:        r.Crystals,
		PriceMinorUnits: r.PriceMinorUnits,
		Currency:        r.Currency,
		ImageURL:        r.ImageURL,
		IsActive:        r.IsActive,
	}
}

// UpdateCrystalPackageRequest is the body for PUT /api/v1/billing/packages/:id.
type UpdateCrystalPackageRequest struct {
	Name            string `json:"name"             binding:"required"`
	Crystals        int32  `json:"crystals"`
	PriceMinorUnits int64  `json:"price_minor_units"`
	Currency        string `json:"currency"         binding:"required"`
	ImageURL        string `json:"image_url"`
	IsActive        bool   `json:"is_active"`
}

// ToModel converts the update request to a readmodel for the command layer.
func (r UpdateCrystalPackageRequest) ToModel(id uuid.UUID) *readmodels.CrystalPackage {
	return &readmodels.CrystalPackage{
		ID:              id,
		Name:            r.Name,
		Crystals:        r.Crystals,
		PriceMinorUnits: r.PriceMinorUnits,
		Currency:        r.Currency,
		ImageURL:        r.ImageURL,
		IsActive:        r.IsActive,
	}
}

func CrystalPackageResponseFromModel(m *readmodels.CrystalPackage) CrystalPackageResponse {
	return CrystalPackageResponse{
		ID:              m.ID,
		Name:            m.Name,
		Crystals:        m.Crystals,
		PriceMinorUnits: m.PriceMinorUnits,
		Currency:        m.Currency,
		ImageURL:        m.ImageURL,
		IsActive:        m.IsActive,
	}
}
