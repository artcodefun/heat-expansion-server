package ports

import (
	"context"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

// ErrNotFound is a generic sentinel for missing records.
var ErrNotFound = errors.New("record not found")

// UserRepository defines the interface for user persistence.
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	// FindByIDForUpdate acquires a row-level lock on the user for the duration of the transaction.
	FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) UserRepository
}

// SectorRepository defines the interface for sector persistence.
type SectorRepository interface {
	Create(ctx context.Context, sector *domain.SectorModel) error
	Update(ctx context.Context, sector *domain.SectorModel) error
	// Finds a sector by its coordinates (x, y)
	FindByCoordinates(ctx context.Context, x int, y int) (*domain.SectorModel, error)
	// Finds and locks a sector row by coordinates within a transaction (SELECT ... FOR UPDATE)
	FindByCoordinatesForUpdate(ctx context.Context, x int, y int) (*domain.SectorModel, error)
	// Returns all sectors
	FindAll(ctx context.Context) ([]*domain.SectorModel, error)
	// Lists coordinates of all non-empty occupied locations (bases, resourceful, dangerous)
	ListOccupiedCoordinates(ctx context.Context) ([]domain.Vector2i, error)
	// Derives occupant location type for a sector by coordinates
	GetLocationTypeByCoordinates(ctx context.Context, x int, y int) (domain.LocationType, error)
	// CountLocationsInRange returns the count of resourceful and dangerous locations within a circle radius.
	CountLocationsInRange(ctx context.Context, x, y, radius int) (resourceful int, dangerous int, err error)
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) SectorRepository
}

// ScanReportRepository defines persistence for sector scan reports.
type ScanReportRepository interface {
	Create(ctx context.Context, report *domain.SectorScanReport) error
	FindByID(ctx context.Context, id int) (*domain.SectorScanReport, error)
	FindByBaseAndCoordinates(ctx context.Context, baseID int, x int, y int) ([]*domain.SectorScanReport, error)
	// RecentReportExistsByScanner checks if a report was produced by a specific building within the last 'since' seconds.
	RecentReportExistsByScanner(ctx context.Context, scannerID uuid.UUID, since int64) (bool, error)
	Delete(ctx context.Context, id int) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) ScanReportRepository
}

// ResourceLocationRepository defines persistence for resource locations per sector.
type ResourceLocationRepository interface {
	Create(ctx context.Context, loc *domain.ResourceLocationModel) error
	FindByID(ctx context.Context, id int) (*domain.ResourceLocationModel, error)
	FindByCoordinates(ctx context.Context, x, y int) (*domain.ResourceLocationModel, error)
	// FindByCoordinatesForUpdate acquires a row-level lock on the resource location for the duration of the transaction.
	FindByCoordinatesForUpdate(ctx context.Context, x, y int) (*domain.ResourceLocationModel, error)
	// FindClosest returns the resource location closest to the provided coordinates.
	FindClosest(ctx context.Context, x, y int) (*domain.ResourceLocationModel, error)
	Update(ctx context.Context, loc *domain.ResourceLocationModel) error
	Delete(ctx context.Context, id int) error
	DeleteByCoordinates(ctx context.Context, x, y int) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) ResourceLocationRepository
}

// DangerousLocationRepository defines persistence for dangerous locations per sector.
type DangerousLocationRepository interface {
	Create(ctx context.Context, loc *domain.DangerousLocationModel) error
	FindByID(ctx context.Context, id int) (*domain.DangerousLocationModel, error)
	FindByCoordinates(ctx context.Context, x, y int) (*domain.DangerousLocationModel, error)
	// FindByCoordinatesForUpdate acquires a row-level lock on the dangerous location for the duration of the transaction.
	FindByCoordinatesForUpdate(ctx context.Context, x, y int) (*domain.DangerousLocationModel, error)
	// FindClosest returns the dangerous location closest to the provided coordinates.
	FindClosest(ctx context.Context, x, y int) (*domain.DangerousLocationModel, error)
	Update(ctx context.Context, loc *domain.DangerousLocationModel) error
	Delete(ctx context.Context, id int) error
	DeleteByCoordinates(ctx context.Context, x, y int) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) DangerousLocationRepository
}

// UserBaseRepository defines the interface for user base persistence.
type UserBaseRepository interface {
	Create(ctx context.Context, base *domain.UserBaseModel) error
	FindByID(ctx context.Context, id int) (*domain.UserBaseModel, error)
	// FindByIDForUpdate acquires a row-level lock on the base for the duration of the transaction.
	FindByIDForUpdate(ctx context.Context, id int) (*domain.UserBaseModel, error)
	// GetOwnerID returns just the owning user ID for the specified base, or ErrNotFound if missing.
	// Implementations should prefer a lightweight SELECT of user_id only.
	GetOwnerID(ctx context.Context, baseID int) (uuid.UUID, error)
	// Update replaces all per-item rows (army/build/tech/storage) with the current aggregate state and updates base stats.
	Update(ctx context.Context, base *domain.UserBaseModel) error
	Delete(ctx context.Context, id int) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.UserBaseModel, error)
	FindByCoordinates(ctx context.Context, x, y int) (*domain.UserBaseModel, error)
	// FindByCoordinatesForUpdate acquires a row-level lock on the base found by coordinates.
	FindByCoordinatesForUpdate(ctx context.Context, x, y int) (*domain.UserBaseModel, error)
	// FindClosest returns the user base closest to the provided coordinates.
	FindClosest(ctx context.Context, x, y int) (*domain.UserBaseModel, error)
	// FindAll returns all user bases (light hydration). Primarily for world generation / balancing.
	FindAll(ctx context.Context) ([]*domain.UserBaseModel, error)
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) UserBaseRepository
}

// BuildPrototypeRepository provides access to building prototypes only.
// Base building queues and inventory are managed via the UserBase repository.
type BuildPrototypeRepository interface {
	// Prototypes
	CreatePrototype(ctx context.Context, proto *domain.BuildItemPrototype) error
	FindPrototypeByID(ctx context.Context, id int) (*domain.BuildItemPrototype, error)
	FindAllPrototypes(ctx context.Context) ([]*domain.BuildItemPrototype, error)
	UpdatePrototype(ctx context.Context, proto *domain.BuildItemPrototype) error
	DeletePrototype(ctx context.Context, id int) error

	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) BuildPrototypeRepository
}

// ArmyPrototypeRepository provides access to army item prototypes only.
// Base inventory and queues are managed through the UserBase aggregate and repository.
type ArmyPrototypeRepository interface {
	// Prototypes
	CreatePrototype(ctx context.Context, proto *domain.ArmyItemPrototype) error
	FindPrototypeByID(ctx context.Context, id int) (*domain.ArmyItemPrototype, error)
	FindAllPrototypes(ctx context.Context) ([]*domain.ArmyItemPrototype, error)
	UpdatePrototype(ctx context.Context, proto *domain.ArmyItemPrototype) error
	DeletePrototype(ctx context.Context, id int) error

	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) ArmyPrototypeRepository
}

// StoragePrototypeRepository provides access to storage item prototypes only.
// Inventory state is managed by the UserBase aggregate and repository.
type StoragePrototypeRepository interface {
	// Prototypes
	CreatePrototype(ctx context.Context, proto *domain.StorageItemPrototype) error
	FindPrototypeByID(ctx context.Context, id int) (*domain.StorageItemPrototype, error)
	FindAllPrototypes(ctx context.Context) ([]*domain.StorageItemPrototype, error)
	UpdatePrototype(ctx context.Context, proto *domain.StorageItemPrototype) error
	DeletePrototype(ctx context.Context, id int) error

	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) StoragePrototypeRepository
}

// TechPrototypeRepository provides access to technology prototypes only.
// Research queues and state are managed via the UserBase aggregate and repository.
type TechPrototypeRepository interface {
	// Prototypes
	CreatePrototype(ctx context.Context, proto *domain.TechItemPrototype) error
	FindPrototypeByID(ctx context.Context, id int) (*domain.TechItemPrototype, error)
	FindAllPrototypes(ctx context.Context) ([]*domain.TechItemPrototype, error)
	UpdatePrototype(ctx context.Context, proto *domain.TechItemPrototype) error
	DeletePrototype(ctx context.Context, id int) error

	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) TechPrototypeRepository
}

// RadarThreatRepository defines the interface for radar threat persistence.
type RadarThreatRepository interface {
	Create(ctx context.Context, threat *domain.RadarThreat) error
	Update(ctx context.Context, threat *domain.RadarThreat) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.RadarThreat, error)
	FindByOperationIDForUpdate(ctx context.Context, opID int) (*domain.RadarThreat, error)
	RadarThreatExists(ctx context.Context, ownerBaseID int, opID int) (bool, error)
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) RadarThreatRepository
}

// MilitaryOperationRepository defines persistence for military operations.
type MilitaryOperationRepository interface {
	Create(ctx context.Context, op *domain.MilitaryOperation) error
	FindByID(ctx context.Context, id int) (*domain.MilitaryOperation, error)
	// FindByIDForUpdate acquires a row-level lock on the operation for the duration of the transaction.
	FindByIDForUpdate(ctx context.Context, id int) (*domain.MilitaryOperation, error)
	Update(ctx context.Context, op *domain.MilitaryOperation) error
	Delete(ctx context.Context, id int) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) MilitaryOperationRepository
}

// TradeOperationRepository defines persistence for trade operations.
type TradeOperationRepository interface {
	Create(ctx context.Context, op *domain.TradeOperation) error
	FindByID(ctx context.Context, id int) (*domain.TradeOperation, error)
	FindByIDForUpdate(ctx context.Context, id int) (*domain.TradeOperation, error)
	Update(ctx context.Context, op *domain.TradeOperation) error
	Delete(ctx context.Context, id int) error
	Tx(tx Transaction) TradeOperationRepository
}

// BlackMarketOfferRepository defines persistence for black market offers.
type BlackMarketOfferRepository interface {
	Create(ctx context.Context, offer *domain.BlackMarketOffer) error
	Update(ctx context.Context, offer *domain.BlackMarketOffer) error
	FindByID(ctx context.Context, id int64) (*domain.BlackMarketOffer, error)
	FindByIDForUpdate(ctx context.Context, id int64) (*domain.BlackMarketOffer, error)
	ListActiveLimitedOffers(ctx context.Context, now int64) ([]*domain.BlackMarketOffer, error)
	ListExpiredLimitedOffers(ctx context.Context, now int64) ([]*domain.BlackMarketOffer, error)
	Tx(tx Transaction) BlackMarketOfferRepository
}

// ActivityRepository defines persistence for activity items (append-only feed).
type ActivityRepository interface {
	Create(ctx context.Context, item *domain.ActivityItem) error
	ExistsForOperation(ctx context.Context, baseID int, kind domain.ActivityKind, opID int) (bool, error)
	ExistsForScanReport(ctx context.Context, reportID int) (bool, error)
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) ActivityRepository
}

type DiplomaticRelationshipRepository interface {
	Create(ctx context.Context, relationship *domain.DiplomaticRelationship) error
	Update(ctx context.Context, relationship *domain.DiplomaticRelationship) error
	FindBetweenUsers(ctx context.Context, userAID, userBID uuid.UUID) (*domain.DiplomaticRelationship, error)
	FindBetweenUsersForUpdate(ctx context.Context, userAID, userBID uuid.UUID) (*domain.DiplomaticRelationship, error)
	Tx(tx Transaction) DiplomaticRelationshipRepository
}

type DiplomaticMessageRepository interface {
	Create(ctx context.Context, message *domain.DiplomaticMessage) error
	ExistsByRequestAndContent(ctx context.Context, requestID uuid.UUID, content domain.TranslationKey) (bool, error)
	FindByRequestAndContent(ctx context.Context, requestID uuid.UUID, content domain.TranslationKey) (*domain.DiplomaticMessage, error)
	FindByID(ctx context.Context, id uuid.UUID) (*domain.DiplomaticMessage, error)
	MarkChatAsRead(ctx context.Context, receiverUserID, senderUserID uuid.UUID) error
	Tx(tx Transaction) DiplomaticMessageRepository
}

type DiplomaticRequestRepository interface {
	Create(ctx context.Context, request *domain.DiplomaticRequest) error
	Update(ctx context.Context, request *domain.DiplomaticRequest) error
	FindByID(ctx context.Context, id uuid.UUID) (*domain.DiplomaticRequest, error)
	FindByIDForUpdate(ctx context.Context, id uuid.UUID) (*domain.DiplomaticRequest, error)
	ExistsPendingByKind(ctx context.Context, userAID, userBID uuid.UUID, kind domain.DiplomaticRequestKind) (bool, error)
	Tx(tx Transaction) DiplomaticRequestRepository
}

// CrystalCreditsRepository tracks credited crystal purchases for idempotency.
type CrystalCreditsRepository interface {
	Insert(ctx context.Context, orderID uuid.UUID, userID uuid.UUID, crystals int, creditedAt int64) error
	Exists(ctx context.Context, orderID uuid.UUID) (bool, error)
	Tx(tx Transaction) CrystalCreditsRepository
}

// AlertRepository defines persistence for high-priority notifications.
type AlertRepository interface {
	Create(ctx context.Context, alert *domain.Alert) error
	ExistsForActivity(ctx context.Context, activityID uuid.UUID) (bool, error)
	MarkAllAsRead(ctx context.Context, userID uuid.UUID) error
	DeleteExpired(ctx context.Context, now int64) error
	// Tx returns a repository instance bound to the provided transaction.
	Tx(tx Transaction) AlertRepository
}
