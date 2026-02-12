package bootstrap

import "github.com/artcodefun/heat-expansion-server/internal/game/application/commands"

// Commands aggregates all command handlers.
type Commands struct {
	Base        *commands.BaseCommands
	Army        *commands.ArmyCommands
	Building    *commands.BuildingCommands
	Tech        *commands.TechCommands
	Storage     *commands.StorageCommands
	Operation   *commands.OperationCommands
	Scanner     *commands.IntelligenceScannerCommands
	Radar       *commands.IntelligenceRadarCommands
	RadarThreat *commands.RadarThreatCommands
	User        *commands.UserCommands
	Activity    *commands.ActivityCommands
	World       *commands.WorldGenerationCommands
	Alert       *commands.AlertCommands
}

// NewCommands constructs all command handlers using provided secondary adapters.

func NewCommands(a *Adapters, as *AppServices) *Commands {
	return &Commands{
		Base:        commands.NewBaseCommands(a.UserBases, a.Sectors, a.BuildPrototypes, a.ArmyPrototypes, a.Content, as.Provisioner, a.OutboxEvents, a.TxMgr),
		Army:        commands.NewArmyCommands(a.UserBases, a.ArmyPrototypes, a.Users, a.OutboxEvents, a.Scheduler, a.TxMgr, as.Access),
		Building:    commands.NewBuildingCommands(a.UserBases, a.BuildPrototypes, a.Users, a.OutboxEvents, a.Scheduler, a.TxMgr, as.Access),
		Tech:        commands.NewTechCommands(a.UserBases, a.TechPrototypes, a.Users, a.OutboxEvents, a.Scheduler, a.TxMgr, as.Access),
		Storage:     commands.NewStorageCommands(a.UserBases, a.Users, a.Sectors, a.StoragePrototypes, a.ArmyPrototypes, a.ResourceLocations, a.DangerousLocations, a.ScanReports, a.OutboxEvents, a.Scheduler, a.TxMgr, as.Access),
		Operation:   commands.NewOperationCommands(a.UserBases, a.Users, a.Sectors, a.MilitaryOps, a.ResourceLocations, a.DangerousLocations, a.ScanReports, a.StoragePrototypes, as.Provisioner, a.Scheduler, a.OutboxEvents, a.TxMgr, as.Access),
		Scanner:     commands.NewIntelligenceScannerCommands(a.UserBases, a.Sectors, a.ResourceLocations, a.DangerousLocations, a.ScanReports, as.Provisioner, a.Scheduler, a.OutboxEvents, a.TxMgr),
		Radar:       commands.NewIntelligenceRadarCommands(a.UserBases, a.MilitaryOps, a.RadarThreats, a.Scheduler, a.OutboxEvents, a.TxMgr),
		RadarThreat: commands.NewRadarThreatCommands(a.RadarThreats, a.OutboxEvents, a.TxMgr),
		User:        commands.NewUserCommands(a.Users, a.Hasher, a.Tokens, a.OutboxEvents, a.TxMgr),
		Activity:    commands.NewActivityCommands(a.Activities, a.MilitaryOps, a.RadarThreats, a.Sectors, a.UserBases, a.ScanReports, a.OutboxEvents, a.TxMgr),
		World:       commands.NewWorldGenerationCommands(a.UserBases, a.Sectors, a.ResourceLocations, a.DangerousLocations, a.StoragePrototypes, a.ArmyPrototypes, a.BuildPrototypes, a.Content, as.Provisioner, a.Scheduler, a.TxMgr),
		Alert:       commands.NewAlertCommands(a.Alerts, as.Access, a.TxMgr),
	}
}
