package cqrs

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

// UserCommands encapsulates mutating user operations.
type UserCommands interface {
}

// BaseCommands encapsulates base creation operations.
type BaseCommands interface {
}

// BuildingCommands encapsulates building queue & production mutations.
type BuildingCommands interface {
	QueueBuilding(ctx context.Context, actor Actor, baseID int, prototypeID int) error
	CancelPendingBuilding(ctx context.Context, actor Actor, baseID int, itemID uuid.UUID) error
	SpeedUpProductionWithCrystals(ctx context.Context, actor Actor, baseID int, buildingItemID uuid.UUID) error
	DeletePresentBuilding(ctx context.Context, actor Actor, baseID int, itemID uuid.UUID) error
}

// ArmyCommands encapsulates army queue & production mutations.
type ArmyCommands interface {
	QueueArmy(ctx context.Context, actor Actor, baseID int, prototypeID int, count int) error
	CancelPendingArmy(ctx context.Context, actor Actor, baseID int, itemID uuid.UUID, count int) error
	SpeedUpArmyProductionWithCrystals(ctx context.Context, actor Actor, baseID int, armyItemID uuid.UUID) error
	DeletePresentArmy(ctx context.Context, actor Actor, baseID int, itemID uuid.UUID, count int) error
}

// TechCommands encapsulates tech research mutations.
type TechCommands interface {
	StartTechResearch(ctx context.Context, actor Actor, baseID int, prototypeID int) error
	SpeedUpTechResearchWithCrystals(ctx context.Context, actor Actor, baseID int, techItemID uuid.UUID) error
}

// StorageCommands encapsulates storage (buff/artifact) mutations.
type StorageCommands interface {
	DeletePresentStorageItem(ctx context.Context, actor Actor, baseID int, itemID uuid.UUID) error
	ActivateBuff(ctx context.Context, actor Actor, baseID int, itemID uuid.UUID) error
	StartIntelDecryption(ctx context.Context, actor Actor, baseID int, itemID uuid.UUID) error
	StartDamagedItemRestoration(ctx context.Context, actor Actor, baseID int, itemID uuid.UUID) error
	ActivateArtifact(ctx context.Context, actor Actor, baseID int, itemID uuid.UUID) error
	DeactivateArtifact(ctx context.Context, actor Actor, baseID int, itemID uuid.UUID) error
	OpenConsumableBox(ctx context.Context, actor Actor, baseID int, itemID uuid.UUID) error
}

// OperationCommands encapsulates military operation life-cycle mutations.
type OperationCommands interface {
	CreateMilitaryOperation(ctx context.Context, actor Actor, opType domain.MilitaryOperationType, sourceBaseID int, targetX int, targetY int, deployments []domain.ArmyDeploymentRequest) (*domain.MilitaryOperation, error)
	CancelMilitaryOperation(ctx context.Context, actor Actor, operationID int) error
	SpeedUpOperationWithCrystals(ctx context.Context, actor Actor, operationID int) error
}

// AlertCommands encapsulates alert notifications management.
type AlertCommands interface {
	MarkAllAsRead(ctx context.Context, baseID int, userID uuid.UUID) error
}
