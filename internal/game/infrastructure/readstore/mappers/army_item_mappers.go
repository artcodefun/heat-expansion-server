package mappers

import (
	"database/sql"
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
)

// Convert prototype row to readmodel prototype
func armyPrototypeFromParts(id int64, name string, category string, faction string, unlock sql.NullInt64, short sql.NullString, full sql.NullString, price []byte, productionTime int64, space int32, imageURL sql.NullString, attack, defence, capacity, stealth, speed int32) readmodels.ArmyItemPrototype {
	var unlockPtr *int
	if unlock.Valid {
		v := int(unlock.Int64)
		unlockPtr = &v
	}
	return readmodels.ArmyItemPrototype{
		ID:                 int(id),
		Name:               name,
		Category:           readmodels.ArmyCategory(category),
		Faction:            readmodels.Faction(faction),
		UnlockTechnologyID: unlockPtr,
		ShortDescription:   nullString(short),
		FullDescription:    nullString(full),
		Price:              priceFromJSON(price),
		ProductionTime:     productionTime,
		Space:              int(space),
		ImageURL:           nullString(imageURL),
		Attack:             int(attack),
		Defence:            int(defence),
		Capacity:           int(capacity),
		Stealth:            int(stealth),
		Speed:              int(speed),
	}
}

func ArmyItemPendingFromRow(r gen.ListPendingArmyItemsRow) readmodels.ArmyItemPending {
	count := 0
	if r.PendingData.Valid {
		var tmp dtos.ArmyPendingDTO
		_ = json.Unmarshal(r.PendingData.RawMessage, &tmp)
		count = tmp.Count
	}
	return readmodels.ArmyItemPending{
		BaseOwnedItem: readmodels.BaseOwnedItem{ID: r.ID, UserBaseID: int(r.BaseID)},
		Prototype:     armyPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.Faction, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ProductionTime, r.Space, r.ImageUrl, r.Attack, r.Defence, r.Capacity, r.Stealth, r.Speed),
		Count:         count,
	}
}

func ArmyItemInProductionFromRow(r gen.ListInProductionArmyItemsRow) readmodels.ArmyItemInProduction {
	var jd dtos.ArmyInProdDTO
	if r.InProdData.Valid {
		_ = json.Unmarshal(r.InProdData.RawMessage, &jd)
	}
	return readmodels.ArmyItemInProduction{
		BaseOwnedItem:     readmodels.BaseOwnedItem{ID: r.ID, UserBaseID: int(r.BaseID)},
		Prototype:         armyPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.Faction, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ProductionTime, r.Space, r.ImageUrl, r.Attack, r.Defence, r.Capacity, r.Stealth, r.Speed),
		StartDate:         jd.StartDate,
		CompletionDate:    jd.CompletionDate,
		CrystalsSkipPrice: jd.CrystalsSkipPrice,
	}
}

func ArmyItemPresentFromRow(r gen.ListPresentArmyItemsRow) readmodels.ArmyItemPresent {
	var jd dtos.ArmyPresentDTO
	if r.PresentData.Valid {
		_ = json.Unmarshal(r.PresentData.RawMessage, &jd)
	}
	refund := readmodels.PriceModel{Credits: jd.Refund.Credits, Iron: jd.Refund.Iron, Titanium: jd.Refund.Titanium, Antimatter: jd.Refund.Antimatter}
	return readmodels.ArmyItemPresent{
		BaseOwnedItem: readmodels.BaseOwnedItem{ID: r.ID, UserBaseID: int(r.BaseID)},
		Prototype:     armyPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.Faction, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ProductionTime, r.Space, r.ImageUrl, r.Attack, r.Defence, r.Capacity, r.Stealth, r.Speed),
		Count:         jd.Count,
		Refund:        refund,
	}
}

func NewArmyItemFromPrototype(p gen.ArmyItemPrototype) readmodels.ArmyItemNew {
	proto := ArmyPrototypeFromModel(p)
	return readmodels.ArmyItemNew{Prototype: proto}
}

func ArmyPrototypeFromModel(p gen.ArmyItemPrototype) readmodels.ArmyItemPrototype {
	return armyPrototypeFromParts(p.ID, p.Name, p.Category, p.Faction, p.UnlockTechnologyID, p.ShortDescription, p.FullDescription, p.Price, p.ProductionTime, p.Space, p.ImageUrl, p.Attack, p.Defence, p.Capacity, p.Stealth, p.Speed)
}
