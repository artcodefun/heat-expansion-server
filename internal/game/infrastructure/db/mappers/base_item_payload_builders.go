package mappers

import (
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/sqlc-dev/pqtype"
)

// Army payloads
func BuildArmyPendingRaw(it domain.ArmyItemPending) pqtype.NullRawMessage {
	dto := dtos.ArmyPendingDTO{Count: it.Count}
	b, _ := json.Marshal(dto)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}

func BuildArmyInProdRaw(it domain.ArmyItemInProduction) pqtype.NullRawMessage {
	dto := dtos.ArmyInProdDTO{StartDate: it.StartDate, CompletionDate: it.CompletionDate, CrystalsSkipPrice: it.CrystalsSkipPrice}
	b, _ := json.Marshal(dto)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}

func BuildArmyPresentRaw(it domain.ArmyItemPresent) pqtype.NullRawMessage {
	dto := dtos.ArmyPresentDTO{Count: it.Count, Refund: dtos.PriceDTOFromDomain(it.Refund)}
	b, _ := json.Marshal(dto)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}

func BuildArmyDeployedRaw(it domain.ArmyItemDeployed) pqtype.NullRawMessage {
	dto := dtos.ArmyDeployedDTOFromDomain(it)
	b, _ := json.Marshal(dto)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}

// Build payloads
func BuildBuildPendingRaw(it domain.BuildItemPending) pqtype.NullRawMessage {
	dto := dtos.BuildPendingDTOFromDomain(it)
	b, _ := json.Marshal(dto)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}

func BuildBuildInProdRaw(it domain.BuildItemInProduction) pqtype.NullRawMessage {
	dto := dtos.BuildInProdDTO{StartDate: it.StartDate, CompletionDate: it.CompletionDate, CrystalsSkipPrice: it.CrystalsSkipPrice}
	b, _ := json.Marshal(dto)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}

func BuildBuildPresentRaw(it domain.BuildItemPresent) pqtype.NullRawMessage {
	dto := dtos.BuildPresentDTO{Refund: dtos.PriceDTOFromDomain(it.Refund)}
	b, _ := json.Marshal(dto)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}

// Tech payloads
func BuildTechInProgressRaw(it domain.TechItemInProgress) pqtype.NullRawMessage {
	dto := dtos.TechInProgressDTO{StartDate: it.StartDate, CompletionDate: it.CompletionDate, CrystalsSkipPrice: it.CrystalsSkipPrice}
	b, _ := json.Marshal(dto)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}

func BuildTechDoneRaw(it domain.TechItemDone) pqtype.NullRawMessage {
	dto := dtos.TechDoneDTO{ResearchedAt: it.ResearchedAt, Level: it.Level}
	b, _ := json.Marshal(dto)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}

// Storage payloads
func BuildStoragePresentRaw(it domain.StorageItemPresent) pqtype.NullRawMessage {
	dto := dtos.StoragePresentDTOFromDomain(it)
	b, _ := json.Marshal(dto)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}

func BuildStorageDeployedRaw(it domain.StorageItemDeployed) pqtype.NullRawMessage {
	dto := dtos.StorageDeployedDTOFromDomain(it)
	b, _ := json.Marshal(dto)
	return pqtype.NullRawMessage{RawMessage: b, Valid: true}
}
