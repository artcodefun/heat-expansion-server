package mappers

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/infrastructure/db/gen"
)

func SectorFromDB(row gen.Sector) *domain.SectorModel {
	return &domain.SectorModel{
		Coordinates: domain.Vector2i{X: int(row.X), Y: int(row.Y)},
		Details: domain.LocationDetails{
			Name:        nullStringToString(&row.Name.String, row.Name.Valid),
			Description: nullStringToString(&row.Description.String, row.Description.Valid),
			ImageURL:    nullStringToString(&row.ImageUrl.String, row.ImageUrl.Valid),
		},
	}
}

// Insert params builder (CreateSectorParams already matches field order except ID which is returned)
func InsertSectorParamsFromDomain(s *domain.SectorModel) gen.CreateSectorParams {
	return gen.CreateSectorParams{
		X:           int32(s.Coordinates.X),
		Y:           int32(s.Coordinates.Y),
		Name:        toNullString(s.Details.Name),
		Description: toNullString(s.Details.Description),
		ImageUrl:    toNullString(s.Details.ImageURL),
	}
}

func UpdateSectorParamsFromDomain(s *domain.SectorModel) gen.UpdateSectorParams {
	return gen.UpdateSectorParams{
		Name:        toNullString(s.Details.Name),
		Description: toNullString(s.Details.Description),
		ImageUrl:    toNullString(s.Details.ImageURL),
		X:           int32(s.Coordinates.X),
		Y:           int32(s.Coordinates.Y),
	}
}
