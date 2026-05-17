package cqrs

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

// UserCommands encapsulates mutating user operations.
type UserCommands interface{}

// BlackMarketCommands encapsulates Black Market mutations.
type BlackMarketCommands interface {
	PurchaseResources(ctx context.Context, actor Actor, baseID int, resourceType domain.ResourceType, crystals int) error
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

// TradeCommands encapsulates trade operation mutations.
type TradeCommands interface {
	CreateTradeOperation(ctx context.Context, actor Actor, senderBaseID int, targetX, targetY int, offeredResources domain.PriceModel, offeredArmyRequests []domain.ArmyDeploymentRequest, offeredStorageItemIDs []uuid.UUID, requestedResources domain.PriceModel, requestedArmyRequests []domain.ArmyDeploymentRequest, requestedStorageItemIDs []uuid.UUID, transportRequests []domain.ArmyDeploymentRequest) (*domain.TradeOperation, error)
	AcceptTradeOperation(ctx context.Context, actor Actor, operationID int) error
	DeclineTradeOperation(ctx context.Context, actor Actor, operationID int) error
	CancelTradeOperationByInitiator(ctx context.Context, actor Actor, operationID int) error
	SpeedUpTradeOperationWithCrystals(ctx context.Context, actor Actor, operationID int) error
}

// AlertCommands encapsulates alert notifications management.
type AlertCommands interface {
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
}

type DiplomacyCommands interface {
	SendInformationalMessage(ctx context.Context, actor Actor, senderBaseID int, receiverUserID uuid.UUID, receiverBaseID *int, content domain.TranslationKey) (*uuid.UUID, error)
	SendRequest(ctx context.Context, actor Actor, senderBaseID int, receiverUserID uuid.UUID, receiverBaseID *int, kind domain.DiplomaticRequestKind) (*uuid.UUID, error)
	DeclareWar(ctx context.Context, actor Actor, senderBaseID int, receiverUserID uuid.UUID, receiverBaseID *int) (*uuid.UUID, error)
	BreakAlliance(ctx context.Context, actor Actor, senderBaseID int, receiverUserID uuid.UUID, receiverBaseID *int) (*uuid.UUID, error)
	MarkChatAsRead(ctx context.Context, actor Actor, otherUserID uuid.UUID) error
	AcceptRequest(ctx context.Context, actor Actor, senderBaseID int, requestID uuid.UUID) error
	RejectRequest(ctx context.Context, actor Actor, senderBaseID int, requestID uuid.UUID) error
}
