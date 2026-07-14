package ports

import (
	"context"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/google/uuid"
)

// BuildingReadRepository exposes lifecycle-segmented build item projections.
type BuildingReadRepository interface {
	ListNewBuildItemsByPrototypeIDs(ctx context.Context, ids []int) ([]*readmodels.BuildItemNew, error)
	ListPendingBuildItems(ctx context.Context, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPending, error)
	ListInProductionBuildItems(ctx context.Context, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemInProduction, error)
	ListPresentBuildItems(ctx context.Context, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPresent, error)
}

// UserReadRepository provides read-only access to user profile data.
type UserReadRepository interface {
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*readmodels.User, error)
}

// BlackMarketReadRepository exposes active offer projections.
type BlackMarketReadRepository interface {
	ListActiveOffers(ctx context.Context, kind *readmodels.BlackMarketOfferKind, limited *bool) ([]*readmodels.BlackMarketOffer, error)
}

// ActivityReadRepository provides activity feed projections.
type ActivityReadRepository interface {
	ListOffenseActivities(ctx context.Context, baseID int, subtype readmodels.OffenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error)
	ListDefenseActivities(ctx context.Context, baseID int, subtype readmodels.DefenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error)
	ListScanActivities(ctx context.Context, baseID int, subtype readmodels.ScanActivitySubtype, limit int) ([]*readmodels.ActivityItem, error)
	ListRadarActivities(ctx context.Context, baseID int, limit int) ([]*readmodels.ActivityItem, error)
	ListTradeActivities(ctx context.Context, baseID int, limit int) ([]*readmodels.ActivityItem, error)
}

// TechReadRepository exposes technology lifecycle projections.
type TechReadRepository interface {
	ListNewTechItemsByPrototypeIDs(ctx context.Context, baseID int, ids []int) ([]*readmodels.TechItemNew, error)
	ListInResearchTechItems(ctx context.Context, baseID int, category readmodels.TechCategory) ([]*readmodels.TechItemInProgress, error)
	ListDoneTechItems(ctx context.Context, baseID int, category readmodels.TechCategory) ([]*readmodels.TechItemDone, error)
}

// OperationReadRepository exposes military operation projections.
type OperationReadRepository interface {
	GetOperation(ctx context.Context, operationID int) (*readmodels.MilitaryOperation, error)
	GetOperationByUUID(ctx context.Context, operationUUID uuid.UUID) (*readmodels.MilitaryOperation, error)
	ListOperationsByBase(ctx context.Context, baseID int) ([]*readmodels.MilitaryOperation, error)
	ListActiveOperations(ctx context.Context, baseID int) ([]*readmodels.MilitaryOperation, error)
}

// TradeOperationReadRepository exposes trade operation projections.
type TradeOperationReadRepository interface {
	GetTradeOperation(ctx context.Context, operationID int) (*readmodels.TradeOperation, error)
	ListActiveTradeOperations(ctx context.Context, baseID int) ([]*readmodels.TradeOperation, error)
}

// SectorReadRepository provides sector scan report projections.
type SectorReadRepository interface {
	GetScansNear(ctx context.Context, baseID int, centerX, centerY, radius int) ([]*readmodels.SectorScanReport, error)
	GetScanReportByID(ctx context.Context, baseID, id int) (*readmodels.SectorScanReport, error)
	GetLatestScanBefore(ctx context.Context, baseID, x, y int, before int64) (*readmodels.SectorScanReport, error)
}

// AlertReadRepository provides high-priority notification projections.
type AlertReadRepository interface {
	ListActiveAlerts(ctx context.Context, userID uuid.UUID) ([]*readmodels.AlertItem, error)
	GetUnreadAlertsCount(ctx context.Context, userID uuid.UUID) (int, error)
}

type DiplomacyReadRepository interface {
	ListRelationships(ctx context.Context, userID uuid.UUID, status *readmodels.DiplomaticStatus) ([]*readmodels.DiplomaticRelationship, error)
	GetRelationship(ctx context.Context, userID, otherUserID uuid.UUID) (*readmodels.DiplomaticRelationship, error)
	ListChats(ctx context.Context, userID uuid.UUID) ([]*readmodels.DiplomaticChat, error)
	GetUnreadMessagesCount(ctx context.Context, userID uuid.UUID) (int, error)
	ListChatMessages(ctx context.Context, userID, otherUserID uuid.UUID) ([]*readmodels.DiplomaticMessage, error)
	ListPendingRequests(ctx context.Context, userID uuid.UUID) ([]*readmodels.DiplomaticRequest, error)
}

// StorageReadRepository exposes storage item / buff projections.
type StorageReadRepository interface {
	ListPresentStorageItems(ctx context.Context, baseID int, category readmodels.StorageCategory) ([]*readmodels.StorageItemPresent, error)
	ListTradeableStorageItems(ctx context.Context, baseID int) ([]*readmodels.StorageItemPresent, error)
}

// BaseReadRepository provides read-only access to base state.
type BaseReadRepository interface {
	GetBase(ctx context.Context, baseID int) (*readmodels.UserBaseModel, error)
	GetBaseOwnerByCoordinates(ctx context.Context, x, y int) (*readmodels.SectorOwner, error)
	GetBaseStats(ctx context.Context, baseID int) (*readmodels.UserBaseStats, error)
	ListUserBases(ctx context.Context, userID uuid.UUID) ([]*readmodels.UserBaseModel, error)
}

// ArmyReadRepository exposes lifecycle-segmented army item projections.
type ArmyReadRepository interface {
	ListNewArmyItemsByPrototypeIDs(ctx context.Context, ids []int) ([]*readmodels.ArmyItemNew, error)
	ListPendingArmyItems(ctx context.Context, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPending, error)
	ListInProductionArmyItems(ctx context.Context, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemInProduction, error)
	ListPresentArmyItems(ctx context.Context, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPresent, error)
	ListPresentArmyItemsAll(ctx context.Context, baseID int) ([]*readmodels.ArmyItemPresent, error)
}

// ArmyPrototypeReadRepository provides read-only access to the army prototype catalog.
type ArmyPrototypeReadRepository interface {
	ListArmyPrototypes(ctx context.Context) ([]*readmodels.ArmyItemPrototype, error)
	GetArmyPrototype(ctx context.Context, id int) (*readmodels.ArmyItemPrototype, error)
}

// BuildPrototypeReadRepository provides read-only access to the build prototype catalog.
type BuildPrototypeReadRepository interface {
	ListBuildPrototypes(ctx context.Context) ([]*readmodels.BuildItemPrototype, error)
	GetBuildPrototype(ctx context.Context, id int) (*readmodels.BuildItemPrototype, error)
}

// StoragePrototypeReadRepository provides read-only access to the storage prototype catalog.
type StoragePrototypeReadRepository interface {
	ListStoragePrototypes(ctx context.Context) ([]*readmodels.StorageItemPrototype, error)
	GetStoragePrototype(ctx context.Context, id int) (*readmodels.StorageItemPrototype, error)
}

// TechPrototypeReadRepository provides read-only access to the tech prototype catalog.
type TechPrototypeReadRepository interface {
	ListTechPrototypes(ctx context.Context) ([]*readmodels.TechItemPrototype, error)
	GetTechPrototype(ctx context.Context, id int) (*readmodels.TechItemPrototype, error)
}

// RadarReadRepository provides read-only access to radar threats.
type RadarReadRepository interface {
	GetRadarThreat(ctx context.Context, id uuid.UUID) (*readmodels.RadarThreat, error)
	ListIncomingThreats(ctx context.Context, baseID int) ([]*readmodels.RadarThreat, error)
}
