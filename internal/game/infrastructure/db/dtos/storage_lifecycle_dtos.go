package dtos

import "github.com/artcodefun/heat-expansion-server/internal/game/domain"

type StoragePresentDTO struct {
	ExpiresAt *int64 `json:"expires_at,omitempty"`
	IsActive  bool   `json:"is_active"`
}

type StorageDeployedDTO struct {
	OperationKind string `json:"operation_kind"`
	OperationID   int    `json:"operation_id"`
	ExpiresAt     *int64 `json:"expires_at,omitempty"`
	IsActive      bool   `json:"is_active"`
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

func StorageDeployedDTOFromDomain(s domain.StorageItemDeployed) StorageDeployedDTO {
	operationKind := string(s.OperationKind)
	if operationKind == "" {
		operationKind = string(domain.OperationKindTrade)
	}
	return StorageDeployedDTO{
		OperationKind: operationKind,
		OperationID:   s.OperationID,
		ExpiresAt:     s.ExpiresAt,
		IsActive:      s.IsActive,
	}
}

func StorageDeployedFromDTO(d StorageDeployedDTO, owned domain.BaseOwnedItem, proto domain.StorageItemPrototype) domain.StorageItemDeployed {
	operationKind := domain.OperationKind(d.OperationKind)
	if operationKind == "" {
		operationKind = domain.OperationKindTrade
	}
	return domain.StorageItemDeployed{
		BaseOwnedItem: owned,
		Prototype:     proto,
		OperationKind: operationKind,
		OperationID:   d.OperationID,
		ExpiresAt:     d.ExpiresAt,
		IsActive:      d.IsActive,
	}
}
