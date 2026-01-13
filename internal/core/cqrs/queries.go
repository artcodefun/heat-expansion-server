package cqrs

import readmodels "github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"

// QueryContext carries caller identity & auth scope for authorization checks.
// Extend later with tenant, trace, feature flags, etc.
type QueryContext struct {
	UserID int
	Roles  []string
}

// UserQueries: user profile & ownership context.
type UserQueries interface {
	GetUserProfile(ctx QueryContext, userID int) (*readmodels.User, error)
}

// BaseQueries: high-level base stats.
type BaseQueries interface {
	GetBaseStats(ctx QueryContext, baseID int) (*readmodels.UserBaseStats, error)
	ListUserBases(ctx QueryContext) ([]*readmodels.UserBaseModel, error)
}

// BuildingQueries: segmented building item views matching client categories.
type BuildingQueries interface {
	ListNewBuildItems(ctx QueryContext, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemNew, error)
	ListPendingBuildItems(ctx QueryContext, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPending, error)
	ListInProductionBuildItems(ctx QueryContext, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemInProduction, error)
	ListPresentBuildItems(ctx QueryContext, baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPresent, error)
}

// ArmyQueries: segmented army item views.
type ArmyQueries interface {
	ListNewArmyItems(ctx QueryContext, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemNew, error)
	ListPendingArmyItems(ctx QueryContext, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPending, error)
	ListInProductionArmyItems(ctx QueryContext, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemInProduction, error)
	ListPresentArmyItems(ctx QueryContext, baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPresent, error)
}

// TechQueries: technology research lifecycle.
type TechQueries interface {
	ListNewTechItems(ctx QueryContext, baseID int, category readmodels.TechCategory) ([]*readmodels.TechItemNew, error)
	ListInResearchTechItems(ctx QueryContext, baseID int, category readmodels.TechCategory) ([]*readmodels.TechItemInProgress, error)
	ListDoneTechItems(ctx QueryContext, baseID int, category readmodels.TechCategory) ([]*readmodels.TechItemDone, error)
}

// StorageQueries: storage & buffs.
type StorageQueries interface {
	ListPresentStorageItems(ctx QueryContext, baseID int, category readmodels.StorageCategory) ([]*readmodels.StorageItemPresent, error)
	// Buffs may be represented via storage prototypes with BuffData activated; adjust as needed.
}

// SectorQueries: sector scan reports only.
type SectorQueries interface {
	GetScansNear(ctx QueryContext, baseID int, centerX, centerY, radius int) ([]*readmodels.SectorScanReport, error)
	GetScanReportByID(ctx QueryContext, baseID, id int) (*readmodels.SectorScanReport, error)
	GetLatestScanBefore(ctx QueryContext, baseID, x, y int, before int64) (*readmodels.SectorScanReport, error)
}

// ActivityQueries: activity feed filtering.
type ActivityQueries interface {
	ListOffenseActivities(ctx QueryContext, baseID int, subtype readmodels.OffenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error)
	ListDefenseActivities(ctx QueryContext, baseID int, subtype readmodels.DefenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error)
	ListScanActivities(ctx QueryContext, baseID int, limit int) ([]*readmodels.ActivityItem, error)
	ListRadarActivities(ctx QueryContext, baseID int, limit int) ([]*readmodels.ActivityItem, error)
	ListTradeActivities(ctx QueryContext, baseID int, limit int) ([]*readmodels.ActivityItem, error)
}

// OperationQueries: operation state listings.
type OperationQueries interface {
	GetOperation(ctx QueryContext, operationID int) (*readmodels.MilitaryOperation, error)
	ListOperationsByBase(ctx QueryContext, baseID int) ([]*readmodels.MilitaryOperation, error)
	ListActiveOperations(ctx QueryContext, baseID int) ([]*readmodels.MilitaryOperation, error)
}

// RadarQueries: incoming threats tracking.
type RadarQueries interface {
	ListIncomingThreats(ctx QueryContext, baseID int) ([]*readmodels.RadarThreat, error)
}
