package commands

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
)

type ActivityCommands struct {
	ActivityRepo  ports.ActivityRepository
	OperationRepo ports.MilitaryOperationRepository
	SectorRepo    ports.SectorRepository
	UserBaseRepo  ports.UserBaseRepository
	ScanRepo      ports.ScanReportRepository
}

func NewActivityCommands(activityRepo ports.ActivityRepository, opRepo ports.MilitaryOperationRepository, sectorRepo ports.SectorRepository, baseRepo ports.UserBaseRepository, scanRepo ports.ScanReportRepository) *ActivityCommands {
	return &ActivityCommands{ActivityRepo: activityRepo, OperationRepo: opRepo, SectorRepo: sectorRepo, UserBaseRepo: baseRepo, ScanRepo: scanRepo}
}

func (c *ActivityCommands) HandleMilitaryOperationStartedEvent(event domain.MilitaryOperationStartedEvent) error {
	op, err := c.OperationRepo.FindByID(event.OperationID)
	if err != nil {
		return err
	}
	item := domain.NewActivityFromOperation(op.SourceBaseID, op, domain.OperationRoleAttacker)
	return c.ActivityRepo.Create(&item)
}

func (c *ActivityCommands) HandleMilitaryOperationResolvedEvent(event domain.MilitaryOperationResolvedEvent) error {
	op, err := c.OperationRepo.FindByID(event.OperationID)
	if err != nil {
		return err
	}
	occType, _ := c.SectorRepo.GetLocationTypeByCoordinates(op.TargetCoordinates.X, op.TargetCoordinates.Y)
	if occType != domain.LocationTypeUserBase {
		return nil
	}
	base, err := c.UserBaseRepo.FindByCoordinates(op.TargetCoordinates.X, op.TargetCoordinates.Y)
	if err != nil {
		return err
	}
	item := domain.NewActivityFromDefense(base.ID, op)
	if ts := event.OccurredAt(); ts != 0 {
		item.CreatedAt = ts
	}
	return c.ActivityRepo.Create(&item)
}

func (c *ActivityCommands) HandleScanReportCreatedEvent(event domain.ScanReportCreatedEvent) error {
	report, err := c.ScanRepo.FindByID(event.ReportID)
	if err != nil {
		return err
	}
	item := domain.NewActivityFromScan(event.BaseID, report)
	return c.ActivityRepo.Create(&item)
}
