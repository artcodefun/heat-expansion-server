package cqrs

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/google/uuid"
)

// CommandContext carries caller identity & auth scope for authorization checks on mutations.
// Mirrors QueryContext; extend with trace IDs, tenant, feature flags as needed.
type CommandContext struct {
	UserID int
	Roles  []string
}

// UserCommands encapsulates mutating user operations.
type UserCommands interface {
	Create(ctx CommandContext, name, email, password string) error
	Authenticate(ctx CommandContext, email, password string) (string, error)
}

// BaseCommands encapsulates base creation operations.
type BaseCommands interface {
	CreateBase(ctx CommandContext, userID int) error
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
}

// OperationCommands encapsulates military operation life-cycle mutations.
type OperationCommands interface {
	CreateMilitaryOperation(ctx CommandContext, opType domain.MilitaryOperationType, sourceBaseID int, targetX int, targetY int, deployments []domain.ArmyDeploymentRequest) (*domain.MilitaryOperation, error)
}
