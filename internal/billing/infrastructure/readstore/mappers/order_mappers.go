package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs/readmodels"
	dbgen "github.com/artcodefun/heat-expansion-server/internal/billing/infrastructure/readstore/gen"
)

func OrderReadModelFromRow(row dbgen.GetOrderByIDRow) readmodels.PurchaseOrder {
	return readmodels.PurchaseOrder{
		ID:               row.ID,
		UserID:           row.UserID,
		PackageID:        row.PackageID,
		Crystals:         int(row.Crystals),
		AmountMinorUnits: row.AmountMinorUnits,
		Currency:         row.Currency,
		Provider:         row.Provider,
		Status:           row.Status,
		ConfirmationURL:  row.ConfirmationUrl,
		CreatedAt:        row.CreatedAt,
	}
}
