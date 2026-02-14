package cqrs

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

// CommandContext carries caller identity & auth scope for authorization checks on mutations.
// Mirrors QueryContext; extend with trace IDs, tenant, feature flags as needed.
type CommandContext struct {
	UserID uuid.UUID
	Roles  []string
}

// UserCommands encapsulates mutating user operations.
type UserCommands interface {
}

// BaseCommands encapsulates base creation operations.
type BaseCommands interface {
	CreateBase(ctx CommandContext, userID uuid.UUID) error
}

// BuildingCommands encapsulates building queue & production mutations.
type BuildingCommands interface {
	QueueBuilding(ctx CommandContext, baseID int, prototypeID int) error
	CancelPendingBuilding(ctx CommandContext, baseID int, itemID uuid.UUID) error
	SpeedUpProductionWithCrystals(ctx CommandContext, baseID int, buildingItemID uuid.UUID) error
	DeletePresentBuilding(ctx CommandContext, baseID int, itemID uuid.UUID) error
}

// ArmyCommands encapsulates army queue & production mutations.
type ArmyCommands interface {
	QueueArmy(ctx CommandContext, baseID int, prototypeID int, count int) error
	CancelPendingArmy(ctx CommandContext, baseID int, itemID uuid.UUID, count int) error
	SpeedUpArmyProductionWithCrystals(ctx CommandContext, baseID int, armyItemID uuid.UUID) error
	DeletePresentArmy(ctx CommandContext, baseID int, itemID uuid.UUID, count int) error
}

// TechCommands encapsulates tech research mutations.
type TechCommands interface {
	StartTechResearch(ctx CommandContext, baseID int, prototypeID int) error
	SpeedUpTechResearchWithCrystals(ctx CommandContext, baseID int, techItemID uuid.UUID) error
}

// StorageCommands encapsulates storage (buff/artifact) mutations.
type StorageCommands interface {
	DeletePresentStorageItem(ctx CommandContext, baseID int, itemID uuid.UUID) error
	ActivateBuff(ctx CommandContext, baseID int, itemID uuid.UUID) error
	StartIntelDecryption(ctx CommandContext, baseID int, itemID uuid.UUID) error
	StartDamagedItemRestoration(ctx CommandContext, baseID int, itemID uuid.UUID) error
	ActivateArtifact(ctx CommandContext, baseID int, itemID uuid.UUID) error
	DeactivateArtifact(ctx CommandContext, baseID int, itemID uuid.UUID) error
	OpenConsumableBox(ctx CommandContext, baseID int, itemID uuid.UUID) error
}

// OperationCommands encapsulates military operation life-cycle mutations.
type OperationCommands interface {
	CreateMilitaryOperation(ctx CommandContext, opType domain.MilitaryOperationType, sourceBaseID int, targetX int, targetY int, deployments []domain.ArmyDeploymentRequest) (*domain.MilitaryOperation, error)
	CancelMilitaryOperation(ctx CommandContext, operationID int) error
	SpeedUpOperationWithCrystals(ctx CommandContext, operationID int) error
}

// AlertCommands encapsulates alert notifications management.
type AlertCommands interface {
	MarkAllAsRead(baseID int, userID uuid.UUID) error
}
