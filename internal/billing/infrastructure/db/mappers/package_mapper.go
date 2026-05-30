package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	"github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/db/gen"
)

func PackageFromRow(row gen.CrystalPackage) *domain.CrystalPackage {
	return &domain.CrystalPackage{
		ID:              row.ID,
		Name:            row.Name,
		Crystals:        int(row.Crystals),
		PriceMinorUnits: row.PriceMinorUnits,
		Currency:        row.Currency,
		ImageURL:        row.ImageUrl,
		IsActive:        row.IsActive,
	}
}
