package commands

import (
	"context"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
)

type IntelligenceRadarCommands struct {
	BaseRepo        ports.UserBaseRepository
	OpRepo          ports.MilitaryOperationRepository
	RadarThreatRepo ports.RadarThreatRepository
	intelService    *domain.IntelligenceService
	Scheduler       ports.Scheduler
	Outbox          ports.OutboxEventRepository
	TxMgr           ports.TransactionManager
}

func NewIntelligenceRadarCommands(baseRepo ports.UserBaseRepository, opRepo ports.MilitaryOperationRepository, radarThreatRepo ports.RadarThreatRepository, scheduler ports.Scheduler, outbox ports.OutboxEventRepository, txMgr ports.TransactionManager) *IntelligenceRadarCommands {
	return &IntelligenceRadarCommands{
		BaseRepo:        baseRepo,
		OpRepo:          opRepo,
		RadarThreatRepo: radarThreatRepo,
		intelService:    domain.NewIntelligenceService(),
		Scheduler:       scheduler,
		Outbox:          outbox,
		TxMgr:           txMgr,
	}
}

func (c *IntelligenceRadarCommands) HandleMilitaryOperationStartedEvent(ctx context.Context, event domain.MilitaryOperationStartedEvent) error {
	op, err := c.OpRepo.FindByID(ctx, event.OperationID)
	if err != nil {
		return err
	}

	targetBase, err := c.BaseRepo.FindByCoordinates(ctx, op.TargetCoordinates.X, op.TargetCoordinates.Y)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return nil // Target is not a user base
		}
		return err
	}

	for _, b := range targetBase.BuildingsPresent {
		if b.Prototype.IntelligenceData != nil && b.Prototype.IntelligenceData.Subtype == domain.IntelligenceSubtypeRadar {
			radius := b.Prototype.IntelligenceData.ScanRange
			if radius <= 0 {
				continue
			}

			detectAt, err := op.TimeBeforeEntersCircle(targetBase.Coordinates, radius)
			if err != nil {
				continue
			}

			if err := c.Scheduler.Schedule(ctx, ports.IntelligenceRadarJob{
				BaseID:      targetBase.ID,
				OperationID: op.ID,
			}, detectAt); err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *IntelligenceRadarCommands) HandleIntelligenceRadarJob(ctx context.Context, job ports.IntelligenceRadarJob) error {
	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		// 1. Idempotency: skip if radar threat already exists for this base and op
		exists, err := c.RadarThreatRepo.Tx(tx).RadarThreatExists(ctx, job.BaseID, job.OperationID)
		if err != nil || exists {
			return err
		}

		op, err := c.OpRepo.Tx(tx).FindByID(ctx, job.OperationID)
		if err != nil {
			return err
		}
		// Only detect if it's still inbound
		if op.Phase != domain.OperationPhaseOutbound {
			return nil
		}

		// 2. Double check if radar still exists (though job should have been filtered at scheduling)
		base, err := c.BaseRepo.Tx(tx).FindByID(ctx, job.BaseID)
		if err != nil {
			return err
		}

		if !c.intelService.ResolveRadarDetection(base, op) {
			return nil
		}

		// 3. Create Radar Threat
		threat := domain.NewRadarThreat(op, job.BaseID)
		if err := c.RadarThreatRepo.Tx(tx).Create(ctx, threat); err != nil {
			return err
		}

		// 4. Save events to outbox
		return c.Outbox.Tx(tx).Save(ctx, threat.PullEvents())
	})
}
