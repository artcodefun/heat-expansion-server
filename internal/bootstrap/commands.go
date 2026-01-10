package bootstrap

import "github.com/artcodefun/heat-expansion-api/internal/core/commands"

// Commands aggregates all command handlers.
type Commands struct {
	Base      *commands.BaseCommands
	Army      *commands.ArmyCommands
	Building  *commands.BuildingCommands
	Tech      *commands.TechCommands
	Storage   *commands.StorageCommands
	Operation *commands.OperationCommands
	Scanner   *commands.IntelligenceScannerCommands
	Radar     *commands.IntelligenceRadarCommands
	User      *commands.UserCommands
	Activity  *commands.ActivityCommands
	World     *commands.WorldGenerationCommands
}

// NewCommands constructs all command handlers using provided secondary adapters.

func NewCommands(a *Adapters, as *AppServices) *Commands {
	return &Commands{
		Base:      commands.NewBaseCommands(a.UserBases, a.Sectors, a.Content, as.Provisioner, a.TxMgr),
		Army:      commands.NewArmyCommands(a.UserBases, a.ArmyPrototypes, a.Users, a.OutboxEvents, a.Scheduler, a.TxMgr, as.Access),
		Building:  commands.NewBuildingCommands(a.UserBases, a.BuildPrototypes, a.Users, a.OutboxEvents, a.Scheduler, a.TxMgr, as.Access),
		Tech:      commands.NewTechCommands(a.UserBases, a.TechPrototypes, a.Users, a.OutboxEvents, a.Scheduler, a.TxMgr, as.Access),
		Storage:   commands.NewStorageCommands(a.UserBases, a.Users, a.StoragePrototypes, a.ArmyPrototypes, a.ResourceLocations, a.DangerousLocations, a.ScanReports, a.OutboxEvents, a.Scheduler, a.TxMgr, as.Access),
		Operation: commands.NewOperationCommands(a.UserBases, a.Users, a.Sectors, a.MilitaryOps, a.ResourceLocations, a.DangerousLocations, a.ScanReports, as.Provisioner, a.Scheduler, a.OutboxEvents, a.TxMgr, as.Access),
		Scanner:   commands.NewIntelligenceScannerCommands(a.UserBases, a.Sectors, a.ResourceLocations, a.DangerousLocations, a.ScanReports, as.Provisioner, a.Scheduler, a.OutboxEvents, a.TxMgr),
		Radar:     commands.NewIntelligenceRadarCommands(a.UserBases, a.MilitaryOps, a.Activities, a.Scheduler, a.OutboxEvents, a.TxMgr),
		User:      commands.NewUserCommands(a.Users, a.Hasher, a.Tokens, a.OutboxEvents, a.TxMgr),
		Activity:  commands.NewActivityCommands(a.Activities, a.MilitaryOps, a.Sectors, a.UserBases, a.ScanReports),
		World:     commands.NewWorldGenerationCommands(a.UserBases, a.Sectors, a.ResourceLocations, a.DangerousLocations, a.Content, as.Provisioner, a.Scheduler, a.TxMgr),
	}
}
