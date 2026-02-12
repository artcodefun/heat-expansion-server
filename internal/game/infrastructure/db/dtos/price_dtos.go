package dtos

import "github.com/artcodefun/heat-expansion-api/internal/game/domain"

type PriceDTO struct {
	Credits    int `json:"credits"`
	Iron       int `json:"iron"`
	Titanium   int `json:"titanium"`
	Antimatter int `json:"antimatter"`
}

func PriceDTOFromDomain(p domain.PriceModel) PriceDTO {
	return PriceDTO{Credits: p.Credits, Iron: p.Iron, Titanium: p.Titanium, Antimatter: p.Antimatter}
}
func PriceFromDTO(d PriceDTO) domain.PriceModel {
	return domain.PriceModel{Credits: d.Credits, Iron: d.Iron, Titanium: d.Titanium, Antimatter: d.Antimatter}
}
