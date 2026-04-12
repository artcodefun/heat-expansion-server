package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
)

// AlertItemFromModel maps an infrastructure alert record to a read model.
func AlertItemFromModel(row gen.Alert) *readmodels.AlertItem {
	alert := &readmodels.AlertItem{
		ID:         row.ID,
		UserID:     row.UserID,
		Kind:       readmodels.AlertKind(row.Kind),
		Title:      row.Title,
		Content:    row.Content,
		IsRead:     row.IsRead,
		CreatedAt:  row.CreatedAt,
		ExpiresAt:  row.ExpiresAt,
		BaseID:     nullBaseIDPtr(row.BaseID),
		ActivityID: nullUUIDPtr(row.ActivityID),
	}
	return alert
}

// AlertItemsFromModels maps slice of infrastructure alert records to read models.
func AlertItemsFromModels(rows []gen.Alert) []*readmodels.AlertItem {
	alerts := make([]*readmodels.AlertItem, 0, len(rows))
	for _, row := range rows {
		alerts = append(alerts, AlertItemFromModel(row))
	}
	return alerts
}
