package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

type StoragePresentDTO struct {
	Refund PriceDTO `json:"refund"`
}

func StoragePresentDTOFromDomain(s domain.StorageItemPresent) StoragePresentDTO {
	return StoragePresentDTO{Refund: PriceDTOFromDomain(s.Refund)}
}

func StoragePresentFromDTO(d StoragePresentDTO, owned domain.BaseOwnedItem, proto domain.StorageItemPrototype) domain.StorageItemPresent {
	return domain.StorageItemPresent{BaseOwnedItem: owned, Prototype: proto, Refund: PriceFromDTO(d.Refund)}
}
