package readmodels

import "github.com/google/uuid"

type AlertKind string

const (
	AlertKindCombat AlertKind = "COMBAT"
	AlertKindIntel  AlertKind = "INTEL"
	AlertKindSystem AlertKind = "SYSTEM"
)

type AlertItem struct {
	ID         uuid.UUID
	BaseID     int
	ActivityID *uuid.UUID
	Kind       AlertKind
	Title      string
	Content    string
	IsRead     bool
	CreatedAt  int64
	ExpiresAt  int64
}
