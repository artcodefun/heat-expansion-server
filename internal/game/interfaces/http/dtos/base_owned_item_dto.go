package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/google/uuid"
)

// BaseOwnedItemDTO mirrors readmodels.BaseOwnedItem for HTTP responses.
type BaseOwnedItemDTO struct {
	ID uuid.UUID `json:"id"`
}

func BaseOwnedItemDTOFromReadModel(owned readmodels.BaseOwnedItem) BaseOwnedItemDTO {
	return BaseOwnedItemDTO{ID: owned.ID}
}
