package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"

type TechCategory string

type TechStatus string

// TechCategory enum for DTOs
const (
	ArmyDTO     TechCategory = "ARMY"
	BuildDTO    TechCategory = "BUILD"
	BaseDTO     TechCategory = "BASE"
	PoliticsDTO TechCategory = "POLITICS"
)

// TechStatus enum for DTOs
const (
	TechNewDTO        TechStatus = "NEW"
	TechInProgressDTO TechStatus = "IN_PROGRESS"
	TechDoneDTO       TechStatus = "DONE"
)

type TechItemDTO struct {
	ID               int           `json:"id"`
	Name             string        `json:"name"`
	Category         TechCategory  `json:"category"`
	ShortDescription string        `json:"short_description"`
	FullDescription  string        `json:"full_description"`
	Price            PriceModelDTO `json:"price"`
	ImageURL         string        `json:"image_url"`
}

type TechItemNewDTO struct {
	TechItemDTO
}

type TechItemInProgressDTO struct {
	TechItemDTO
	TaskID            int `json:"task_id"`
	StartDate         int `json:"start_date"`
	CompletionDate    int `json:"completion_date"`
	CrystalsSkipPrice int `json:"crystals_skip_price"`
}

type TechItemDoneDTO struct {
	TechItemDTO
}

type TechItemCombinedDTO struct {
	New        []TechItemNewDTO        `json:"new"`
	InProgress []TechItemInProgressDTO `json:"in_progress"`
	Done       []TechItemDoneDTO       `json:"done"`
}

func mapTechPrototype(proto readmodels.TechItemPrototype) TechItemDTO {
	return TechItemDTO{
		ID:               proto.ID,
		Name:             proto.Name,
		Category:         TechCategory(proto.Category),
		ShortDescription: proto.ShortDescription,
		FullDescription:  proto.FullDescription,
		Price:            PriceModelFromReadModel(proto.Price),
		ImageURL:         proto.ImageURL,
	}
}

func TechItemsNewFromReadModels(items []*readmodels.TechItemNew) []TechItemNewDTO {
	out := make([]TechItemNewDTO, 0, len(items))
	for _, item := range items {
		out = append(out, TechItemNewDTO{TechItemDTO: mapTechPrototype(item.Prototype)})
	}
	return out
}

func TechItemsInProgressFromReadModels(items []*readmodels.TechItemInProgress) []TechItemInProgressDTO {
	out := make([]TechItemInProgressDTO, 0, len(items))
	for _, item := range items {
		out = append(out, TechItemInProgressDTO{
			TechItemDTO:       mapTechPrototype(item.Prototype),
			TaskID:            0,
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
		out = append(out, TechItemDoneDTO{TechItemDTO: mapTechPrototype(item.Prototype)})
	}
	return out
}
