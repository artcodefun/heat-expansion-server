package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

type Vector2iDTO struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func Vector2iDTOFromDomain(v domain.Vector2i) Vector2iDTO {
	return Vector2iDTO{X: v.X, Y: v.Y}
}

func (v Vector2iDTO) ToDomain() domain.Vector2i {
	return domain.Vector2i{X: v.X, Y: v.Y}
}
