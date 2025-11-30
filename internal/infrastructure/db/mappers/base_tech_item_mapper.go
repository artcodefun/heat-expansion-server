package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
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
