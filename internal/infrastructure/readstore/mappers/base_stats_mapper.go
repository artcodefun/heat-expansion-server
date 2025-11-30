package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
)

func UserBaseStatsFromDTO(dto dtos.BaseStatsDTO, calcTs int64) readmodels.UserBaseStats {
	return readmodels.UserBaseStats{
		Credits:              dto.Credits,
		CreditsCapacity:      dto.CreditsCapacity,
		CreditsProduction:    dto.CreditsProduction,
		Iron:                 dto.Iron,
		IronCapacity:         dto.IronCapacity,
		IronProduction:       dto.IronProduction,
		Titanium:             dto.Titanium,
		TitaniumCapacity:     dto.TitaniumCapacity,
		TitaniumProduction:   dto.TitaniumProduction,
		Antimatter:           dto.Antimatter,
		AntimatterCapacity:   dto.AntimatterCapacity,
		AntimatterProduction: dto.AntimatterProduction,
		Defence:              dto.Defence,
		Attack:               dto.Attack,
		Space:                dto.Space,
		SpaceCapacity:        dto.SpaceCapacity,
		CalculationTimestamp: calcTs,
	}
}
