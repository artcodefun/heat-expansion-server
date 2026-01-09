package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/domain"

type StoragePresentDTO struct {
	Refund      PriceDTO `json:"refund"`
	ActivatedAt *int64   `json:"activated_at,omitempty"`
}

func StoragePresentDTOFromDomain(s domain.StorageItemPresent) StoragePresentDTO {
	return StoragePresentDTO{
		Refund:      PriceDTOFromDomain(s.Refund),
		ActivatedAt: s.ActivatedAt,
	}
}

func StoragePresentFromDTO(d StoragePresentDTO, owned domain.BaseOwnedItem, proto domain.StorageItemPrototype) domain.StorageItemPresent {
	return domain.StorageItemPresent{
		BaseOwnedItem: owned,
		Prototype:     proto,
		Refund:        PriceFromDTO(d.Refund),
		ActivatedAt:   d.ActivatedAt,
	}
}
