package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/dtos"
)

func UserBaseStatsFromDTO(dto dtos.BaseStatsDTO, calcTs int64) readmodels.UserBaseStats {
	return readmodels.UserBaseStats{
		Credits:               int(dto.Credits),
		CreditsCapacity:       dto.CreditsCapacity,
		CreditsProduction:     dto.CreditsProduction,
		Iron:                  int(dto.Iron),
		IronCapacity:          dto.IronCapacity,
		IronProduction:        dto.IronProduction,
		Titanium:              int(dto.Titanium),
		TitaniumCapacity:      dto.TitaniumCapacity,
		TitaniumProduction:    dto.TitaniumProduction,
		Antimatter:            int(dto.Antimatter),
		AntimatterCapacity:    dto.AntimatterCapacity,
		AntimatterProduction:  dto.AntimatterProduction,
		Defence:               dto.Defence,
		Attack:                dto.Attack,
		Space:                 dto.Space,
		MaxSpace:              dto.MaxSpace,
		MaxOperations:         dto.MaxOperations,
		MaxActiveBuffs:        dto.MaxActiveBuffs,
		MaxActiveArtifacts:    dto.MaxActiveArtifacts,
		MaxBuildingProduction: dto.MaxBuildingProduction,
		MaxActiveRestorations: dto.MaxActiveRestorations,
		MaxActiveDecryptions:  dto.MaxActiveDecryptions,
		CalculationTimestamp:  calcTs,
	}
}
