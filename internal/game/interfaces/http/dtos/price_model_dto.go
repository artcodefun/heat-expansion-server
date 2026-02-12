package dtos

import "github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs/readmodels"

type PriceModelDTO struct {
	Credits    int `json:"credits"`
	Iron       int `json:"iron"`
	Titanium   int `json:"titanium"`
	Antimatter int `json:"antimatter"`
}

func PriceModelFromReadModel(m readmodels.PriceModel) PriceModelDTO {
	return PriceModelDTO{
		Credits:    m.Credits,
		Iron:       m.Iron,
		Titanium:   m.Titanium,
		Antimatter: m.Antimatter,
	}
}
