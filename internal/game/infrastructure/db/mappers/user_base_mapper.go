package mappers

import (
	"database/sql"
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/game/domain"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/game/infrastructure/db/gen"
)

func UserBaseFromDB(b gen.UserBase) *domain.UserBaseModel {
	// Unmarshal stats JSON
	var dto dtos.BaseStatsDTO
	_ = json.Unmarshal(b.Stats, &dto)

	ub := &domain.UserBaseModel{
		ID:          int(b.ID),
		Coordinates: domain.Vector2i{X: int(b.SectorX), Y: int(b.SectorY)},
		UserID:      int(b.UserID),
		Stats:       dtos.BaseStatsFromDTO(dto, b.StatsCalcTimestamp),
	}
	// Location details
	ub.Name = nullStringToString(&b.Name.String, b.Name.Valid)
	ub.Description = nullStringToString(&b.Description.String, b.Description.Valid)
	ub.ImageURL = nullStringToString(&b.ImageUrl.String, b.ImageUrl.Valid)

	// Item collections are populated by the repository as needed. Default to empty slices.
	ub.ArmiesPending = []domain.ArmyItemPending{}
	ub.ArmiesPresent = []domain.ArmyItemPresent{}
	ub.ArmiesInProduction = []domain.ArmyItemInProduction{}
	ub.ArmiesDeployed = []domain.ArmyItemDeployed{}
	ub.BuildingsPending = []domain.BuildItemPending{}
	ub.BuildingsPresent = []domain.BuildItemPresent{}
	ub.BuildingsInProduction = []domain.BuildItemInProduction{}
	ub.TechnologiesInProgress = []domain.TechItemInProgress{}
	ub.TechnologiesDone = []domain.TechItemDone{}
	ub.StorageItemsPresent = []domain.StorageItemPresent{}
	return ub
}

func InsertBaseParamsFromDomain(base *domain.UserBaseModel) gen.CreateBaseParams {
	dto := dtos.BaseStatsDTOFromDomain(base.Stats)
	stats, _ := json.Marshal(dto)
	return gen.CreateBaseParams{
		UserID:             int64(base.UserID),
		SectorX:            int32(base.Coordinates.X),
		SectorY:            int32(base.Coordinates.Y),
		Name:               toNullString(base.Name),
		Description:        toNullString(base.Description),
		ImageUrl:           toNullString(base.ImageURL),
		Stats:              stats,
		StatsCalcTimestamp: base.Stats.CalculationTimestamp,
	}
}

func UpdateBaseParamsFromDomain(base *domain.UserBaseModel) gen.UpdateBaseParams {
	dto := dtos.BaseStatsDTOFromDomain(base.Stats)
	stats, _ := json.Marshal(dto)
	return gen.UpdateBaseParams{
		ID:                 int64(base.ID),
		Name:               toNullString(base.Name),
		Description:        toNullString(base.Description),
		ImageUrl:           toNullString(base.ImageURL),
		Stats:              stats,
		StatsCalcTimestamp: base.Stats.CalculationTimestamp,
	}
}

// Helper to construct sql.NullString from plain string (empty -> NULL)
func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}
