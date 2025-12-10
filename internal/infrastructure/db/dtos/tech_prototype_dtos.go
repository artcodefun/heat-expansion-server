package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

type TechnologyEffectDTO struct {
	Type  string `json:"type"`
	Value int    `json:"value"`
}

func TechnologyEffectDTOFromDomain(e domain.TechnologyEffect) TechnologyEffectDTO {
	return TechnologyEffectDTO{Type: string(e.EffectType), Value: e.Value}
}
func TechnologyEffectFromDTO(d TechnologyEffectDTO) domain.TechnologyEffect {
	return domain.TechnologyEffect{EffectType: domain.EffectType(d.Type), Value: d.Value}
}
