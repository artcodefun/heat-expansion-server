package mappers

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/infrastructure/readstore/gen"
)

// UserBaseFromBasicRow maps basic base fields to readmodels.UserBaseModel.
func UserBaseFromBasicRow(r gen.ListUserBasesRow) readmodels.UserBaseModel {
	return readmodels.UserBaseModel{
		ID:     int(r.ID),
		UserID: int(r.UserID),
		Coordinates: readmodels.Vector2i{
			X: int(r.SectorX),
			Y: int(r.SectorY),
		},
		LocationDetails: readmodels.LocationDetails{
			Name:        nullString(r.Name),
			Description: nullString(r.Description),
			ImageURL:    nullString(r.ImageUrl),
		},
	}
}
