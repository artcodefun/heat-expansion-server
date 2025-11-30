package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
)

func HydrateArmyItems(base *domain.UserBaseModel, rows []gen.BaseArmyItem, proto map[int]*domain.ArmyItemPrototype) {
	for _, r := range rows {
		p := proto[int(r.PrototypeID)]
		if p == nil {
			continue
		}
		owned := domain.BaseOwnedItem{ID: r.ID, UserBaseID: base.ID}
		switch r.Status {
		case string(domain.ArmyStatusPending):
			var d dtos.ArmyPendingDTO
			unmarshalIfValid(r.PendingData, &d)
			base.ArmiesPending = append(base.ArmiesPending, dtos.ArmyPendingFromDTO(d, owned, *p))
		case string(domain.ArmyStatusInProduction):
			var d dtos.ArmyInProdDTO
			unmarshalIfValid(r.InProdData, &d)
			base.ArmiesInProduction = append(base.ArmiesInProduction, dtos.ArmyInProductionFromDTO(d, owned, *p))
		case string(domain.ArmyStatusPresent):
			var d dtos.ArmyPresentDTO
			unmarshalIfValid(r.PresentData, &d)
			base.ArmiesPresent = append(base.ArmiesPresent, dtos.ArmyPresentFromDTO(d, owned, *p))
		case string(domain.ArmyStatusDeployed):
			var d dtos.ArmyDeployedDTO
			unmarshalIfValid(r.DeployedData, &d)
			base.ArmiesDeployed = append(base.ArmiesDeployed, dtos.ArmyDeployedFromDTO(d, owned, *p))
		}
	}
}
