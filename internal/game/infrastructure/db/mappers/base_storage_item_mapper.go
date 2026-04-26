package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/gen"
	"github.com/sqlc-dev/pqtype"
)

func HydrateStorageItems(base *domain.UserBaseModel, rows []gen.ListBaseStorageItemsRow, proto map[int]*domain.StorageItemPrototype) {
	for _, r := range rows {
		p := proto[int(r.PrototypeID)]
		if p == nil {
			continue
		}
		owned := domain.BaseOwnedItem{ID: r.ID, UserBaseID: base.ID}
		switch r.Status {
		case string(domain.StorageStatusPresent):
			var d dtos.StoragePresentDTO
			unmarshalIfValid(r.PresentData, &d)
			base.StorageItemsPresent = append(base.StorageItemsPresent, dtos.StoragePresentFromDTO(d, owned, *p))
		case string(domain.StorageStatusDeployed):
			var d dtos.StorageDeployedDTO
			unmarshalIfValid(r.DeployedData, &d)
			base.StorageItemsDeployed = append(base.StorageItemsDeployed, dtos.StorageDeployedFromDTO(d, owned, *p))
		}
	}
}

// DehydrateStorageItems converts present storage items into insert params
// for the base_storage_items table.
func DehydrateStorageItems(base *domain.UserBaseModel) []gen.InsertBaseStorageItemParams {
	now := domain.NowUnix()
	out := make([]gen.InsertBaseStorageItemParams, 0, len(base.StorageItemsPresent)+len(base.StorageItemsDeployed))

	for _, it := range base.StorageItemsPresent {
		presentRaw := BuildStoragePresentRaw(it)
		stateJSON, _ := json.Marshal(map[string]any{}) // placeholder empty state
		out = append(out, gen.InsertBaseStorageItemParams{
			ID:           it.ID,
			BaseID:       int64(base.ID),
			PrototypeID:  int64(it.Prototype.ID),
			Status:       string(domain.StorageStatusPresent),
			PresentData:  presentRaw,
			DeployedData: pqNullRawInvalid(),
			State:        stateJSON,
			CreatedAt:    now,
		})
	}

	for _, it := range base.StorageItemsDeployed {
		deployedRaw := BuildStorageDeployedRaw(it)
		stateJSON, _ := json.Marshal(map[string]any{})
		out = append(out, gen.InsertBaseStorageItemParams{
			ID:           it.ID,
			BaseID:       int64(base.ID),
			PrototypeID:  int64(it.Prototype.ID),
			Status:       string(domain.StorageStatusDeployed),
			PresentData:  pqNullRawInvalid(),
			DeployedData: deployedRaw,
			State:        stateJSON,
			CreatedAt:    now,
		})
	}
	return out
}

func pqNullRawInvalid() pqtype.NullRawMessage {
	return pqtype.NullRawMessage{Valid: false}
}
