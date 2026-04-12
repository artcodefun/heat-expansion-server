package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
)

// InsertAlertParamsFromDomain maps a domain.Alert into sqlc params for insert.
func InsertAlertParamsFromDomain(alert *domain.Alert) gen.InsertAlertParams {
	params := gen.InsertAlertParams{
		ID:         alert.ID,
		UserID:     alert.UserID,
		Kind:       string(alert.Kind),
		Title:      alert.Title,
		Content:    alert.Content,
		IsRead:     alert.IsRead,
		CreatedAt:  alert.CreatedAt,
		ExpiresAt:  alert.ExpiresAt,
		BaseID:     nullableBaseID(alert.BaseID),
		ActivityID: nullableUUID(alert.ActivityID),
	}
	return params
}
