package mappers

import (
	"database/sql"
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
	"github.com/sqlc-dev/pqtype"
)

// Prototype parts conversion
func buildPrototypeFromParts(id int64, name, category, faction string, unlock sql.NullInt64, short, full sql.NullString, price []byte, productionTime int64, space int32, imageURL sql.NullString,
	control, resources, defense, military, intelligence pqtype.NullRawMessage) readmodels.BuildItemPrototype {
	var unlockPtr *int
	if unlock.Valid {
		v := int(unlock.Int64)
		unlockPtr = &v
	}
	return readmodels.BuildItemPrototype{
		ID:                 int(id),
		Name:               name,
		Category:           readmodels.BuildCategory(category),
		Faction:            readmodels.Faction(faction),
		UnlockTechnologyID: unlockPtr,
		ShortDescription:   nullString(short),
		FullDescription:    nullString(full),
		Price:              priceFromJSON(price),
		ProductionTime:     productionTime,
		Space:              int(space),
		ImageURL:           nullString(imageURL),
		ControlData:        controlBuildingDataFromJSON(control),
		ResourcesData:      resourcesBuildingDataFromJSON(resources),
		DefenseData:        defenseBuildingDataFromJSON(defense),
		MilitaryData:       militaryBuildingDataFromJSON(military),
		IntelligenceData:   intelligenceBuildingDataFromJSON(intelligence),
	}
}

func NewBuildItemFromPrototype(p gen.BuildItemPrototype) readmodels.BuildItemNew {
	proto := BuildPrototypeFromModel(p)
	return readmodels.BuildItemNew{Prototype: proto}
}

func BuildPrototypeFromModel(p gen.BuildItemPrototype) readmodels.BuildItemPrototype {
	proto := buildPrototypeFromParts(p.ID, p.Name, p.Category, p.Faction, p.UnlockTechnologyID, p.ShortDescription, p.FullDescription, p.Price, p.ProductionTime, p.Space, p.ImageUrl, p.ControlData, p.ResourcesData, p.DefenseData, p.MilitaryData, p.IntelligenceData)
	proto.CreationSources = creationSourcesFromJSON(p.CreationSources)
	return proto
}

func BuildPrototypesFromModels(rows []gen.BuildItemPrototype) []*readmodels.BuildItemPrototype {
	out := make([]*readmodels.BuildItemPrototype, len(rows))
	for i, p := range rows {
		v := BuildPrototypeFromModel(p)
		out[i] = &v
	}
	return out
}

func BuildItemPendingFromRow(r gen.ListPendingBuildItemsRow) readmodels.BuildItemPending {
	return readmodels.BuildItemPending{BaseOwnedItem: readmodels.BaseOwnedItem{ID: r.ID, UserBaseID: int(r.BaseID)}, Prototype: buildPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.Faction, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ProductionTime, r.Space, r.ImageUrl, r.ControlData, r.ResourcesData, r.DefenseData, r.MilitaryData, r.IntelligenceData)}
}

func BuildItemInProductionFromRow(r gen.ListInProductionBuildItemsRow) readmodels.BuildItemInProduction {
	var jd dtos.BuildInProdDTO
	if r.InProdData.Valid {
		_ = json.Unmarshal(r.InProdData.RawMessage, &jd)
	}
	return readmodels.BuildItemInProduction{BaseOwnedItem: readmodels.BaseOwnedItem{ID: r.ID, UserBaseID: int(r.BaseID)}, Prototype: buildPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.Faction, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ProductionTime, r.Space, r.ImageUrl, r.ControlData, r.ResourcesData, r.DefenseData, r.MilitaryData, r.IntelligenceData), StartDate: jd.StartDate, CompletionDate: jd.CompletionDate, CrystalsSkipPrice: jd.CrystalsSkipPrice}
}

func BuildItemPresentFromRow(r gen.ListPresentBuildItemsRow) readmodels.BuildItemPresent {
	var jd dtos.BuildPresentDTO
	if r.PresentData.Valid {
		_ = json.Unmarshal(r.PresentData.RawMessage, &jd)
	}
	refund := readmodels.PriceModel{Credits: jd.Refund.Credits, Iron: jd.Refund.Iron, Titanium: jd.Refund.Titanium, Antimatter: jd.Refund.Antimatter}
	return readmodels.BuildItemPresent{BaseOwnedItem: readmodels.BaseOwnedItem{ID: r.ID, UserBaseID: int(r.BaseID)}, Prototype: buildPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.Faction, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ProductionTime, r.Space, r.ImageUrl, r.ControlData, r.ResourcesData, r.DefenseData, r.MilitaryData, r.IntelligenceData), Refund: refund}
}

// Build prototype detail helpers: JSONB (DTO shape) -> readmodels.* data

func controlBuildingDataFromJSON(nm pqtype.NullRawMessage) *readmodels.ControlBuildingData {
	if !nm.Valid {
		return nil
	}
	var d dtos.ControlBuildingDataDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	return &readmodels.ControlBuildingData{Subtype: readmodels.ControlSubtype(d.Subtype)}
}

func resourcesBuildingDataFromJSON(nm pqtype.NullRawMessage) *readmodels.ResourcesBuildingData {
	if !nm.Valid {
		return nil
	}
	var d dtos.ResourcesBuildingDataDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	return &readmodels.ResourcesBuildingData{
		CreditsProduction:    d.CreditsProduction,
		IronProduction:       d.IronProduction,
		TitaniumProduction:   d.TitaniumProduction,
		AntimatterProduction: d.AntimatterProduction,
		CreditsCapacity:      d.CreditsCapacity,
		IronCapacity:         d.IronCapacity,
		TitaniumCapacity:     d.TitaniumCapacity,
		AntimatterCapacity:   d.AntimatterCapacity,
	}
}

func defenseBuildingDataFromJSON(nm pqtype.NullRawMessage) *readmodels.DefenseBuildingData {
	if !nm.Valid {
		return nil
	}
	var d dtos.DefenseBuildingDataDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	return &readmodels.DefenseBuildingData{
		DefenceBonus: d.DefenceBonus,
	}
}

func militaryBuildingDataFromJSON(nm pqtype.NullRawMessage) *readmodels.MilitaryBuildingData {
	if !nm.Valid {
		return nil
	}
	var d dtos.MilitaryBuildingDataDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	return &readmodels.MilitaryBuildingData{UnlockArmyCategory: readmodels.ArmyCategory(d.UnlockArmyCategory)}
}

func intelligenceBuildingDataFromJSON(nm pqtype.NullRawMessage) *readmodels.IntelligenceBuildingData {
	if !nm.Valid {
		return nil
	}
	var d dtos.IntelligenceBuildingDataDTO
	if err := json.Unmarshal(nm.RawMessage, &d); err != nil {
		return nil
	}
	return &readmodels.IntelligenceBuildingData{
		Subtype:         readmodels.IntelligenceSubtype(d.Subtype),
		StealthStrength: d.StealthStrength,
		ScanRange:       d.ScanRange,
		ScanCooldown:    d.ScanCooldown,
	}
}
