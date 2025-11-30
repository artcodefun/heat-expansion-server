package mappers

import (
	"database/sql"
	"encoding/json"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/readstore/gen"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// Prototype parts conversion
func buildPrototypeFromParts(id int64, name, category string, unlock sql.NullInt64, short, full sql.NullString, price []byte, productionTime int64, space int32, imageURL sql.NullString,
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
		UnlockTechnologyID: unlockPtr,
		ShortDescription:   nullString(short),
		FullDescription:    nullString(full),
		Price:              priceFromJSON(price),
		ProductionTime:     productionTime,
		Space:              int(space),
		ImageURL:           nullString(imageURL),
		ControlData:        jsonToNullRaw[readmodels.ControlBuildingData](control),
		ResourcesData:      jsonToNullRaw[readmodels.ResourcesBuildingData](resources),
		DefenseData:        jsonToNullRaw[readmodels.DefenseBuildingData](defense),
		MilitaryData:       jsonToNullRaw[readmodels.MilitaryBuildingData](military),
		IntelligenceData:   jsonToNullRaw[readmodels.IntelligenceBuildingData](intelligence),
	}
}

func NewBuildItemFromPrototype(p gen.BuildItemPrototype) readmodels.BuildItemNew {
	proto := buildPrototypeFromParts(p.ID, p.Name, p.Category, p.UnlockTechnologyID, p.ShortDescription, p.FullDescription, p.Price, p.ProductionTime, p.Space, p.ImageUrl, p.ControlData, p.ResourcesData, p.DefenseData, p.MilitaryData, p.IntelligenceData)
	return readmodels.BuildItemNew{Prototype: proto}
}

func BuildItemPendingFromRow(r gen.ListPendingBuildItemsRow) readmodels.BuildItemPending {
	return readmodels.BuildItemPending{BaseOwnedItem: readmodels.BaseOwnedItem{ID: uuid.UUID(r.ID), UserBaseID: int(r.BaseID)}, Prototype: buildPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ProductionTime, r.Space, r.ImageUrl, r.ControlData, r.ResourcesData, r.DefenseData, r.MilitaryData, r.IntelligenceData)}
}

func BuildItemInProductionFromRow(r gen.ListInProductionBuildItemsRow) readmodels.BuildItemInProduction {
	var jd dtos.BuildInProdDTO
	if r.InProdData.Valid {
		_ = json.Unmarshal(r.InProdData.RawMessage, &jd)
	}
	return readmodels.BuildItemInProduction{BaseOwnedItem: readmodels.BaseOwnedItem{ID: uuid.UUID(r.ID), UserBaseID: int(r.BaseID)}, Prototype: buildPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ProductionTime, r.Space, r.ImageUrl, r.ControlData, r.ResourcesData, r.DefenseData, r.MilitaryData, r.IntelligenceData), StartDate: jd.StartDate, CompletionDate: jd.CompletionDate, CrystalsSkipPrice: jd.CrystalsSkipPrice}
}

func BuildItemPresentFromRow(r gen.ListPresentBuildItemsRow) readmodels.BuildItemPresent {
	var jd dtos.BuildPresentDTO
	if r.PresentData.Valid {
		_ = json.Unmarshal(r.PresentData.RawMessage, &jd)
	}
	refund := readmodels.PriceModel{Credits: jd.Refund.Credits, Iron: jd.Refund.Iron, Titanium: jd.Refund.Titanium, Antimatter: jd.Refund.Antimatter}
	return readmodels.BuildItemPresent{BaseOwnedItem: readmodels.BaseOwnedItem{ID: uuid.UUID(r.ID), UserBaseID: int(r.BaseID)}, Prototype: buildPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ProductionTime, r.Space, r.ImageUrl, r.ControlData, r.ResourcesData, r.DefenseData, r.MilitaryData, r.IntelligenceData), Refund: refund}
}
