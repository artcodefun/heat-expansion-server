package ports

import (
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/google/uuid"
)

// ErrNotFound is a generic sentinel for missing records.
var ErrNotFound = errors.New("record not found")

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id int) (*domain.User, error)
	// FindByIDForUpdate acquires a row-level lock on the user for the duration of the transaction.
	FindByIDForUpdate(id int) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id int) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) UserRepository
}

// SectorRepository defines the interface for sector persistence.
type SectorRepository interface {
	Create(sector *domain.SectorModel) error
	Update(sector *domain.SectorModel) error
	// Finds a sector by its coordinates (x, y)
	FindByCoordinates(x int, y int) (*domain.SectorModel, error)
	// Finds and locks a sector row by coordinates within a transaction (SELECT ... FOR UPDATE)
	FindByCoordinatesForUpdate(x int, y int) (*domain.SectorModel, error)
	// Returns all sectors
	FindAll() ([]*domain.SectorModel, error)
	// Lists coordinates of all non-empty occupied locations (bases, resourceful, dangerous)
	ListOccupiedCoordinates() ([]domain.Vector2i, error)
	// Derives occupant location type for a sector by coordinates
	GetLocationTypeByCoordinates(x int, y int) (domain.LocationType, error)
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) SectorRepository
}

// ScanReportRepository defines persistence for sector scan reports.
type ScanReportRepository interface {
	Create(report *domain.SectorScanReport) error
	FindByID(id int) (*domain.SectorScanReport, error)
	FindByBaseAndCoordinates(baseID int, x int, y int) ([]*domain.SectorScanReport, error)
	// GetLatestScansByBase returns the latest scan reports for the base, ordered by CreatedAt desc.
	// Implementations may choose an internal cap; expose pagination later if needed.
	GetLatestScansByBase(baseID int) ([]*domain.SectorScanReport, error)
	// RecentReportExistsByScanner checks if a report was produced by a specific building within the last 'since' seconds.
	RecentReportExistsByScanner(scannerID uuid.UUID, since int64) (bool, error)
	// FindByBaseWithinArea returns all scan reports for a base whose sector coordinates fall within
	// the inclusive radius (Euclidean) of the provided center. This may be implemented efficiently
	// with a join against sectors; a naive implementation can load latest scans and filter in memory.
	FindByBaseWithinArea(baseID int, centerX int, centerY int, radius int) ([]*domain.SectorScanReport, error)
	Delete(id int) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) ScanReportRepository
}

// ResourceLocationRepository defines persistence for resource locations per sector.
type ResourceLocationRepository interface {
	Create(loc *domain.ResourceLocationModel) error
	FindByID(id int) (*domain.ResourceLocationModel, error)
	FindByCoordinates(x, y int) (*domain.ResourceLocationModel, error)
	// FindByCoordinatesForUpdate acquires a row-level lock on the resource location for the duration of the transaction.
	FindByCoordinatesForUpdate(x, y int) (*domain.ResourceLocationModel, error)
	// FindClosest returns the resource location closest to the provided coordinates.
	FindClosest(x, y int) (*domain.ResourceLocationModel, error)
	Update(loc *domain.ResourceLocationModel) error
	Delete(id int) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) ResourceLocationRepository
}

// DangerousLocationRepository defines persistence for dangerous locations per sector.
type DangerousLocationRepository interface {
	Create(loc *domain.DangerousLocationModel) error
	FindByID(id int) (*domain.DangerousLocationModel, error)
	FindByCoordinates(x, y int) (*domain.DangerousLocationModel, error)
	// FindByCoordinatesForUpdate acquires a row-level lock on the dangerous location for the duration of the transaction.
	FindByCoordinatesForUpdate(x, y int) (*domain.DangerousLocationModel, error)
	// FindClosest returns the dangerous location closest to the provided coordinates.
	FindClosest(x, y int) (*domain.DangerousLocationModel, error)
	Update(loc *domain.DangerousLocationModel) error
	Delete(id int) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) DangerousLocationRepository
}

// UserBaseRepository defines the interface for user base persistence.
type UserBaseRepository interface {
	Create(base *domain.UserBaseModel) error
	FindByID(id int) (*domain.UserBaseModel, error)
	// FindByIDForUpdate acquires a row-level lock on the base for the duration of the transaction.
	FindByIDForUpdate(id int) (*domain.UserBaseModel, error)
	// GetOwnerID returns just the owning user ID for the specified base, or ErrNotFound if missing.
	// Implementations should prefer a lightweight SELECT of user_id only.
	GetOwnerID(baseID int) (int, error)
	// Update replaces all per-item rows (army/build/tech/storage) with the current aggregate state and updates base stats.
	Update(base *domain.UserBaseModel) error
	Delete(id int) error
	FindByUserID(userID int) ([]*domain.UserBaseModel, error)
	FindByCoordinates(x, y int) (*domain.UserBaseModel, error)
	// FindByCoordinatesForUpdate acquires a row-level lock on the base found by coordinates.
	FindByCoordinatesForUpdate(x, y int) (*domain.UserBaseModel, error)
	// FindClosest returns the user base closest to the provided coordinates.
	FindClosest(x, y int) (*domain.UserBaseModel, error)
	// FindAll returns all user bases (light hydration). Primarily for world generation / balancing.
	FindAll() ([]*domain.UserBaseModel, error)
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) UserBaseRepository
}

// BuildPrototypeRepository provides access to building prototypes only.
// Base building queues and inventory are managed via the UserBase repository.
type BuildPrototypeRepository interface {
	// Prototypes
	CreatePrototype(proto *domain.BuildItemPrototype) error
	FindPrototypeByID(id int) (*domain.BuildItemPrototype, error)
	FindAllPrototypes() ([]*domain.BuildItemPrototype, error)
	UpdatePrototype(proto *domain.BuildItemPrototype) error
	DeletePrototype(id int) error

	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) BuildPrototypeRepository
}

// ArmyPrototypeRepository provides access to army item prototypes only.
// Base inventory and queues are managed through the UserBase aggregate and repository.
type ArmyPrototypeRepository interface {
	// Prototypes
	CreatePrototype(proto *domain.ArmyItemPrototype) error
	FindPrototypeByID(id int) (*domain.ArmyItemPrototype, error)
	FindAllPrototypes() ([]*domain.ArmyItemPrototype, error)
	UpdatePrototype(proto *domain.ArmyItemPrototype) error
	DeletePrototype(id int) error

	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) ArmyPrototypeRepository
}

// StoragePrototypeRepository provides access to storage item prototypes only.
// Inventory state is managed by the UserBase aggregate and repository.
type StoragePrototypeRepository interface {
	// Prototypes
	CreatePrototype(proto *domain.StorageItemPrototype) error
	FindPrototypeByID(id int) (*domain.StorageItemPrototype, error)
	FindAllPrototypes() ([]*domain.StorageItemPrototype, error)
	UpdatePrototype(proto *domain.StorageItemPrototype) error
	DeletePrototype(id int) error

	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) StoragePrototypeRepository
}

// TechPrototypeRepository provides access to technology prototypes only.
// Research queues and state are managed via the UserBase aggregate and repository.
type TechPrototypeRepository interface {
	// Prototypes
	CreatePrototype(proto *domain.TechItemPrototype) error
	FindPrototypeByID(id int) (*domain.TechItemPrototype, error)
	FindAllPrototypes() ([]*domain.TechItemPrototype, error)
	UpdatePrototype(proto *domain.TechItemPrototype) error
	DeletePrototype(id int) error

	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) TechPrototypeRepository
}

// RadarThreatRepository defines the interface for radar threat persistence.
type RadarThreatRepository interface {
	Create(threat *domain.RadarThreat) error
	Update(threat *domain.RadarThreat) error
	FindByID(id uuid.UUID) (*domain.RadarThreat, error)
	FindByOperationID(opID int) (*domain.RadarThreat, error)
	RadarThreatExists(ownerBaseID int, opID int) (bool, error)
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) RadarThreatRepository
}

// MilitaryOperationRepository defines persistence for military operations.
type MilitaryOperationRepository interface {
	Create(op *domain.MilitaryOperation) error
	FindByID(id int) (*domain.MilitaryOperation, error)
	// FindByIDForUpdate acquires a row-level lock on the operation for the duration of the transaction.
	FindByIDForUpdate(id int) (*domain.MilitaryOperation, error)
	Update(op *domain.MilitaryOperation) error
	Delete(id int) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) MilitaryOperationRepository
}

// ActivityRepository defines persistence for activity items (append-only feed).
type ActivityRepository interface {
	Create(item *domain.ActivityItem) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) ActivityRepository
}
