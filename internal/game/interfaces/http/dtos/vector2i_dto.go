package dtos

import "github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs/readmodels"

// Vector2iDTO mirrors readmodels.Vector2i for serialization.
type Vector2iDTO struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func Vector2iFromReadModel(v readmodels.Vector2i) Vector2iDTO {
	return Vector2iDTO{X: v.X, Y: v.Y}
}
