package dtos

import "github.com/artcodefun/heat-expansion-api/internal/game/application/cqrs/readmodels"

// UserBaseDTO is a lightweight representation of a user's base for listings.
type UserBaseDTO struct {
	ID          int         `json:"id"`
	Coordinates Vector2iDTO `json:"coordinates"`
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	ImageURL    string      `json:"image_url,omitempty"`
}

func UserBaseFromReadModel(m *readmodels.UserBaseModel) UserBaseDTO {
	return UserBaseDTO{
		ID:          m.ID,
		Coordinates: Vector2iFromReadModel(m.Coordinates),
		Name:        m.LocationDetails.Name,
		Description: m.LocationDetails.Description,
		ImageURL:    m.LocationDetails.ImageURL,
	}
}

func UserBasesFromReadModels(models []*readmodels.UserBaseModel) []UserBaseDTO {
	out := make([]UserBaseDTO, 0, len(models))
	for _, m := range models {
		out = append(out, UserBaseFromReadModel(m))
	}
	return out
}
