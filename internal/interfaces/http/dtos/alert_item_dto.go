package dtos

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
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
	BaseID     int        `json:"baseId"`
	ActivityID *uuid.UUID `json:"activityId,omitempty"`
	Kind       AlertKind  `json:"kind"`
	Title      string     `json:"title"`
	Content    string     `json:"content"`
	IsRead     bool       `json:"isRead"`
	CreatedAt  int64      `json:"createdAt"`
	ExpiresAt  int64      `json:"expiresAt"`
}

func AlertItemDTOFromReadModel(a *readmodels.AlertItem) AlertItemDTO {
	return AlertItemDTO{
		ID:         a.ID,
		BaseID:     a.BaseID,
		ActivityID: a.ActivityID,
		Kind:       AlertKind(a.Kind),
		Title:      a.Title,
		Content:    a.Content,
		IsRead:     a.IsRead,
		CreatedAt:  a.CreatedAt,
		ExpiresAt:  a.ExpiresAt,
	}
}

func AlertItemsFromReadModels(items []*readmodels.AlertItem) []AlertItemDTO {
	out := make([]AlertItemDTO, 0, len(items))
	for _, item := range items {
		out = append(out, AlertItemDTOFromReadModel(item))
	}
	return out
}
