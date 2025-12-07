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
	ListNewTechItems(ctx QueryContext, baseID int) ([]*readmodels.TechItemNew, error)
	ListInResearchTechItems(ctx QueryContext, baseID int) ([]*readmodels.TechItemInProgress, error)
	ListDoneTechItems(ctx QueryContext, baseID int) ([]*readmodels.TechItemDone, error)
}

// StorageQueries: storage & buffs.
type StorageQueries interface {
	ListPresentStorageItems(ctx QueryContext, baseID int) ([]*readmodels.StorageItemPresent, error)
	// Buffs may be represented via storage prototypes with BuffData activated; adjust as needed.
}

// SectorQueries: sector intelligence & scans.
type SectorQueries interface {
	GetSector(ctx QueryContext, x, y int) (*readmodels.SectorModel, error)
	GetLatestScans(ctx QueryContext, baseID int) ([]*readmodels.SectorScanReport, error)
	GetScansNear(ctx QueryContext, baseID int, centerX, centerY, radius int) ([]*readmodels.SectorScanReport, error)
	// Map-related spatial summaries merged from MapQueries
	ListOccupiedCoordinates(ctx QueryContext) ([]readmodels.Vector2i, error)
	ListSectorsInRadius(ctx QueryContext, centerX, centerY, radius int) ([]*readmodels.SectorModel, error)
}

// ActivityQueries: activity feed filtering.
type ActivityQueries interface {
	ListActivities(ctx QueryContext, baseID int, limit int) ([]*readmodels.ActivityItem, error)
	ListMilitaryActivities(ctx QueryContext, baseID int, limit int) ([]*readmodels.ActivityItem, error)
}

// OperationQueries: operation state listings.
type OperationQueries interface {
	GetOperation(ctx QueryContext, operationID int) (*readmodels.MilitaryOperation, error)
	ListOperationsByBase(ctx QueryContext, baseID int) ([]*readmodels.MilitaryOperation, error)
	ListActiveOperations(ctx QueryContext, baseID int) ([]*readmodels.MilitaryOperation, error)
}
