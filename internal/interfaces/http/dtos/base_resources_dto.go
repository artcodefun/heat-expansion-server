package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"

type BaseResourcesDTO struct {
	Credits            int     `json:"credits"`
	CreditsCapacity    int     `json:"credits_capacity"`
	CreditsProduction  float64 `json:"credits_production"`
	Iron               int     `json:"iron"`
	IronCapacity       int     `json:"iron_capacity"`
	IronProduction     float64 `json:"iron_production"`
	Titanium           int     `json:"titanium"`
	TitaniumCapacity   int     `json:"titanium_capacity"`
	TitaniumProduction float64 `json:"titanium_production"`
	Antimatter         int     `json:"antimatter"`
	AntimatterCapacity int     `json:"antimatter_capacity"`
	Defence            int     `json:"defence"`
	Attack             int     `json:"attack"`
	Space              int     `json:"space"`
	SpaceCapacity      int     `json:"space_capacity"`
	Crystals           int     `json:"crystals"`
}

// BaseResourcesFromReadModel maps a readmodel.UserBaseStats to BaseResourcesDTO.
func BaseResourcesFromReadModel(m *readmodels.UserBaseStats) BaseResourcesDTO {
	return BaseResourcesDTO{
		Credits:            m.Credits,
		CreditsCapacity:    m.CreditsCapacity,
		CreditsProduction:  m.CreditsProduction,
		Iron:               m.Iron,
		IronCapacity:       m.IronCapacity,
		IronProduction:     m.IronProduction,
		Titanium:           m.Titanium,
		TitaniumCapacity:   m.TitaniumCapacity,
		TitaniumProduction: m.TitaniumProduction,
		Antimatter:         m.Antimatter,
		AntimatterCapacity: m.AntimatterCapacity,
		Defence:            m.Defence,
		Attack:             m.Attack,
		Space:              m.Space,
		SpaceCapacity:      m.SpaceCapacity,
	}
}
