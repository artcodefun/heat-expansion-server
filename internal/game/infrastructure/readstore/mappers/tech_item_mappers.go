package mappers

import (
	"database/sql"
	"encoding/json"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
)

func techPrototypeFromParts(id int64, name, category string, unlock sql.NullInt64, short, full sql.NullString, price []byte, researchTime int64, imageURL sql.NullString, improvement []byte) readmodels.TechItemPrototype {
	var unlockPtr *int
	if unlock.Valid {
		v := int(unlock.Int64)
		unlockPtr = &v
	}
	return readmodels.TechItemPrototype{
		ID:                 int(id),
		Name:               name,
		Category:           readmodels.TechCategory(category),
		UnlockTechnologyID: unlockPtr,
		ShortDescription:   nullString(short),
		FullDescription:    nullString(full),
		Price:              priceFromJSON(price),
		ResearchTime:       researchTime,
		ImageURL:           nullString(imageURL),
		Improvement:        techImprovementFromJSON(improvement),
	}
}

func NewTechItemFromPrototype(p gen.TechItemPrototype, currentLevel int) readmodels.TechItemNew {
	proto := techPrototypeFromParts(p.ID, p.Name, p.Category, p.UnlockTechnologyID, p.ShortDescription, p.FullDescription, p.Price, p.ResearchTime, p.ImageUrl, p.Improvement.RawMessage)
	multiplier := 1.0 + float64(currentLevel)*0.5
	return readmodels.TechItemNew{
		Prototype:           proto,
		CurrentLevel:        currentLevel,
		CurrentPrice:        proto.Price.MultiplyFloat(multiplier),
		CurrentResearchTime: int64(float64(proto.ResearchTime) * multiplier),
	}
}

func TechItemInProgressFromRow(r gen.ListInResearchTechItemsRow, currentLevel int) readmodels.TechItemInProgress {
	var jd dtos.TechInProgressDTO
	if r.InProgressData.Valid {
		_ = json.Unmarshal(r.InProgressData.RawMessage, &jd)
	}
	return readmodels.TechItemInProgress{
		BaseOwnedItem:     readmodels.BaseOwnedItem{ID: r.ID, UserBaseID: int(r.BaseID)},
		Prototype:         techPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ResearchTime, r.ImageUrl, r.Improvement.RawMessage),
		StartDate:         jd.StartDate,
		CompletionDate:    jd.CompletionDate,
		CrystalsSkipPrice: jd.CrystalsSkipPrice,
		CurrentLevel:      currentLevel,
	}
}

func TechItemDoneFromRow(r gen.ListDoneTechItemsRow) readmodels.TechItemDone {
	var jd dtos.TechDoneDTO
	if r.DoneData.Valid {
		_ = json.Unmarshal(r.DoneData.RawMessage, &jd)
	}
	return readmodels.TechItemDone{
		BaseOwnedItem: readmodels.BaseOwnedItem{ID: r.ID, UserBaseID: int(r.BaseID)},
		Prototype:     techPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ResearchTime, r.ImageUrl, r.Improvement.RawMessage),
		ResearchedAt:  jd.ResearchedAt,
		Level:         jd.Level,
	}
}

func TechItemDoneFromAllRow(r gen.ListDoneTechItemsAllRow) readmodels.TechItemDone {
	var jd dtos.TechDoneDTO
	if r.DoneData.Valid {
		_ = json.Unmarshal(r.DoneData.RawMessage, &jd)
	}
	return readmodels.TechItemDone{
		BaseOwnedItem: readmodels.BaseOwnedItem{ID: r.ID, UserBaseID: int(r.BaseID)},
		Prototype:     techPrototypeFromParts(r.ProtoID, r.Name, r.Category, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ResearchTime, r.ImageUrl, r.Improvement.RawMessage),
		ResearchedAt:  jd.ResearchedAt,
		Level:         jd.Level,
	}
}

func TechPrototypeFromModel(r gen.TechItemPrototype) readmodels.TechItemPrototype {
	return techPrototypeFromParts(r.ID, r.Name, r.Category, r.UnlockTechnologyID, r.ShortDescription, r.FullDescription, r.Price, r.ResearchTime, r.ImageUrl, r.Improvement.RawMessage)
}

func TechPrototypesFromModels(rows []gen.TechItemPrototype) []*readmodels.TechItemPrototype {
	dst := make([]*readmodels.TechItemPrototype, len(rows))
	for i, r := range rows {
		v := TechPrototypeFromModel(r)
		dst[i] = &v
	}
	return dst
}

func techImprovementFromJSON(b []byte) *readmodels.TechImprovement {
	if len(b) == 0 {
		return nil
	}
	var d dtos.TechImprovementDTO
	if err := json.Unmarshal(b, &d); err != nil {
		return nil
	}
	return &readmodels.TechImprovement{
		Type:     readmodels.ImprovementType(d.Type),
		Value:    d.Value,
		MaxLevel: d.MaxLevel,
	}
}
