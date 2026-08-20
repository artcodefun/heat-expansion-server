package readmodels

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

type AlertKind string

const (
	AlertKindCombat     AlertKind = "COMBAT"
	AlertKindIntel      AlertKind = "INTEL"
	AlertKindSystem     AlertKind = "SYSTEM"
	AlertKindDiplomatic AlertKind = "DIPLOMATIC"
	AlertKindTrade      AlertKind = "TRADE"
)

type AlertItem struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	BaseID     *int
	ActivityID *uuid.UUID
	Kind       AlertKind
	Title      domain.TranslationKey
	Content    domain.TranslationKey
	IsRead     bool
	CreatedAt  int64
	ExpiresAt  int64
}
