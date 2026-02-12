package dtos

import "github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"

type BaseResourcesDTO struct {
	Credits              float64 `json:"credits"`
	CreditsCapacity      int     `json:"credits_capacity"`
	CreditsProduction    float64 `json:"credits_production"`
	Iron                 float64 `json:"iron"`
	IronCapacity         int     `json:"iron_capacity"`
	IronProduction       float64 `json:"iron_production"`
	Titanium             float64 `json:"titanium"`
	TitaniumCapacity     int     `json:"titanium_capacity"`
	TitaniumProduction   float64 `json:"titanium_production"`
	Antimatter           float64 `json:"antimatter"`
	AntimatterCapacity   int     `json:"antimatter_capacity"`
	AntimatterProduction float64 `json:"antimatter_production"`
	Defence              int     `json:"defence"`
	Attack               int     `json:"attack"`
	Space                int     `json:"space"`
	MaxSpace             int     `json:"max_space"`
	CalculationTimestamp int64   `json:"calculation_timestamp"`
}

// BaseResourcesFromReadModel maps a readmodel.UserBaseStats to BaseResourcesDTO.
func BaseResourcesFromReadModel(m *readmodels.UserBaseStats) BaseResourcesDTO {
	return BaseResourcesDTO{
		Credits:              m.Credits,
		CreditsCapacity:      m.CreditsCapacity,
		CreditsProduction:    m.CreditsProduction,
		Iron:                 m.Iron,
		IronCapacity:         m.IronCapacity,
		IronProduction:       m.IronProduction,
		Titanium:             m.Titanium,
		TitaniumCapacity:     m.TitaniumCapacity,
		TitaniumProduction:   m.TitaniumProduction,
		Antimatter:           m.Antimatter,
		AntimatterCapacity:   m.AntimatterCapacity,
		AntimatterProduction: m.AntimatterProduction,
		Defence:              m.Defence,
		Attack:               m.Attack,
		Space:                m.Space,
		MaxSpace:             m.MaxSpace,
		CalculationTimestamp: m.CalculationTimestamp,
	}
}
