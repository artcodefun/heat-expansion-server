package dtos

import "github.com/artcodefun/heat-expansion-api/internal/game/core/domain"

type StoragePresentDTO struct {
	ExpiresAt *int64 `json:"expires_at,omitempty"`
	IsActive  bool   `json:"is_active"`
}

func StoragePresentDTOFromDomain(s domain.StorageItemPresent) StoragePresentDTO {
	return StoragePresentDTO{
		ExpiresAt: s.ExpiresAt,
		IsActive:  s.IsActive,
	}
}

func StoragePresentFromDTO(d StoragePresentDTO, owned domain.BaseOwnedItem, proto domain.StorageItemPrototype) domain.StorageItemPresent {
	return domain.StorageItemPresent{
		BaseOwnedItem: owned,
		Prototype:     proto,
		ExpiresAt:     d.ExpiresAt,
		IsActive:      d.IsActive,
	}
}
