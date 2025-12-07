package dtos

import (
	"strings"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
)

type BuildCategory string

type BuildStatus string

// BuildCategory enum values
const (
	Control      BuildCategory = "CONTROL"
	Resources    BuildCategory = "RESOURCES"
	Defense      BuildCategory = "DEFENSE"
	Military     BuildCategory = "MILITARY"
	Intelligence BuildCategory = "INTELLIGENCE"
)

// BuildStatus enum values
const (
	BuildNew          BuildStatus = "NEW"
	BuildPending      BuildStatus = "PENDING"
	BuildInProduction BuildStatus = "IN_PRODUCTION"
	BuildPresent      BuildStatus = "PRESENT"
)

type BuildItemPrototypeDTO struct {
	ID               int           `json:"id"`
	Name             string        `json:"name"`
	Category         BuildCategory `json:"category"`
	ShortDescription string        `json:"short_description"`
	FullDescription  string        `json:"full_description"`
	Price            PriceModelDTO `json:"price"`
	Space            int           `json:"space"`
	ImageURL         string        `json:"image_url"`
	ProductionTime   int           `json:"production_time"`
}

type BuildItemNewDTO struct {
	BuildItemPrototypeDTO
}

type BuildItemPendingDTO struct {
	BaseOwnedItemDTO
	Prototype BuildItemPrototypeDTO `json:"prototype"`
}

type BuildItemInProductionDTO struct {
	BaseOwnedItemDTO
	Prototype         BuildItemPrototypeDTO `json:"prototype"`
	StartDate         int                   `json:"start_date"`
	CompletionDate    int                   `json:"completion_date"`
	CrystalsSkipPrice int                   `json:"crystals_skip_price"`
}

type BuildItemPresentDTO struct {
	BaseOwnedItemDTO
	Prototype BuildItemPrototypeDTO `json:"prototype"`
	Refund    PriceModelDTO         `json:"refund"`
}

func mapBuildItemPrototype(proto readmodels.BuildItemPrototype) BuildItemPrototypeDTO {
	return BuildItemPrototypeDTO{
		ID:               proto.ID,
		Name:             proto.Name,
		Category:         BuildCategory(proto.Category),
		ShortDescription: proto.ShortDescription,
		FullDescription:  proto.FullDescription,
		Price:            PriceModelFromReadModel(proto.Price),
		Space:            proto.Space,
		ImageURL:         proto.ImageURL,
		ProductionTime:   int(proto.ProductionTime),
	}
}

func BuildItemsNewFromReadModels(items []*readmodels.BuildItemNew) []BuildItemNewDTO {
	out := make([]BuildItemNewDTO, 0, len(items))
	for _, item := range items {
		out = append(out, BuildItemNewDTO{BuildItemPrototypeDTO: mapBuildItemPrototype(item.Prototype)})
	}
	return out
}

func BuildItemsPendingFromReadModels(items []*readmodels.BuildItemPending) []BuildItemPendingDTO {
	out := make([]BuildItemPendingDTO, 0, len(items))
	for _, item := range items {
		out = append(out, BuildItemPendingDTO{
			BaseOwnedItemDTO: BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:        mapBuildItemPrototype(item.Prototype),
		})
	}
	return out
}

func BuildItemsInProductionFromReadModels(items []*readmodels.BuildItemInProduction) []BuildItemInProductionDTO {
	out := make([]BuildItemInProductionDTO, 0, len(items))
	for _, item := range items {
		out = append(out, BuildItemInProductionDTO{
			BaseOwnedItemDTO:  BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:         mapBuildItemPrototype(item.Prototype),
			StartDate:         int(item.StartDate),
			CompletionDate:    int(item.CompletionDate),
			CrystalsSkipPrice: item.CrystalsSkipPrice,
		})
	}
	return out
}

func BuildItemsPresentFromReadModels(items []*readmodels.BuildItemPresent) []BuildItemPresentDTO {
	out := make([]BuildItemPresentDTO, 0, len(items))
	for _, item := range items {
		out = append(out, BuildItemPresentDTO{
			BaseOwnedItemDTO: BaseOwnedItemDTOFromReadModel(item.BaseOwnedItem),
			Prototype:        mapBuildItemPrototype(item.Prototype),
			Refund:           PriceModelFromReadModel(item.Refund),
		})
	}
	return out
}

// BuildCategoryFromDTO normalizes a DTO category string to the read-model type.
func BuildCategoryFromDTO(value string) readmodels.BuildCategory {
	return readmodels.BuildCategory(strings.ToUpper(strings.TrimSpace(value)))
}
