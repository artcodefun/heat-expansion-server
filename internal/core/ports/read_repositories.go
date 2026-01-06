package ports

import "github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"

// BuildingReadRepository exposes lifecycle-segmented build item projections.
type BuildingReadRepository interface {
	ListNewBuildItemsByPrototypeIDs(ids []int) ([]*readmodels.BuildItemNew, error)
	ListPendingBuildItems(baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPending, error)
	ListInProductionBuildItems(baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemInProduction, error)
	ListPresentBuildItems(baseID int, category readmodels.BuildCategory) ([]*readmodels.BuildItemPresent, error)
}

// UserReadRepository provides read-only access to user profile data.
type UserReadRepository interface {
	GetUserProfile(userID int) (*readmodels.User, error)
}

// ActivityReadRepository provides activity feed projections.
type ActivityReadRepository interface {
	ListOffenseActivities(baseID int, subtype readmodels.OffenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error)
	ListDefenseActivities(baseID int, subtype readmodels.DefenseActivitySubtype, limit int) ([]*readmodels.ActivityItem, error)
	ListScanActivities(baseID int, limit int) ([]*readmodels.ActivityItem, error)
	ListRadarActivities(baseID int, limit int) ([]*readmodels.ActivityItem, error)
	ListTradeActivities(baseID int, limit int) ([]*readmodels.ActivityItem, error)
}

// TechReadRepository exposes technology lifecycle projections.
type TechReadRepository interface {
	ListNewTechItemsByPrototypeIDs(ids []int) ([]*readmodels.TechItemNew, error)
	ListInResearchTechItems(baseID int) ([]*readmodels.TechItemInProgress, error)
	ListDoneTechItems(baseID int) ([]*readmodels.TechItemDone, error)
}

// OperationReadRepository exposes military operation projections.
type OperationReadRepository interface {
	GetOperation(operationID int) (*readmodels.MilitaryOperation, error)
	ListOperationsByBase(baseID int) ([]*readmodels.MilitaryOperation, error)
	ListActiveOperations(baseID int) ([]*readmodels.MilitaryOperation, error)
}

// SectorReadRepository provides sector scan report projections.
type SectorReadRepository interface {
	GetScansNear(baseID int, centerX, centerY, radius int) ([]*readmodels.SectorScanReport, error)
	GetScanReportByID(baseID, id int) (*readmodels.SectorScanReport, error)
	GetLatestScanBefore(baseID, x, y int, before int64) (*readmodels.SectorScanReport, error)
}

// StorageReadRepository exposes storage item / buff projections.
type StorageReadRepository interface {
	ListPresentStorageItems(baseID int) ([]*readmodels.StorageItemPresent, error)
}

// BaseReadRepository provides read-only access to base state.
type BaseReadRepository interface {
	GetBaseStats(baseID int) (*readmodels.UserBaseStats, error)
	ListUserBases(userID int) ([]*readmodels.UserBaseModel, error)
}

// ArmyReadRepository exposes lifecycle-segmented army item projections.
type ArmyReadRepository interface {
	ListNewArmyItemsByPrototypeIDs(ids []int) ([]*readmodels.ArmyItemNew, error)
	ListPendingArmyItems(baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPending, error)
	ListInProductionArmyItems(baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemInProduction, error)
	ListPresentArmyItems(baseID int, category readmodels.ArmyCategory) ([]*readmodels.ArmyItemPresent, error)
}
