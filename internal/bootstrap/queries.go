package bootstrap

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/queries"
	"github.com/artcodefun/heat-expansion-api/internal/core/services"
)

// Queries aggregates all query facades.
type Queries struct {
	Base      *queries.BaseQueries
	Army      *queries.ArmyQueries
	Building  *queries.BuildingQueries
	Tech      *queries.TechQueries
	Storage   *queries.StorageQueries
	Sector    *queries.SectorQueries
	Operation *queries.OperationQueries
	Activity  *queries.ActivityQueries
	User      *queries.UserQueries
}

// NewQueries builds query facades using read repositories and shared services.
func NewQueries(a *Adapters) *Queries {
	access := services.NewAccessControlService(a.UserBases)
	return &Queries{
		Base:      queries.NewBaseQueries(a.BaseRead, access),
		Army:      queries.NewArmyQueries(a.ArmyRead, access),
		Building:  queries.NewBuildingQueries(a.BuildingRead, access),
		Tech:      queries.NewTechQueries(a.TechRead, access),
		Storage:   queries.NewStorageQueries(a.StorageRead, access),
		Sector:    queries.NewSectorQueries(a.SectorRead, access),
		Operation: queries.NewOperationQueries(a.OperationRead, access),
		Activity:  queries.NewActivityQueries(a.ActivityRead, access),
		User:      queries.NewUserQueries(a.UserRead),
	}
}
