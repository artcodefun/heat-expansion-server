package domain

import (
	"time"

	"github.com/google/uuid"
)

type AlertKind string

const (
	AlertKindCombat AlertKind = "COMBAT"
	AlertKindIntel  AlertKind = "INTEL"
	AlertKindSystem AlertKind = "SYSTEM"
)

type Alert struct {
	ID         uuid.UUID
	BaseID     int
	ActivityID *uuid.UUID
	Kind       AlertKind
	Title      TranslationKey
	Content    TranslationKey
	IsRead     bool
	CreatedAt  int64
	ExpiresAt  int64
}

func NewAlert(baseID int, activityID *uuid.UUID, kind AlertKind, title, content TranslationKey, ttl time.Duration) *Alert {
	now := time.Now().Unix()
	return &Alert{
		ID:         uuid.Must(uuid.NewV7()),
		BaseID:     baseID,
		ActivityID: activityID,
		Kind:       kind,
		Title:      title,
		Content:    content,
		IsRead:     false,
		CreatedAt:  now,
		ExpiresAt:  now + int64(ttl.Seconds()),
	}
}
