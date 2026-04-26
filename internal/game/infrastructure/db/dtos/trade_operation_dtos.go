package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

type TradeArmyItemSnapDTO struct {
	PrototypeID int `json:"prototype_id"`
	Count       int `json:"count"`
	Capacity    int `json:"capacity"`
}

func TradeArmyItemSnapDTOFromDomain(s domain.TradeArmyItemSnap) TradeArmyItemSnapDTO {
	return TradeArmyItemSnapDTO{PrototypeID: s.PrototypeID, Count: s.Count, Capacity: s.Capacity}
}

func TradeArmyItemSnapFromDTO(d TradeArmyItemSnapDTO) domain.TradeArmyItemSnap {
	return domain.TradeArmyItemSnap{PrototypeID: d.PrototypeID, Count: d.Count, Capacity: d.Capacity}
}

type TradeStorageItemSnapDTO struct {
	ItemID      uuid.UUID              `json:"item_id"`
	PrototypeID int                    `json:"prototype_id"`
	Category    domain.StorageCategory `json:"category"`
}

func TradeStorageItemSnapDTOFromDomain(s domain.TradeStorageItemSnap) TradeStorageItemSnapDTO {
	return TradeStorageItemSnapDTO{ItemID: s.ItemID, PrototypeID: s.PrototypeID, Category: s.Category}
}

func TradeStorageItemSnapFromDTO(d TradeStorageItemSnapDTO) domain.TradeStorageItemSnap {
	return domain.TradeStorageItemSnap{ItemID: d.ItemID, PrototypeID: d.PrototypeID, Category: d.Category}
}

type TradePayloadDTO struct {
	Resources PriceDTO                  `json:"resources"`
	Storage   []TradeStorageItemSnapDTO `json:"storage"`
	Army      []TradeArmyItemSnapDTO    `json:"army"`
}

func TradePayloadDTOFromDomain(p domain.TradePayload) TradePayloadDTO {
	storage := make([]TradeStorageItemSnapDTO, 0, len(p.Storage))
	for _, s := range p.Storage {
		storage = append(storage, TradeStorageItemSnapDTOFromDomain(s))
	}
	army := make([]TradeArmyItemSnapDTO, 0, len(p.Army))
	for _, a := range p.Army {
		army = append(army, TradeArmyItemSnapDTOFromDomain(a))
	}
	return TradePayloadDTO{Resources: PriceDTOFromDomain(p.Resources), Storage: storage, Army: army}
}

func TradePayloadFromDTO(d TradePayloadDTO) domain.TradePayload {
	storage := make([]domain.TradeStorageItemSnap, 0, len(d.Storage))
	for _, s := range d.Storage {
		storage = append(storage, TradeStorageItemSnapFromDTO(s))
	}
	army := make([]domain.TradeArmyItemSnap, 0, len(d.Army))
	for _, a := range d.Army {
		army = append(army, TradeArmyItemSnapFromDTO(a))
	}
	return domain.TradePayload{Resources: PriceFromDTO(d.Resources), Storage: storage, Army: army}
}
