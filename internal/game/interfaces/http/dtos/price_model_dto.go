package dtos

import "github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"

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

func PriceModelFromDTO(m PriceModelDTO) readmodels.PriceModel {
	return readmodels.PriceModel{
		Credits:    m.Credits,
		Iron:       m.Iron,
		Titanium:   m.Titanium,
		Antimatter: m.Antimatter,
	}
}
