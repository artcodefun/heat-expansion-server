package ports

import "github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"

// BuildingReadRepository exposes lifecycle-segmented build item projections.
type BuildingReadRepository interface {
	ListNewBuildItems(baseID int, category string) ([]*readmodels.BuildItemNew, error)
	ListPendingBuildItems(baseID int, category string) ([]*readmodels.BuildItemPending, error)
	ListInProductionBuildItems(baseID int, category string) ([]*readmodels.BuildItemInProduction, error)
	ListPresentBuildItems(baseID int, category string) ([]*readmodels.BuildItemPresent, error)
}

// UserReadRepository provides read-only access to user profile data.
type UserReadRepository interface {
	GetUserProfile(userID int) (*readmodels.User, error)
}

// ActivityReadRepository provides activity feed projections.
type ActivityReadRepository interface {
	ListActivities(baseID int, limit int) ([]*readmodels.ActivityItem, error)
	ListMilitaryActivities(baseID int, limit int) ([]*readmodels.ActivityItem, error)
}

// TechReadRepository exposes technology lifecycle projections.
type TechReadRepository interface {
	ListNewTechItems(baseID int) ([]*readmodels.TechItemNew, error)
	ListInResearchTechItems(baseID int) ([]*readmodels.TechItemInProgress, error)
	ListDoneTechItems(baseID int) ([]*readmodels.TechItemDone, error)
}

// OperationReadRepository exposes military operation projections.
type OperationReadRepository interface {
	GetOperation(operationID int) (*readmodels.MilitaryOperation, error)
	ListOperationsByBase(baseID int) ([]*readmodels.MilitaryOperation, error)
	ListActiveOperations(baseID int) ([]*readmodels.MilitaryOperation, error)
}

// SectorReadRepository provides sector and scan report projections.
type SectorReadRepository interface {
	GetSector(x, y int) (*readmodels.SectorModel, error)
	GetLatestScans(baseID int) ([]*readmodels.SectorScanReport, error)
	GetScansNear(baseID int, centerX, centerY, radius int) ([]*readmodels.SectorScanReport, error)
	ListOccupiedCoordinates() ([]readmodels.Vector2i, error)
	ListSectorsInRadius(centerX, centerY, radius int) ([]*readmodels.SectorModel, error)
}

// StorageReadRepository exposes storage item / buff projections.
type StorageReadRepository interface {
	ListPresentStorageItems(baseID int) ([]*readmodels.StorageItemPresent, error)
}

// BaseReadRepository provides read-only access to base state.
type BaseReadRepository interface {
	GetBaseStats(baseID int) (*readmodels.UserBaseStats, error)
}

// ArmyReadRepository exposes lifecycle-segmented army item projections.
type ArmyReadRepository interface {
	ListNewArmyItems(baseID int, category string) ([]*readmodels.ArmyItemNew, error)
	ListPendingArmyItems(baseID int, category string) ([]*readmodels.ArmyItemPending, error)
	ListInProductionArmyItems(baseID int, category string) ([]*readmodels.ArmyItemInProduction, error)
	ListPresentArmyItems(baseID int, category string) ([]*readmodels.ArmyItemPresent, error)
}
