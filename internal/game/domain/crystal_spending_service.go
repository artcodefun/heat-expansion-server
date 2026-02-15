package domain

import (
	"math"

	"github.com/google/uuid"
)

// CrystalSpendingService is a pure domain service for crystal-based speedup operations.

type CrystalSpendingService struct{}

func NewCrystalSpendingService() *CrystalSpendingService {
	return &CrystalSpendingService{}
}

// SpeedUpBuildingProduction deducts crystals from user and speeds up building production in the base aggregate.
func (s *CrystalSpendingService) SpeedUpBuildingProduction(user *User, base *UserBaseModel, buildingItemID uuid.UUID) error {
	idx := -1
	for i, item := range base.BuildingsInProduction {
		if item.BaseOwnedItem.ID == buildingItemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return NewError("error.domain.building.in_production_not_found", H{"item_id": buildingItemID})
	}

	item := base.BuildingsInProduction[idx]

	remaining := item.CompletionDate - NowUnix()
	total := item.CompletionDate - item.StartDate
	fraction := float64(remaining) / float64(total)
	crystals := int(math.Ceil(float64(item.CrystalsSkipPrice) * fraction))
	if crystals < 1 {
		crystals = 1 // Minimum price
	}
	if err := user.SpendCrystals(crystals); err != nil {
		return err
	}
	if err := base.SpeedUpBuildingProduction(buildingItemID); err != nil {
		return err
	}

	return nil
}

// SpeedUpArmyProduction deducts crystals from user and speeds up army production in the base aggregate.
func (s *CrystalSpendingService) SpeedUpArmyProduction(user *User, base *UserBaseModel, armyItemID uuid.UUID) error {
	idx := -1
	for i, item := range base.ArmiesInProduction {
		if item.BaseOwnedItem.ID == armyItemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return NewError("error.domain.army.in_production_not_found", H{"item_id": armyItemID})
	}

	item := base.ArmiesInProduction[idx]

	remaining := item.CompletionDate - NowUnix()
	total := item.CompletionDate - item.StartDate
	fraction := float64(remaining) / float64(total)
	crystals := int(math.Ceil(float64(item.CrystalsSkipPrice) * fraction))
	if crystals < 1 {
		crystals = 1 // Minimum price
	}
	if err := user.SpendCrystals(crystals); err != nil {
		return err
	}
	if err := base.SpeedUpArmyProduction(armyItemID); err != nil {
		return err
	}

	return nil
}

// SpeedUpTechResearch deducts crystals from user and speeds up tech research in the base aggregate.
func (s *CrystalSpendingService) SpeedUpTechResearch(user *User, base *UserBaseModel, techItemID uuid.UUID) error {
	idx := -1
	for i, item := range base.TechnologiesInProgress {
		if item.BaseOwnedItem.ID == techItemID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return NewError("error.domain.tech.in_progress_not_found", H{"item_id": techItemID})
	}

	item := base.TechnologiesInProgress[idx]

	remaining := item.CompletionDate - NowUnix()
	total := item.CompletionDate - item.StartDate
	fraction := float64(remaining) / float64(total)
	crystals := int(math.Ceil(float64(item.CrystalsSkipPrice) * fraction))
	if crystals < 1 {
		crystals = 1 // Minimum price
	}
	if err := user.SpendCrystals(crystals); err != nil {
		return err
	}
	if err := base.SpeedUpTechResearch(techItemID); err != nil {
		return err
	}

	return nil
}

// SpeedUpOperation deducts crystals from user and speeds up an in-flight military operation.
// It supports both outbound and return legs by fast-forwarding the operation to arrival.
func (s *CrystalSpendingService) SpeedUpOperation(user *User, op *MilitaryOperation) error {
	// Only operations currently traveling can be sped up.
	if op.Phase != OperationPhaseOutbound && op.Phase != OperationPhaseReturning {
		return NewError("error.domain.operation.not_in_travel_phase", nil)
	}

	now := NowUnix()
	var departAt, arriveAt int64
	switch op.Phase {
	case OperationPhaseOutbound:
		departAt = op.OutboundDepartAt
		arriveAt = op.OutboundArriveAt
	case OperationPhaseReturning:
		departAt = op.ReturnDepartAt
		arriveAt = op.ReturnArriveAt
	}

	// Nothing meaningful to speed up if timing is invalid or already at/after arrival.
	if arriveAt <= departAt || now >= arriveAt {
		return NewError("error.domain.operation.no_travel_remaining", nil)
	}

	total := arriveAt - departAt
	remaining := arriveAt - now
	if total <= 0 || remaining <= 0 {
		return NewError("error.domain.operation.no_travel_remaining", nil)
	}

	fraction := float64(remaining) / float64(total)
	crystals := int(math.Ceil(float64(op.CrystalsSkipPrice) * fraction))
	if crystals < 1 {
		crystals = 1 // Minimum price
	}
	if err := user.SpendCrystals(crystals); err != nil {
		return err
	}

	// Fast-forward the leg and trigger normal arrival logic.
	switch op.Phase {
	case OperationPhaseOutbound:
		// Set arrival to now and invoke arrival handler.
		op.OutboundArriveAt = now
		op.OnArrive()
	case OperationPhaseReturning:
		op.ReturnArriveAt = now
		op.OnReturnArrive()
	}

	return nil
}
