package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"

type TechCategory string

type TechStatus string

// TechCategory enum values
const (
	Army     TechCategory = "ARMY"
	Build    TechCategory = "BUILD"
	Base     TechCategory = "BASE"
	Politics TechCategory = "POLITICS"
)

// TechStatus enum values
const (
	TechNew        TechStatus = "NEW"
	TechInProgress TechStatus = "IN_PROGRESS"
	TechDone       TechStatus = "DONE"
)

type TechItemPrototypeDTO struct {
	ID               int           `json:"id"`
	Name             string        `json:"name"`
	Category         TechCategory  `json:"category"`
	ShortDescription string        `json:"short_description"`
	FullDescription  string        `json:"full_description"`
	Price            PriceModelDTO `json:"price"`
	ImageURL         string        `json:"image_url"`
	ResearchTime     int           `json:"research_time"`
}

type TechItemNewDTO struct {
	TechItemPrototypeDTO
}

type TechItemInProgressDTO struct {
	BaseOwnedItemDTO
	Prototype         TechItemPrototypeDTO `json:"prototype"`
	StartDate         int                  `json:"start_date"`
	CompletionDate    int                  `json:"completion_date"`
	CrystalsSkipPrice int                  `json:"crystals_skip_price"`
}

type TechItemDoneDTO struct {
	BaseOwnedItemDTO
	Prototype TechItemPrototypeDTO `json:"prototype"`
}

type TechItemCombinedDTO struct {
	New        []TechItemNewDTO        `json:"new"`
	InProgress []TechItemInProgressDTO `json:"in_progress"`
	Done       []TechItemDoneDTO       `json:"done"`
}

func mapTechPrototype(proto readmodels.TechItemPrototype) TechItemPrototypeDTO {
	return TechItemPrototypeDTO{
		ID:               proto.ID,
		Name:             proto.Name,
		Category:         TechCategory(proto.Category),
		ShortDescription: proto.ShortDescription,
		FullDescription:  proto.FullDescription,
		Price:            PriceModelFromReadModel(proto.Price),
		ImageURL:         proto.ImageURL,
		ResearchTime:     int(proto.ResearchTime),
	}
}

func TechItemsNewFromReadModels(items []*readmodels.TechItemNew) []TechItemNewDTO {
	out := make([]TechItemNewDTO, 0, len(items))
	for _, item := range items {
		out = append(out, TechItemNewDTO{TechItemPrototypeDTO: mapTechPrototype(item.Prototype)})
	}
	return out
}

func TechItemsInProgressFromReadModels(items []*readmodels.TechItemInProgress) []TechItemInProgressDTO {
	out := make([]TechItemInProgressDTO, 0, len(items))
	for _, item := range items {
		out = append(out, TechItemInProgressDTO{
			BaseOwnedItemDTO:  BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:         mapTechPrototype(item.Prototype),
			StartDate:         int(item.StartDate),
			CompletionDate:    int(item.CompletionDate),
			CrystalsSkipPrice: item.CrystalsSkipPrice,
		})
	}
	return out
}

func TechItemsDoneFromReadModels(items []*readmodels.TechItemDone) []TechItemDoneDTO {
	out := make([]TechItemDoneDTO, 0, len(items))
	for _, item := range items {
		out = append(out, TechItemDoneDTO{
			BaseOwnedItemDTO: BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:        mapTechPrototype(item.Prototype),
		})
	}
	return out
}
