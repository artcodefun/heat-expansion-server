package bootstrap

import "github.com/artcodefun/heat-expansion-api/internal/game/application/queries"

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
	Radar     *queries.RadarQueries
	User      *queries.UserQueries
	Alert     *queries.AlertQueries
}

// NewQueries builds query facades using read repositories and shared services.

func NewQueries(a *Adapters, as *AppServices) *Queries {
	return &Queries{
		Base:      queries.NewBaseQueries(a.BaseRead, as.Access),
		Army:      queries.NewArmyQueries(a.ArmyRead, a.ArmyPrototypes, a.UserBases, as.Access),
		Building:  queries.NewBuildingQueries(a.BuildingRead, a.BuildPrototypes, a.UserBases, as.Access),
		Tech:      queries.NewTechQueries(a.TechRead, a.TechPrototypes, a.UserBases, as.Access),
		Storage:   queries.NewStorageQueries(a.StorageRead, as.Access),
		Sector:    queries.NewSectorQueries(a.SectorRead, as.Access),
		Operation: queries.NewOperationQueries(a.OperationRead, as.Access),
		Activity:  queries.NewActivityQueries(a.ActivityRead, as.Access),
		Radar:     queries.NewRadarQueries(a.RadarRead, as.Access),
		User:      queries.NewUserQueries(a.UserRead),
		Alert:     queries.NewAlertQueries(a.AlertRead, as.Access)}
}
