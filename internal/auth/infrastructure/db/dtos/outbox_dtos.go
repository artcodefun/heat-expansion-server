package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/auth/domain"
	"github.com/google/uuid"
)

// =========================
// Domain Event DTOs
// =========================

type AccountRegisteredEventDTO struct {
	ID         uuid.UUID `json:"id"`
	OccurredAt int64     `json:"occurred_at"`
	AccountID  uuid.UUID `json:"account_id"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
}

func AccountRegisteredEventDTOFromDomain(e domain.AccountRegisteredEvent) AccountRegisteredEventDTO {
	return AccountRegisteredEventDTO{
		ID:         e.ID(),
		OccurredAt: e.OccurredAt(),
		AccountID:  e.AccountID,
		Name:       e.Name,
		Email:      e.Email,
	}
}

func AccountRegisteredEventFromDTO(d AccountRegisteredEventDTO) domain.AccountRegisteredEvent {
	return domain.AccountRegisteredEvent{
		BasicEvent: domain.BasicEvent{
			EventID:   d.ID,
			Timestamp: d.OccurredAt,
		},
		AccountID: d.AccountID,
		Name:      d.Name,
		Email:     d.Email,
	}
}
