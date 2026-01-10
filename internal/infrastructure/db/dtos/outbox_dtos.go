package dtos

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/domain"
	"github.com/artcodefun/heat-expansion-api/internal/core/ports"
	"github.com/google/uuid"
)

// =========================
// Job DTOs
// =========================

type MoveBuildQueueJobDTO struct {
	BaseID int `json:"base_id"`
}

func MoveBuildQueueJobDTOFromDomain(j ports.MoveBuildQueueJob) MoveBuildQueueJobDTO {
	return MoveBuildQueueJobDTO{BaseID: j.BaseID}
}

func MoveBuildQueueJobFromDTO(d MoveBuildQueueJobDTO) ports.MoveBuildQueueJob {
	return ports.MoveBuildQueueJob{BaseID: d.BaseID}
}

type MoveArmyQueueJobDTO struct {
	BaseID int `json:"base_id"`
}

func MoveArmyQueueJobDTOFromDomain(j ports.MoveArmyQueueJob) MoveArmyQueueJobDTO {
	return MoveArmyQueueJobDTO{BaseID: j.BaseID}
}

func MoveArmyQueueJobFromDTO(d MoveArmyQueueJobDTO) ports.MoveArmyQueueJob {
	return ports.MoveArmyQueueJob{BaseID: d.BaseID}
}

type MoveTechQueueJobDTO struct {
	BaseID int `json:"base_id"`
}

func MoveTechQueueJobDTOFromDomain(j ports.MoveTechQueueJob) MoveTechQueueJobDTO {
	return MoveTechQueueJobDTO{BaseID: j.BaseID}
}

func MoveTechQueueJobFromDTO(d MoveTechQueueJobDTO) ports.MoveTechQueueJob {
	return ports.MoveTechQueueJob{BaseID: d.BaseID}
}

type DeleteExpiredBuffJobDTO struct {
	BaseID int       `json:"base_id"`
	ItemID uuid.UUID `json:"item_id"`
}

func DeleteExpiredBuffJobDTOFromDomain(j ports.DeleteExpiredBuffJob) DeleteExpiredBuffJobDTO {
	return DeleteExpiredBuffJobDTO{BaseID: j.BaseID, ItemID: j.ItemID}
}

func DeleteExpiredBuffJobFromDTO(d DeleteExpiredBuffJobDTO) ports.DeleteExpiredBuffJob {
	return ports.DeleteExpiredBuffJob{BaseID: d.BaseID, ItemID: d.ItemID}
}

type RestoreDamagedItemJobDTO struct {
	BaseID int       `json:"base_id"`
	ItemID uuid.UUID `json:"item_id"`
}

func RestoreDamagedItemJobDTOFromDomain(j ports.RestoreDamagedItemJob) RestoreDamagedItemJobDTO {
	return RestoreDamagedItemJobDTO{BaseID: j.BaseID, ItemID: j.ItemID}
}

func RestoreDamagedItemJobFromDTO(d RestoreDamagedItemJobDTO) ports.RestoreDamagedItemJob {
	return ports.RestoreDamagedItemJob{BaseID: d.BaseID, ItemID: d.ItemID}
}

type DecryptIntelItemJobDTO struct {
	BaseID int       `json:"base_id"`
	ItemID uuid.UUID `json:"item_id"`
}

func DecryptIntelItemJobDTOFromDomain(j ports.DecryptIntelItemJob) DecryptIntelItemJobDTO {
	return DecryptIntelItemJobDTO{BaseID: j.BaseID, ItemID: j.ItemID}
}

func DecryptIntelItemJobFromDTO(d DecryptIntelItemJobDTO) ports.DecryptIntelItemJob {
	return ports.DecryptIntelItemJob{BaseID: d.BaseID, ItemID: d.ItemID}
}

type UpdateMilitaryOperationJobDTO struct {
	OperationID int `json:"operation_id"`
}

func UpdateMilitaryOperationJobDTOFromDomain(j ports.UpdateMilitaryOperationJob) UpdateMilitaryOperationJobDTO {
	return UpdateMilitaryOperationJobDTO{OperationID: j.OperationID}
}

func UpdateMilitaryOperationJobFromDTO(d UpdateMilitaryOperationJobDTO) ports.UpdateMilitaryOperationJob {
	return ports.UpdateMilitaryOperationJob{OperationID: d.OperationID}
}

type SpawnNearbyLocationsJobDTO struct{}

func SpawnNearbyLocationsJobDTOFromDomain(j ports.SpawnNearbyLocationsJob) SpawnNearbyLocationsJobDTO {
	return SpawnNearbyLocationsJobDTO{}
}

func SpawnNearbyLocationsJobFromDTO(SpawnNearbyLocationsJobDTO) ports.SpawnNearbyLocationsJob {
	return ports.SpawnNearbyLocationsJob{}
}

type IntelligenceScanJobDTO struct {
	BaseID     int       `json:"base_id"`
	BuildingID uuid.UUID `json:"building_id"`
}

func IntelligenceScanJobDTOFromDomain(j ports.IntelligenceScanJob) IntelligenceScanJobDTO {
	return IntelligenceScanJobDTO{BaseID: j.BaseID, BuildingID: j.BuildingID}
}

func IntelligenceScanJobFromDTO(d IntelligenceScanJobDTO) ports.IntelligenceScanJob {
	return ports.IntelligenceScanJob{BaseID: d.BaseID, BuildingID: d.BuildingID}
}

type IntelligenceRadarJobDTO struct {
	BaseID      int `json:"base_id"`
	OperationID int `json:"operation_id"`
}

func IntelligenceRadarJobDTOFromDomain(j ports.IntelligenceRadarJob) IntelligenceRadarJobDTO {
	return IntelligenceRadarJobDTO{BaseID: j.BaseID, OperationID: j.OperationID}
}

func IntelligenceRadarJobFromDTO(d IntelligenceRadarJobDTO) ports.IntelligenceRadarJob {
	return ports.IntelligenceRadarJob{BaseID: d.BaseID, OperationID: d.OperationID}
}

// =========================
// Event DTOs
// =========================

// Each event DTO includes an OccurredAt field so the original event timestamp
// is captured in the payload. When decoding, we currently reconstruct domain
// events using their public constructors, which assign a fresh BasicEvent.

type UserAccountCreatedEventDTO struct {
	OccurredAt int64 `json:"occurred_at"`
	UserID     int   `json:"user_id"`
}

func UserAccountCreatedEventDTOFromDomain(e domain.UserAccountCreatedEvent) UserAccountCreatedEventDTO {
	return UserAccountCreatedEventDTO{OccurredAt: e.OccurredAt(), UserID: e.UserID}
}

func UserAccountCreatedEventFromDTO(d UserAccountCreatedEventDTO) domain.UserAccountCreatedEvent {
	return domain.NewUserAccountCreatedEvent(d.UserID)
}

type BuildingProductionStartedEventDTO struct {
	OccurredAt     int64     `json:"occurred_at"`
	BaseID         int       `json:"base_id"`
	ItemID         uuid.UUID `json:"item_id"`
	CompletionDate int64     `json:"completion_date"`
}

func BuildingProductionStartedEventDTOFromDomain(e domain.BuildingProductionStartedEvent) BuildingProductionStartedEventDTO {
	return BuildingProductionStartedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID, CompletionDate: e.CompletionDate}
}

func BuildingProductionStartedEventFromDTO(d BuildingProductionStartedEventDTO) domain.BuildingProductionStartedEvent {
	return domain.NewBuildingProductionStartedEvent(d.BaseID, d.ItemID, d.CompletionDate)
}

type BuildingProductionFinishedEventDTO struct {
	OccurredAt    int64     `json:"occurred_at"`
	BaseID        int       `json:"base_id"`
	ItemID        uuid.UUID `json:"item_id"`
	PresentItemID uuid.UUID `json:"present_item_id"`
}

func BuildingProductionFinishedEventDTOFromDomain(e domain.BuildingProductionFinishedEvent) BuildingProductionFinishedEventDTO {
	return BuildingProductionFinishedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID, PresentItemID: e.PresentItemID}
}

func BuildingProductionFinishedEventFromDTO(d BuildingProductionFinishedEventDTO) domain.BuildingProductionFinishedEvent {
	return domain.NewBuildingProductionFinishedEvent(d.BaseID, d.ItemID, d.PresentItemID)
}

type BuildingProductionCancelledEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func BuildingProductionCancelledEventDTOFromDomain(e domain.BuildingProductionCancelledEvent) BuildingProductionCancelledEventDTO {
	return BuildingProductionCancelledEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func BuildingProductionCancelledEventFromDTO(d BuildingProductionCancelledEventDTO) domain.BuildingProductionCancelledEvent {
	return domain.NewBuildingProductionCancelledEvent(d.BaseID, d.ItemID)
}

type BuildingProductionSpeedupEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func BuildingProductionSpeedupEventDTOFromDomain(e domain.BuildingProductionSpeedupEvent) BuildingProductionSpeedupEventDTO {
	return BuildingProductionSpeedupEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func BuildingProductionSpeedupEventFromDTO(d BuildingProductionSpeedupEventDTO) domain.BuildingProductionSpeedupEvent {
	return domain.NewBuildingProductionSpeedupEvent(d.BaseID, d.ItemID)
}

type BuildingPresentDeletedEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func BuildingPresentDeletedEventDTOFromDomain(e domain.BuildingPresentDeletedEvent) BuildingPresentDeletedEventDTO {
	return BuildingPresentDeletedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func BuildingPresentDeletedEventFromDTO(d BuildingPresentDeletedEventDTO) domain.BuildingPresentDeletedEvent {
	return domain.NewBuildingPresentDeletedEvent(d.BaseID, d.ItemID)
}

type ArmyProductionPendingEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
	Count      int       `json:"count"`
}

func ArmyProductionPendingEventDTOFromDomain(e domain.ArmyProductionPendingEvent) ArmyProductionPendingEventDTO {
	return ArmyProductionPendingEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID, Count: e.Count}
}

func ArmyProductionPendingEventFromDTO(d ArmyProductionPendingEventDTO) domain.ArmyProductionPendingEvent {
	return domain.NewArmyProductionPendingEvent(d.BaseID, d.ItemID, d.Count)
}

type ArmyProductionStartedEventDTO struct {
	OccurredAt     int64     `json:"occurred_at"`
	BaseID         int       `json:"base_id"`
	ItemID         uuid.UUID `json:"item_id"`
	CompletionDate int64     `json:"completion_date"`
}

func ArmyProductionStartedEventDTOFromDomain(e domain.ArmyProductionStartedEvent) ArmyProductionStartedEventDTO {
	return ArmyProductionStartedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID, CompletionDate: e.CompletionDate}
}

func ArmyProductionStartedEventFromDTO(d ArmyProductionStartedEventDTO) domain.ArmyProductionStartedEvent {
	return domain.NewArmyProductionStartedEvent(d.BaseID, d.ItemID, d.CompletionDate)
}

type ArmyProductionFinishedEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func ArmyProductionFinishedEventDTOFromDomain(e domain.ArmyProductionFinishedEvent) ArmyProductionFinishedEventDTO {
	return ArmyProductionFinishedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func ArmyProductionFinishedEventFromDTO(d ArmyProductionFinishedEventDTO) domain.ArmyProductionFinishedEvent {
	return domain.NewArmyProductionFinishedEvent(d.BaseID, d.ItemID)
}

type ArmyProductionCancelledEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
	Count      int       `json:"count"`
}

func ArmyProductionCancelledEventDTOFromDomain(e domain.ArmyProductionCancelledEvent) ArmyProductionCancelledEventDTO {
	return ArmyProductionCancelledEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID, Count: e.Count}
}

func ArmyProductionCancelledEventFromDTO(d ArmyProductionCancelledEventDTO) domain.ArmyProductionCancelledEvent {
	return domain.NewArmyProductionCancelledEvent(d.BaseID, d.ItemID, d.Count)
}

type ArmyProductionSpeedupEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func ArmyProductionSpeedupEventDTOFromDomain(e domain.ArmyProductionSpeedupEvent) ArmyProductionSpeedupEventDTO {
	return ArmyProductionSpeedupEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func ArmyProductionSpeedupEventFromDTO(d ArmyProductionSpeedupEventDTO) domain.ArmyProductionSpeedupEvent {
	return domain.NewArmyProductionSpeedupEvent(d.BaseID, d.ItemID)
}

type ArmyPresentDeletedEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
	Count      int       `json:"count"`
}

func ArmyPresentDeletedEventDTOFromDomain(e domain.ArmyPresentDeletedEvent) ArmyPresentDeletedEventDTO {
	return ArmyPresentDeletedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID, Count: e.Count}
}

func ArmyPresentDeletedEventFromDTO(d ArmyPresentDeletedEventDTO) domain.ArmyPresentDeletedEvent {
	return domain.NewArmyPresentDeletedEvent(d.BaseID, d.ItemID, d.Count)
}

type TechResearchStartedEventDTO struct {
	OccurredAt     int64     `json:"occurred_at"`
	BaseID         int       `json:"base_id"`
	ItemID         uuid.UUID `json:"item_id"`
	PrototypeID    int       `json:"prototype_id"`
	CompletionDate int64     `json:"completion_date"`
}

func TechResearchStartedEventDTOFromDomain(e domain.TechResearchStartedEvent) TechResearchStartedEventDTO {
	return TechResearchStartedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID, PrototypeID: e.PrototypeID, CompletionDate: e.CompletionDate}
}

func TechResearchStartedEventFromDTO(d TechResearchStartedEventDTO) domain.TechResearchStartedEvent {
	return domain.NewTechResearchStartedEvent(d.BaseID, d.ItemID, d.PrototypeID, d.CompletionDate)
}

type TechResearchFinishedEventDTO struct {
	OccurredAt  int64     `json:"occurred_at"`
	BaseID      int       `json:"base_id"`
	ItemID      uuid.UUID `json:"item_id"`
	PrototypeID int       `json:"prototype_id"`
}

func TechResearchFinishedEventDTOFromDomain(e domain.TechResearchFinishedEvent) TechResearchFinishedEventDTO {
	return TechResearchFinishedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID, PrototypeID: e.PrototypeID}
}

func TechResearchFinishedEventFromDTO(d TechResearchFinishedEventDTO) domain.TechResearchFinishedEvent {
	return domain.NewTechResearchFinishedEvent(d.BaseID, d.ItemID, d.PrototypeID)
}

type TechResearchSpeedupEventDTO struct {
	OccurredAt  int64     `json:"occurred_at"`
	BaseID      int       `json:"base_id"`
	ItemID      uuid.UUID `json:"item_id"`
	PrototypeID int       `json:"prototype_id"`
}

func TechResearchSpeedupEventDTOFromDomain(e domain.TechResearchSpeedupEvent) TechResearchSpeedupEventDTO {
	return TechResearchSpeedupEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID, PrototypeID: e.PrototypeID}
}

func TechResearchSpeedupEventFromDTO(d TechResearchSpeedupEventDTO) domain.TechResearchSpeedupEvent {
	return domain.NewTechResearchSpeedupEvent(d.BaseID, d.ItemID, d.PrototypeID)
}

type StorageItemPresentDeletedEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func StorageItemPresentDeletedEventDTOFromDomain(e domain.StorageItemPresentDeletedEvent) StorageItemPresentDeletedEventDTO {
	return StorageItemPresentDeletedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func StorageItemPresentDeletedEventFromDTO(d StorageItemPresentDeletedEventDTO) domain.StorageItemPresentDeletedEvent {
	return domain.NewStorageItemPresentDeletedEvent(d.BaseID, d.ItemID)
}

type BuffActivatedEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func BuffActivatedEventDTOFromDomain(e domain.BuffActivatedEvent) BuffActivatedEventDTO {
	return BuffActivatedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func BuffActivatedEventFromDTO(d BuffActivatedEventDTO) domain.BuffActivatedEvent {
	return domain.NewBuffActivatedEvent(d.BaseID, d.ItemID)
}

type IntelDecryptionStartedEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func IntelDecryptionStartedEventDTOFromDomain(e domain.IntelDecryptionStartedEvent) IntelDecryptionStartedEventDTO {
	return IntelDecryptionStartedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func IntelDecryptionStartedEventFromDTO(d IntelDecryptionStartedEventDTO) domain.IntelDecryptionStartedEvent {
	return domain.NewIntelDecryptionStartedEvent(d.BaseID, d.ItemID)
}

type IntelDecryptionFinishedEventDTO struct {
	OccurredAt int64                     `json:"occurred_at"`
	BaseID     int                       `json:"base_id"`
	ItemID     uuid.UUID                 `json:"item_id"`
	IntelType  domain.HiddenLocationType `json:"intel_type"`
}

func IntelDecryptionFinishedEventDTOFromDomain(e domain.IntelDecryptionFinishedEvent) IntelDecryptionFinishedEventDTO {
	return IntelDecryptionFinishedEventDTO{
		OccurredAt: e.OccurredAt(),
		BaseID:     e.BaseID,
		ItemID:     e.ItemID,
		IntelType:  e.IntelType,
	}
}

func IntelDecryptionFinishedEventFromDTO(d IntelDecryptionFinishedEventDTO) domain.IntelDecryptionFinishedEvent {
	return domain.NewIntelDecryptionFinishedEvent(d.BaseID, d.ItemID, d.IntelType)
}

type DamagedItemRestorationStartedEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func DamagedItemRestorationStartedEventDTOFromDomain(e domain.DamagedItemRestorationStartedEvent) DamagedItemRestorationStartedEventDTO {
	return DamagedItemRestorationStartedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func DamagedItemRestorationStartedEventFromDTO(d DamagedItemRestorationStartedEventDTO) domain.DamagedItemRestorationStartedEvent {
	return domain.NewDamagedItemRestorationStartedEvent(d.BaseID, d.ItemID)
}

type DamagedItemRestoredEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func DamagedItemRestoredEventDTOFromDomain(e domain.DamagedItemRestoredEvent) DamagedItemRestoredEventDTO {
	return DamagedItemRestoredEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func DamagedItemRestoredEventFromDTO(d DamagedItemRestoredEventDTO) domain.DamagedItemRestoredEvent {
	return domain.NewDamagedItemRestoredEvent(d.BaseID, d.ItemID)
}

type ArtifactActivatedEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func ArtifactActivatedEventDTOFromDomain(e domain.ArtifactActivatedEvent) ArtifactActivatedEventDTO {
	return ArtifactActivatedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func ArtifactActivatedEventFromDTO(d ArtifactActivatedEventDTO) domain.ArtifactActivatedEvent {
	return domain.NewArtifactActivatedEvent(d.BaseID, d.ItemID)
}

type ArtifactDeactivatedEventDTO struct {
	OccurredAt int64     `json:"occurred_at"`
	BaseID     int       `json:"base_id"`
	ItemID     uuid.UUID `json:"item_id"`
}

func ArtifactDeactivatedEventDTOFromDomain(e domain.ArtifactDeactivatedEvent) ArtifactDeactivatedEventDTO {
	return ArtifactDeactivatedEventDTO{OccurredAt: e.OccurredAt(), BaseID: e.BaseID, ItemID: e.ItemID}
}

func ArtifactDeactivatedEventFromDTO(d ArtifactDeactivatedEventDTO) domain.ArtifactDeactivatedEvent {
	return domain.NewArtifactDeactivatedEvent(d.BaseID, d.ItemID)
}

type MilitaryOperationStartedEventDTO struct {
	OccurredAt       int64 `json:"occurred_at"`
	OperationID      int   `json:"operation_id"`
	OutboundArriveAt int64 `json:"outbound_arrive_at"`
}

func MilitaryOperationStartedEventDTOFromDomain(e domain.MilitaryOperationStartedEvent) MilitaryOperationStartedEventDTO {
	return MilitaryOperationStartedEventDTO{OccurredAt: e.OccurredAt(), OperationID: e.OperationID, OutboundArriveAt: e.OutboundArriveAt}
}

func MilitaryOperationStartedEventFromDTO(d MilitaryOperationStartedEventDTO) domain.MilitaryOperationStartedEvent {
	return domain.NewMilitaryOperationStartedEvent(d.OperationID, d.OutboundArriveAt)
}

type MilitaryOperationArrivedEventDTO struct {
	OccurredAt  int64 `json:"occurred_at"`
	OperationID int   `json:"operation_id"`
}

func MilitaryOperationArrivedEventDTOFromDomain(e domain.MilitaryOperationArrivedEvent) MilitaryOperationArrivedEventDTO {
	return MilitaryOperationArrivedEventDTO{OccurredAt: e.OccurredAt(), OperationID: e.OperationID}
}

func MilitaryOperationArrivedEventFromDTO(d MilitaryOperationArrivedEventDTO) domain.MilitaryOperationArrivedEvent {
	return domain.NewMilitaryOperationArrivedEvent(d.OperationID)
}

type MilitaryOperationResolvedEventDTO struct {
	OccurredAt  int64                          `json:"occurred_at"`
	OperationID int                            `json:"operation_id"`
	Result      domain.MilitaryOperationResult `json:"result"`
}

func MilitaryOperationResolvedEventDTOFromDomain(e domain.MilitaryOperationResolvedEvent) MilitaryOperationResolvedEventDTO {
	return MilitaryOperationResolvedEventDTO{OccurredAt: e.OccurredAt(), OperationID: e.OperationID, Result: e.Result}
}

func MilitaryOperationResolvedEventFromDTO(d MilitaryOperationResolvedEventDTO) domain.MilitaryOperationResolvedEvent {
	return domain.NewMilitaryOperationResolvedEvent(d.OperationID, d.Result)
}

type MilitaryOperationReturnStartedEventDTO struct {
	OccurredAt     int64 `json:"occurred_at"`
	OperationID    int   `json:"operation_id"`
	ReturnArriveAt int64 `json:"return_arrive_at"`
}

func MilitaryOperationReturnStartedEventDTOFromDomain(e domain.MilitaryOperationReturnStartedEvent) MilitaryOperationReturnStartedEventDTO {
	return MilitaryOperationReturnStartedEventDTO{OccurredAt: e.OccurredAt(), OperationID: e.OperationID, ReturnArriveAt: e.ReturnArriveAt}
}

func MilitaryOperationReturnStartedEventFromDTO(d MilitaryOperationReturnStartedEventDTO) domain.MilitaryOperationReturnStartedEvent {
	return domain.NewMilitaryOperationReturnStartedEvent(d.OperationID, d.ReturnArriveAt)
}

type MilitaryOperationReturnArrivedEventDTO struct {
	OccurredAt  int64 `json:"occurred_at"`
	OperationID int   `json:"operation_id"`
}

func MilitaryOperationReturnArrivedEventDTOFromDomain(e domain.MilitaryOperationReturnArrivedEvent) MilitaryOperationReturnArrivedEventDTO {
	return MilitaryOperationReturnArrivedEventDTO{OccurredAt: e.OccurredAt(), OperationID: e.OperationID}
}

func MilitaryOperationReturnArrivedEventFromDTO(d MilitaryOperationReturnArrivedEventDTO) domain.MilitaryOperationReturnArrivedEvent {
	return domain.NewMilitaryOperationReturnArrivedEvent(d.OperationID)
}

type ScanReportCreatedEventDTO struct {
	OccurredAt        int64 `json:"occurred_at"`
	ReportID          int   `json:"report_id"`
	BaseID            int   `json:"base_id"`
	SourceOperationID int   `json:"source_operation_id"`
}

func ScanReportCreatedEventDTOFromDomain(e domain.ScanReportCreatedEvent) ScanReportCreatedEventDTO {
	return ScanReportCreatedEventDTO{OccurredAt: e.OccurredAt(), ReportID: e.ReportID, BaseID: e.BaseID, SourceOperationID: e.SourceOperationID}
}

func ScanReportCreatedEventFromDTO(d ScanReportCreatedEventDTO) domain.ScanReportCreatedEvent {
	return domain.NewScanReportCreatedEvent(d.ReportID, d.BaseID, d.SourceOperationID)
}
