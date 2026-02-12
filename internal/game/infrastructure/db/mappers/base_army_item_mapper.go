package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/game/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/gen"
	"github.com/sqlc-dev/pqtype"
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

// DehydrateArmyItems converts the in-memory aggregate collections into insert params
// for the base_army_items table. This keeps persistence logic thin in the repo.
func DehydrateArmyItems(base *domain.UserBaseModel) []gen.InsertBaseArmyItemParams {
	now := domain.NowUnix()
	out := make([]gen.InsertBaseArmyItemParams, 0,
		len(base.ArmiesPending)+len(base.ArmiesInProduction)+len(base.ArmiesPresent)+len(base.ArmiesDeployed))

	// Pending
	for _, it := range base.ArmiesPending {
		pRaw := BuildArmyPendingRaw(it)
		out = append(out, gen.InsertBaseArmyItemParams{
			ID:           it.ID,
			BaseID:       int64(base.ID),
			PrototypeID:  int64(it.Prototype.ID),
			Status:       string(domain.ArmyStatusPending),
			PendingData:  pRaw,
			InProdData:   pqtype.NullRawMessage{Valid: false},
			PresentData:  pqtype.NullRawMessage{Valid: false},
			DeployedData: pqtype.NullRawMessage{Valid: false},
			CreatedAt:    now,
		})
	}
	// In Production
	for _, it := range base.ArmiesInProduction {
		prodRaw := BuildArmyInProdRaw(it)
		out = append(out, gen.InsertBaseArmyItemParams{
			ID:           it.ID,
			BaseID:       int64(base.ID),
			PrototypeID:  int64(it.Prototype.ID),
			Status:       string(domain.ArmyStatusInProduction),
			PendingData:  pqtype.NullRawMessage{Valid: false},
			InProdData:   prodRaw,
			PresentData:  pqtype.NullRawMessage{Valid: false},
			DeployedData: pqtype.NullRawMessage{Valid: false},
			CreatedAt:    now,
		})
	}
	// Present
	for _, it := range base.ArmiesPresent {
		presentRaw := BuildArmyPresentRaw(it)
		out = append(out, gen.InsertBaseArmyItemParams{
			ID:           it.ID,
			BaseID:       int64(base.ID),
			PrototypeID:  int64(it.Prototype.ID),
			Status:       string(domain.ArmyStatusPresent),
			PendingData:  pqtype.NullRawMessage{Valid: false},
			InProdData:   pqtype.NullRawMessage{Valid: false},
			PresentData:  presentRaw,
			DeployedData: pqtype.NullRawMessage{Valid: false},
			CreatedAt:    now,
		})
	}
	// Deployed
	for _, it := range base.ArmiesDeployed {
		depRaw := BuildArmyDeployedRaw(it)
		out = append(out, gen.InsertBaseArmyItemParams{
			ID:           it.ID,
			BaseID:       int64(base.ID),
			PrototypeID:  int64(it.Prototype.ID),
			Status:       string(domain.ArmyStatusDeployed),
			PendingData:  pqtype.NullRawMessage{Valid: false},
			InProdData:   pqtype.NullRawMessage{Valid: false},
			PresentData:  pqtype.NullRawMessage{Valid: false},
			DeployedData: depRaw,
			CreatedAt:    now,
		})
	}
	return out
}
