package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
	"github.com/sqlc-dev/pqtype"
)

func HydrateTechItems(base *domain.UserBaseModel, rows []gen.BaseTechItem, proto map[int]*domain.TechItemPrototype) {
	for _, r := range rows {
		p := proto[int(r.PrototypeID)]
		if p == nil {
			continue
		}
		owned := domain.BaseOwnedItem{ID: r.ID, UserBaseID: base.ID}
		switch r.Status {
		case string(domain.TechStatusInProgress):
			var d dtos.TechInProgressDTO
			unmarshalIfValid(r.InProgressData, &d)
			base.TechnologiesInProgress = append(base.TechnologiesInProgress, dtos.TechInProgressFromDTO(d, owned, *p))
		case string(domain.TechStatusDone):
			var d dtos.TechDoneDTO
			unmarshalIfValid(r.DoneData, &d)
			base.TechnologiesDone = append(base.TechnologiesDone, dtos.TechDoneFromDTO(d, owned, *p))
		}
	}
}

// DehydrateTechItems converts the in-memory tech collections into insert params
// for the base_tech_items table.
func DehydrateTechItems(base *domain.UserBaseModel) []gen.InsertBaseTechItemParams {
	now := domain.NowUnix()
	out := make([]gen.InsertBaseTechItemParams, 0,
		len(base.TechnologiesInProgress)+len(base.TechnologiesDone))

	// In Progress
	for _, it := range base.TechnologiesInProgress {
		inProgRaw := BuildTechInProgressRaw(it)
		out = append(out, gen.InsertBaseTechItemParams{
			ID:             it.ID,
			BaseID:         int64(base.ID),
			PrototypeID:    int64(it.Prototype.ID),
			Status:         string(domain.TechStatusInProgress),
			InProgressData: inProgRaw,
			DoneData:       pqtype.NullRawMessage{Valid: false},
			CreatedAt:      now,
		})
	}
	// Done
	for _, it := range base.TechnologiesDone {
		doneRaw := BuildTechDoneRaw(it)
		out = append(out, gen.InsertBaseTechItemParams{
			ID:             it.ID,
			BaseID:         int64(base.ID),
			PrototypeID:    int64(it.Prototype.ID),
			Status:         string(domain.TechStatusDone),
			InProgressData: pqtype.NullRawMessage{Valid: false},
			DoneData:       doneRaw,
			CreatedAt:      now,
		})
	}
	return out
}
