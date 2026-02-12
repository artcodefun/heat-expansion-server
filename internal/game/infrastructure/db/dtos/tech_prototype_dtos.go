package dtos

import "github.com/artcodefun/heat-expansion-server/internal/game/domain"

type TechImprovementDTO struct {
	Type     string `json:"type"`
	Value    int    `json:"value"`
	MaxLevel *int   `json:"max_level"`
}

func TechImprovementDTOFromDomain(e *domain.TechImprovement) *TechImprovementDTO {
	if e == nil {
		return nil
	}
	return &TechImprovementDTO{Type: string(e.Type), Value: e.Value, MaxLevel: e.MaxLevel}
}

func TechImprovementFromDTO(d *TechImprovementDTO) *domain.TechImprovement {
	if d == nil {
		return nil
	}
	return &domain.TechImprovement{Type: domain.ImprovementType(d.Type), Value: d.Value, MaxLevel: d.MaxLevel}
}
