package dtos

import (
	"fmt"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
)

// UserBaseDTO is a lightweight representation of a user's base for listings.
type UserBaseDTO struct {
	ID          int         `json:"id"`
	Coordinates Vector2iDTO `json:"coordinates"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	ImageURL    string      `json:"image_url,omitempty"`
}

func UserBaseFromReadModel(m *readmodels.UserBaseModel, tr ports.Translator, locale string) UserBaseDTO {
	return UserBaseDTO{
		ID:          m.ID,
		Coordinates: Vector2iFromReadModel(m.Coordinates),
		Name:        fmt.Sprintf("%s #%d", tr.T(locale, m.LocationDetails.Name, nil), m.ID),
		Description: tr.T(locale, m.LocationDetails.Description, nil),
		ImageURL:    m.LocationDetails.ImageURL,
	}
}

func UserBasesFromReadModels(models []*readmodels.UserBaseModel, tr ports.Translator, locale string) []UserBaseDTO {
	out := make([]UserBaseDTO, 0, len(models))
	for _, m := range models {
		out = append(out, UserBaseFromReadModel(m, tr, locale))
	}
	return out
}
