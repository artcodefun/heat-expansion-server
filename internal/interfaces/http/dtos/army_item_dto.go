package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"

type ArmyCategory string

type ArmyStatus string

// ArmyCategory enum for DTOs
const (
	InfantryDTO  ArmyCategory = "INFANTRY"
	ArmoredDTO   ArmyCategory = "ARMORED"
	ArtilleryDTO ArmyCategory = "ARTILLERY"
	AviationDTO  ArmyCategory = "AVIATION"
	SpyDTO       ArmyCategory = "SPY"
	SpecialDTO   ArmyCategory = "SPECIAL"
)

// ArmyStatus enum for DTOs
const (
	ArmyNewDTO          ArmyStatus = "NEW"
	ArmyPendingDTO      ArmyStatus = "PENDING"
	ArmyInProductionDTO ArmyStatus = "IN_PRODUCTION"
	ArmyPresentDTO      ArmyStatus = "PRESENT"
)

type ArmyItemDTO struct {
	ID               int           `json:"id"`
	Name             string        `json:"name"`
	Category         ArmyCategory  `json:"category"`
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
}

type ArmyItemNewDTO struct {
	ArmyItemDTO
}

type ArmyItemPendingDTO struct {
	ArmyItemDTO
	Count int `json:"count"`
}

type ArmyItemInProductionDTO struct {
	ArmyItemDTO
	TaskID            int `json:"task_id"`
	StartDate         int `json:"start_date"`
	CompletionDate    int `json:"completion_date"`
	CrystalsSkipPrice int `json:"crystals_skip_price"`
}

type ArmyItemPresentDTO struct {
	ArmyItemDTO
	Count  int           `json:"count"`
	Refund PriceModelDTO `json:"refund"`
}

func mapArmyPrototype(proto readmodels.ArmyItemPrototype) ArmyItemDTO {
	return ArmyItemDTO{
		ID:               proto.ID,
		Name:             proto.Name,
		Category:         ArmyCategory(proto.Category),
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
	}
}

func ArmyItemsNewFromReadModels(items []*readmodels.ArmyItemNew) []ArmyItemNewDTO {
	out := make([]ArmyItemNewDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ArmyItemNewDTO{ArmyItemDTO: mapArmyPrototype(item.Prototype)})
	}
	return out
}

func ArmyItemsPendingFromReadModels(items []*readmodels.ArmyItemPending) []ArmyItemPendingDTO {
	out := make([]ArmyItemPendingDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ArmyItemPendingDTO{ArmyItemDTO: mapArmyPrototype(item.Prototype), Count: item.Count})
	}
	return out
}

func ArmyItemsInProductionFromReadModels(items []*readmodels.ArmyItemInProduction) []ArmyItemInProductionDTO {
	out := make([]ArmyItemInProductionDTO, 0, len(items))
	for _, item := range items {
		out = append(out, ArmyItemInProductionDTO{
			ArmyItemDTO:       mapArmyPrototype(item.Prototype),
			TaskID:            0,
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
			ArmyItemDTO: mapArmyPrototype(item.Prototype),
			Count:       item.Count,
			Refund:      PriceModelFromReadModel(item.Refund),
		})
	}
	return out
}
