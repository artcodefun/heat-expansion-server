package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/sqlc-dev/pqtype"
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
			var d dtos.BuildPendingDTO
			unmarshalIfValid(r.PendingData, &d)
			base.BuildingsPending = append(base.BuildingsPending, dtos.BuildPendingFromDTO(d, owned, *p))
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

// DehydrateBuildItems converts the in-memory aggregate collections into insert params
// for the base_build_items table.
func DehydrateBuildItems(base *domain.UserBaseModel) []gen.InsertBaseBuildItemParams {
	now := domain.NowUnix()
	out := make([]gen.InsertBaseBuildItemParams, 0,
		len(base.BuildingsPending)+len(base.BuildingsInProduction)+len(base.BuildingsPresent))

	// Pending (store empty JSON object to satisfy constraint)
	for _, it := range base.BuildingsPending {
		pendingRaw := BuildBuildPendingRaw(it)
		out = append(out, gen.InsertBaseBuildItemParams{
			ID:          it.ID,
			BaseID:      int64(base.ID),
			PrototypeID: int64(it.Prototype.ID),
			Status:      string(domain.BuildStatusPending),
			PendingData: pendingRaw,
			InProdData:  pqtype.NullRawMessage{Valid: false},
			PresentData: pqtype.NullRawMessage{Valid: false},
			CreatedAt:   now,
		})
	}
	// In Production
	for _, it := range base.BuildingsInProduction {
		prodRaw := BuildBuildInProdRaw(it)
		out = append(out, gen.InsertBaseBuildItemParams{
			ID:          it.ID,
			BaseID:      int64(base.ID),
			PrototypeID: int64(it.Prototype.ID),
			Status:      string(domain.BuildStatusInProduction),
			PendingData: pqtype.NullRawMessage{Valid: false},
			InProdData:  prodRaw,
			PresentData: pqtype.NullRawMessage{Valid: false},
			CreatedAt:   now,
		})
	}
	// Present
	for _, it := range base.BuildingsPresent {
		presentRaw := BuildBuildPresentRaw(it)
		out = append(out, gen.InsertBaseBuildItemParams{
			ID:          it.ID,
			BaseID:      int64(base.ID),
			PrototypeID: int64(it.Prototype.ID),
			Status:      string(domain.BuildStatusPresent),
			PendingData: pqtype.NullRawMessage{Valid: false},
			InProdData:  pqtype.NullRawMessage{Valid: false},
			PresentData: presentRaw,
			CreatedAt:   now,
		})
	}
	return out
}
