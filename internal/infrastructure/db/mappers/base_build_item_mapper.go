package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
)

func HydrateBuildItems(base *domain.UserBaseModel, rows []gen.BaseBuildItem, proto map[int]*domain.BuildItemPrototype) {
	for _, r := range rows {
		p := proto[int(r.PrototypeID)]
		if p == nil {
			continue
		}
		owned := domain.BaseOwnedItem{ID: r.ID, UserBaseID: base.ID}
		switch r.Status {
		case string(domain.BuildStatusPending):
			base.BuildingsPending = append(base.BuildingsPending, domain.BuildItemPending{
				BaseOwnedItem: owned,
				Prototype:     *p,
			})
		case string(domain.BuildStatusInProduction):
			var d dtos.BuildInProdDTO
			unmarshalIfValid(r.InProdData, &d)
			base.BuildingsInProduction = append(base.BuildingsInProduction, dtos.BuildInProductionFromDTO(d, owned, *p))
		case string(domain.BuildStatusPresent):
			var d dtos.BuildPresentDTO
			unmarshalIfValid(r.PresentData, &d)
			base.BuildingsPresent = append(base.BuildingsPresent, dtos.BuildPresentFromDTO(d, owned, *p))
		}
	}
}
