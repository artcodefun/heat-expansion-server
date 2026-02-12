package dtos

import "github.com/artcodefun/heat-expansion-server/internal/game/domain"

type ArmyPendingDTO struct {
	Count int `json:"count"`
}

func ArmyPendingDTOFromDomain(a domain.ArmyItemPending) ArmyPendingDTO {
	return ArmyPendingDTO{Count: a.Count}
}

func ArmyPendingFromDTO(d ArmyPendingDTO, owned domain.BaseOwnedItem, proto domain.ArmyItemPrototype) domain.ArmyItemPending {
	return domain.ArmyItemPending{BaseOwnedItem: owned, Prototype: proto, Count: d.Count}
}

type ArmyInProdDTO struct {
	StartDate         int64 `json:"start_date"`
	CompletionDate    int64 `json:"completion_date"`
	CrystalsSkipPrice int   `json:"crystals_skip_price"`
}

func ArmyInProdDTOFromDomain(a domain.ArmyItemInProduction) ArmyInProdDTO {
	return ArmyInProdDTO{StartDate: a.StartDate, CompletionDate: a.CompletionDate, CrystalsSkipPrice: a.CrystalsSkipPrice}
}

func ArmyInProductionFromDTO(d ArmyInProdDTO, owned domain.BaseOwnedItem, proto domain.ArmyItemPrototype) domain.ArmyItemInProduction {
	return domain.ArmyItemInProduction{BaseOwnedItem: owned, Prototype: proto, StartDate: d.StartDate, CompletionDate: d.CompletionDate, CrystalsSkipPrice: d.CrystalsSkipPrice}
}

type ArmyPresentDTO struct {
	Count  int      `json:"count"`
	Refund PriceDTO `json:"refund"`
}

func ArmyPresentDTOFromDomain(a domain.ArmyItemPresent) ArmyPresentDTO {
	return ArmyPresentDTO{Count: a.Count, Refund: PriceDTOFromDomain(a.Refund)}
}

func ArmyPresentFromDTO(d ArmyPresentDTO, owned domain.BaseOwnedItem, proto domain.ArmyItemPrototype) domain.ArmyItemPresent {
	return domain.ArmyItemPresent{BaseOwnedItem: owned, Prototype: proto, Count: d.Count, Refund: PriceFromDTO(d.Refund)}
}

type ArmyDeployedDTO struct {
	OperationID int `json:"operation_id"`
	Count       int `json:"count"`
}

func ArmyDeployedDTOFromDomain(a domain.ArmyItemDeployed) ArmyDeployedDTO {
	return ArmyDeployedDTO{OperationID: a.OperationID, Count: a.Count}
}

func ArmyDeployedFromDTO(d ArmyDeployedDTO, owned domain.BaseOwnedItem, proto domain.ArmyItemPrototype) domain.ArmyItemDeployed {
	return domain.ArmyItemDeployed{BaseOwnedItem: owned, Prototype: proto, OperationID: d.OperationID, Count: d.Count}
}
