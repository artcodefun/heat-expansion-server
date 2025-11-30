package dtos

import "github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"

type StorageCategory string

// StorageCategory enum for DTOs
const (
	BuffDTO     StorageCategory = "BUFF"
	ArtifactDTO StorageCategory = "ARTIFACT"
)

type StorageItemDTO struct {
	ID               int             `json:"id"`
	Name             string          `json:"name"`
	Category         StorageCategory `json:"category"`
	ShortDescription string          `json:"short_description"`
	FullDescription  string          `json:"full_description"`
	ImageURL         string          `json:"image_url"`
}

type StorageItemPresentDTO struct {
	StorageItemDTO
	ItemID string        `json:"item_id"`
	Refund PriceModelDTO `json:"refund"`
}

func mapStorageItemPrototype(proto readmodels.StorageItemPrototype) StorageItemDTO {
	return StorageItemDTO{
		ID:               proto.ID,
		Name:             proto.Name,
		Category:         StorageCategory(proto.Category),
		ShortDescription: proto.ShortDescription,
		FullDescription:  proto.FullDescription,
		ImageURL:         proto.ImageURL,
	}
}

func StorageItemsPresentFromReadModels(items []*readmodels.StorageItemPresent) []StorageItemPresentDTO {
	out := make([]StorageItemPresentDTO, 0, len(items))
	for _, item := range items {
		out = append(out, StorageItemPresentDTO{
			StorageItemDTO: mapStorageItemPrototype(item.Prototype),
			ItemID:         item.BaseOwnedItem.ID.String(),
			Refund:         PriceModelFromReadModel(item.Refund),
		})
	}
	return out
}
