package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/google/uuid"
)

type AlertKind string

const (
	AlertKindCombat AlertKind = "COMBAT"
	AlertKindIntel  AlertKind = "INTEL"
	AlertKindSystem AlertKind = "SYSTEM"
)

type AlertItemDTO struct {
	ID         uuid.UUID  `json:"id"`
	UserID     uuid.UUID  `json:"userId"`
	BaseID     int        `json:"baseId"`
	ActivityID *uuid.UUID `json:"activityId,omitempty"`
	Kind       AlertKind  `json:"kind"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	IsRead     bool       `json:"isRead"`
	CreatedAt  int64      `json:"createdAt"`
	ExpiresAt  int64      `json:"expiresAt"`
}

func AlertItemDTOFromReadModel(a *readmodels.AlertItem, tr ports.Translator, locale string) AlertItemDTO {
	return AlertItemDTO{
		ID:         a.ID,
		UserID:     a.UserID,
		BaseID:     a.BaseID,
		ActivityID: a.ActivityID,
		Kind:       AlertKind(a.Kind),
		Title:      tr.T(locale, a.Title, nil),
		Content:    tr.T(locale, a.Content, nil),
		IsRead:     a.IsRead,
		CreatedAt:  a.CreatedAt,
		ExpiresAt:  a.ExpiresAt,
	}
}

func AlertItemsFromReadModels(items []*readmodels.AlertItem, tr ports.Translator, locale string) []AlertItemDTO {
	out := make([]AlertItemDTO, 0, len(items))
	for _, item := range items {
		out = append(out, AlertItemDTOFromReadModel(item, tr, locale))
	}
	return out
}
