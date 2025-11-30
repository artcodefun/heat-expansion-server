package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

type BuildInProdDTO struct {
	StartDate         int64 `json:"start_date"`
	CompletionDate    int64 `json:"completion_date"`
	CrystalsSkipPrice int   `json:"crystals_skip_price"`
}

func BuildInProdDTOFromDomain(b domain.BuildItemInProduction) BuildInProdDTO {
	return BuildInProdDTO{StartDate: b.StartDate, CompletionDate: b.CompletionDate, CrystalsSkipPrice: b.CrystalsSkipPrice}
}

func BuildInProductionFromDTO(d BuildInProdDTO, owned domain.BaseOwnedItem, proto domain.BuildItemPrototype) domain.BuildItemInProduction {
	return domain.BuildItemInProduction{BaseOwnedItem: owned, Prototype: proto, StartDate: d.StartDate, CompletionDate: d.CompletionDate, CrystalsSkipPrice: d.CrystalsSkipPrice}
}

type BuildPresentDTO struct {
	Refund PriceDTO `json:"refund"`
}

func BuildPresentDTOFromDomain(b domain.BuildItemPresent) BuildPresentDTO {
	return BuildPresentDTO{Refund: PriceDTOFromDomain(b.Refund)}
}

func BuildPresentFromDTO(d BuildPresentDTO, owned domain.BaseOwnedItem, proto domain.BuildItemPrototype) domain.BuildItemPresent {
	return domain.BuildItemPresent{BaseOwnedItem: owned, Prototype: proto, Refund: PriceFromDTO(d.Refund)}
}
