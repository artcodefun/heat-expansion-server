package dtos

import (
	"strings"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
)

type ArmyCategory string

// ArmyCategory enum values
const (
	Infantry  ArmyCategory = "INFANTRY"
	Armored   ArmyCategory = "ARMORED"
	Artillery ArmyCategory = "ARTILLERY"
	Aviation  ArmyCategory = "AVIATION"
	Spy       ArmyCategory = "SPY"
	Special   ArmyCategory = "SPECIAL"
)

type ArmyStatus string

// ArmyStatus enum values
const (
	ArmyNew          ArmyStatus = "NEW"
	ArmyPending      ArmyStatus = "PENDING"
	ArmyInProduction ArmyStatus = "IN_PRODUCTION"
	ArmyPresent      ArmyStatus = "PRESENT"
)

type Faction string

const (
	FactionExoCoalition      Faction = "EXO_COALITION"   // Playable (Human)
	FactionMarauders         Faction = "MARAUDERS"       // NPC: Credits
	FactionFerrousSwarm      Faction = "FERROUS_SWARM"   // NPC: Iron
	FactionTitanArachnids    Faction = "TITAN_ARACHNIDS" // NPC: Titanium
	FactionVoidEcho          Faction = "VOID_ECHO"       // NPC: Antimatter
	FactionCustodianProtocol Faction = "CUSTODIAN"       // NPC: Dangerous (Artifacts)
	FactionScorchWalkers     Faction = "SCORCH_WALKERS"  // NPC: Dangerous (Buffs)
	FactionObsidianSentinels Faction = "OBSIDIAN"        // NPC: Dangerous (Trophies)
	FactionNeuralWormApex    Faction = "NEURAL_WORM"     // NPC: Dangerous (Intel)
)

type ArmyItemPrototypeDTO struct {
	ID               int           `json:"id"`
	Name             string        `json:"name"`
	Category         ArmyCategory  `json:"category"`
	Faction          Faction       `json:"faction"`
	ShortDescription string        `json:"short_description"`
	FullDescription  string        `json:"full_description"`
	Price            PriceModelDTO `json:"price"`
	Space            int           `json:"space"`
	ImageURL         string        `json:"image_url"`
	Attack           int           `json:"attack"`
	Defence          int           `json:"defence"`
	Capacity         int           `json:"capacity"`
	Stealth          int           `json:"stealth"`
	Speed            int           `json:"speed"`
	ProductionTime   int           `json:"production_time"`
}

type ArmyItemNewDTO struct {
	ArmyItemPrototypeDTO
}

type ArmyItemPendingDTO struct {
	BaseOwnedItemDTO
	Prototype ArmyItemPrototypeDTO `json:"prototype"`
	Count     int                  `json:"count"`
}

type ArmyItemInProductionDTO struct {
	BaseOwnedItemDTO
	Prototype         ArmyItemPrototypeDTO `json:"prototype"`
	StartDate         int                  `json:"start_date"`
	CompletionDate    int                  `json:"completion_date"`
	CrystalsSkipPrice int                  `json:"crystals_skip_price"`
}

type ArmyItemPresentDTO struct {
	BaseOwnedItemDTO
	Prototype ArmyItemPrototypeDTO `json:"prototype"`
	Count     int                  `json:"count"`
	Refund    PriceModelDTO        `json:"refund"`
}

func mapArmyPrototype(proto readmodels.ArmyItemPrototype) ArmyItemPrototypeDTO {
	return ArmyItemPrototypeDTO{
		ID:               proto.ID,
		Name:             proto.Name,
		Category:         ArmyCategory(proto.Category),
		Faction:          Faction(proto.Faction),
		ShortDescription: proto.ShortDescription,
		FullDescription:  proto.FullDescription,
		Price:            PriceModelFromReadModel(proto.Price),
		Space:            proto.Space,
		ImageURL:         proto.ImageURL,
		Attack:           proto.Attack,
		Defence:          proto.Defence,
		Capacity:         proto.Capacity,
		Stealth:          proto.Stealth,
		Speed:            proto.Speed,
		ProductionTime:   int(proto.ProductionTime),
	}
}

func ArmyItemsNewFromReadModels(items []*readmodels.ArmyItemNew) []ArmyItemNewDTO {
	out := make([]ArmyItemNewDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ArmyItemNewDTO{ArmyItemPrototypeDTO: mapArmyPrototype(item.Prototype)})
	}
	return out
}

func ArmyItemsPendingFromReadModels(items []*readmodels.ArmyItemPending) []ArmyItemPendingDTO {
	out := make([]ArmyItemPendingDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ArmyItemPendingDTO{
			BaseOwnedItemDTO: BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:        mapArmyPrototype(item.Prototype),
			Count:            item.Count,
		})
	}
	return out
}

func ArmyItemsInProductionFromReadModels(items []*readmodels.ArmyItemInProduction) []ArmyItemInProductionDTO {
	out := make([]ArmyItemInProductionDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ArmyItemInProductionDTO{
			BaseOwnedItemDTO:  BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:         mapArmyPrototype(item.Prototype),
			StartDate:         int(item.StartDate),
			CompletionDate:    int(item.CompletionDate),
			CrystalsSkipPrice: item.CrystalsSkipPrice,
		})
	}
	return out
}

func ArmyItemsPresentFromReadModels(items []*readmodels.ArmyItemPresent) []ArmyItemPresentDTO {
	out := make([]ArmyItemPresentDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ArmyItemPresentDTO{
			BaseOwnedItemDTO: BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:        mapArmyPrototype(item.Prototype),
			Count:            item.Count,
			Refund:           PriceModelFromReadModel(item.Refund),
		})
	}
	return out
}

// ArmyCategoryFromDTO normalizes a request category string to the read-model type.
func ArmyCategoryFromDTO(value string) readmodels.ArmyCategory {
	return readmodels.ArmyCategory(strings.ToUpper(strings.TrimSpace(value)))
}
