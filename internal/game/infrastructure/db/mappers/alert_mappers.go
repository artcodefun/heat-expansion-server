package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/gen"
	"github.com/google/uuid"
)

// InsertAlertParamsFromDomain maps a domain.Alert into sqlc params for insert.
func InsertAlertParamsFromDomain(alert *domain.Alert) gen.InsertAlertParams {
	params := gen.InsertAlertParams{
		ID:        alert.ID,
		BaseID:    int64(alert.BaseID),
		Kind:      string(alert.Kind),
		Title:     alert.Title,
		Content:   alert.Content,
		IsRead:    alert.IsRead,
		CreatedAt: alert.CreatedAt,
		ExpiresAt: alert.ExpiresAt,
	}
	if alert.ActivityID != nil {
		params.ActivityID = uuid.NullUUID{UUID: *alert.ActivityID, Valid: true}
	}
	return params
}
