package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/gen"
)

func PackageReadModelFromRow(row dbgen.ListActivePackagesRow) readmodels.CrystalPackage {
	return readmodels.CrystalPackage{
		ID:              row.ID,
		Crystals:        int(row.Crystals),
		PriceMinorUnits: row.PriceMinorUnits,
		Currency:        row.Currency,
		ImageURL:        row.ImageUrl,
	}
}

func PackageReadModelsFromRows(rows []dbgen.ListActivePackagesRow) []readmodels.CrystalPackage {
	out := make([]readmodels.CrystalPackage, len(rows))
	for i, row := range rows {
		out[i] = PackageReadModelFromRow(row)
	}
	return out
}
