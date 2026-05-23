package cqrs

import (
	"context"

	readmodels "github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

// UserQueries: user profile & ownership context.
type UserQueries interface {
	GetUserProfile(ctx context.Context, actor Actor, userID uuid.UUID) (*readmodels.User, error)
}

type BlackMarketQueries interface {
	ListResourceRates(ctx context.Context, actor Actor, baseID int) ([]*readmodels.BlackMarketResourceRate, error)
	ListActiveOffers(ctx context.Context, actor Actor, baseID int, kind *domain.BlackMarketOfferKind, limited *bool) ([]*readmodels.BlackMarketOffer, error)
}

// BaseQueries: high-level base stats.
type BaseQueries interface {
	GetBaseStats(ctx context.Context, actor Actor, baseID int) (*readmodels.UserBaseStats, error)
	ListUserBases(ctx context.Context, actor Actor) ([]*readmodels.UserBaseModel, error)
}

// BuildingQueries: segmented building item views matching client categories.
type BuildingQueries interface {
	ListNewBuildItems(ctx context.Context, actor Actor, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemNew, error)
	ListPendingBuildItems(ctx context.Context, actor Actor, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPending, error)
	ListInProductionBuildItems(ctx context.Context, actor Actor, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemInProduction, error)
	ListPresentBuildItems(ctx context.Context, actor Actor, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPresent, error)
}

// ArmyQueries: segmented army item views.
type ArmyQueries interface {
	ListNewArmyItems(ctx context.Context, actor Actor, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemNew, error)
	ListPendingArmyItems(ctx context.Context, actor Actor, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPending, error)
	ListInProductionArmyItems(ctx context.Context, actor Actor, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemInProduction, error)
	ListPresentArmyItems(ctx context.Context, actor Actor, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPresent, error)
}

// TechQueries: technology research lifecycle.
type TechQueries interface {
	ListNewTechItems(ctx context.Context, actor Actor, baseID int, category readmodels.TechCategory) ([]*readmodels.TechItemNew, error)
	ListInResearchTechItems(ctx context.Context, actor Actor, baseID int, category readmodels.TechCategory) ([]*readmodels.TechItemInProgress, error)
	ListDoneTechItems(ctx context.Context, actor Actor, baseID int, category readmodels.TechCategory) ([]*readmodels.TechItemDone, error)
}

// StorageQueries: storage & buffs.
type StorageQueries interface {
	ListPresentStorageItems(ctx context.Context, actor Actor, baseID int, category readmodels.StorageCategory) ([]*readmodels.StorageItemPresent, error)
	// Buffs may be represented via storage prototypes with BuffData activated; adjust as needed.
}

// TradeQueries: ally/base-owner inventory reads for trade preparation.
type TradeQueries interface {
	GetTradeInfo(ctx context.Context, actor Actor, targetX, targetY int) (*readmodels.TradeInfo, error)
	GetTradeOperation(ctx context.Context, actor Actor, baseID int, operationID int) (*readmodels.TradeOperation, error)
	ListActiveTradeOperations(ctx context.Context, actor Actor, baseID int) ([]*readmodels.TradeOperation, error)
}

// SectorQueries: sector scan reports only.
type SectorQueries interface {
	GetScansNear(ctx context.Context, actor Actor, baseID int, centerX, centerY, radius int) ([]*readmodels.SectorScanReport, error)
	GetScanReportByID(ctx context.Context, actor Actor, baseID, id int) (*readmodels.SectorScanReport, error)
	GetLatestScanBefore(ctx context.Context, actor Actor, baseID, x, y int, before int64) (*readmodels.SectorScanReport, error)
}

// ActivityQueries: activity feed filtering.
type ActivityQueries interface {
	ListOffenseActivities(ctx context.Context, actor Actor, baseID int, subtype readmodels.OffenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error)
	ListDefenseActivities(ctx context.Context, actor Actor, baseID int, subtype readmodels.DefenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error)
	ListScanActivities(ctx context.Context, actor Actor, baseID int, subtype readmodels.ScanActivitySubtype, limit int) ([]*readmodels.ActivityItem, error)
	ListRadarActivities(ctx context.Context, actor Actor, baseID int, limit int) ([]*readmodels.ActivityItem, error)
	ListTradeActivities(ctx context.Context, actor Actor, baseID int, limit int) ([]*readmodels.ActivityItem, error)
}

// OperationQueries: operation state listings.
type OperationQueries interface {
	GetOperation(ctx context.Context, actor Actor, operationID int) (*readmodels.MilitaryOperation, error)
	GetOperationByUUID(ctx context.Context, actor Actor, operationUUID uuid.UUID) (*readmodels.MilitaryOperation, error)
	ListOperationsByBase(ctx context.Context, actor Actor, baseID int) ([]*readmodels.MilitaryOperation, error)
	ListActiveOperations(ctx context.Context, actor Actor, baseID int) ([]*readmodels.MilitaryOperation, error)
}

// RadarQueries: incoming threats tracking.
type RadarQueries interface {
	ListIncomingThreats(ctx context.Context, actor Actor, baseID int) ([]*readmodels.RadarThreat, error)
}

// AlertQueries: high-priority notification feed.
type AlertQueries interface {
	ListActiveAlerts(ctx context.Context, actor Actor) ([]*readmodels.AlertItem, error)
	GetUnreadAlertsCount(ctx context.Context, actor Actor) (int, error)
}

type DiplomacyQueries interface {
	ListRelationships(ctx context.Context, actor Actor, status *readmodels.DiplomaticStatus) ([]*readmodels.DiplomaticRelationship, error)
	GetRelationship(ctx context.Context, actor Actor, otherUserID uuid.UUID) (*readmodels.DiplomaticRelationship, error)
	ListChats(ctx context.Context, actor Actor) ([]*readmodels.DiplomaticChat, error)
	GetUnreadMessagesCount(ctx context.Context, actor Actor) (int, error)
	ListChatMessages(ctx context.Context, actor Actor, otherUserID uuid.UUID) ([]*readmodels.DiplomaticMessage, error)
	ListPendingRequests(ctx context.Context, actor Actor) ([]*readmodels.DiplomaticRequest, error)
}
