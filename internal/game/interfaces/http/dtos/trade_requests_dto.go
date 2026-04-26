package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/google/uuid"
)

type TradeBaseURI struct {
	BaseID int `uri:"baseId" binding:"required,min=1"`
}

type TradeBaseRequest = Request[TradeBaseURI, None, None]

type tradeOperationCreateBody struct {
	TargetX                 *int                       `json:"target_x" binding:"required"`
	TargetY                 *int                       `json:"target_y" binding:"required"`
	OfferedResources        PriceModelDTO              `json:"offered_resources" binding:"required"`
	OfferedArmy             []ArmyDeploymentRequestDTO `json:"offered_army" binding:"omitempty,dive"`
	OfferedStorageItemIDs   []UuidStr                  `json:"offered_storage_item_ids" binding:"omitempty,dive,uuid"`
	RequestedResources      PriceModelDTO              `json:"requested_resources" binding:"required"`
	RequestedArmy           []ArmyDeploymentRequestDTO `json:"requested_army" binding:"omitempty,dive"`
	RequestedStorageItemIDs []UuidStr                  `json:"requested_storage_item_ids" binding:"omitempty,dive,uuid"`
	TransportUnits          []ArmyDeploymentRequestDTO `json:"transport_units" binding:"omitempty,dive"`
}

type tradeInfoQuery struct {
	TargetX *int `form:"targetX" binding:"required"`
	TargetY *int `form:"targetY" binding:"required"`
}

type TradeInfoRequest = Request[TradeBaseURI, tradeInfoQuery, None]

type TradeOperationCreateRequest = Request[TradeBaseURI, None, tradeOperationCreateBody]

type tradeBaseOperationURI struct {
	BaseID      int `uri:"baseId" binding:"required,min=1"`
	OperationID int `uri:"operationId" binding:"required,min=1"`
}

type TradeBaseOperationIDRequest = Request[tradeBaseOperationURI, None, None]

func UUIDs(items []UuidStr) []uuid.UUID {
	out := make([]uuid.UUID, 0, len(items))
	for _, item := range items {
		out = append(out, item.Uuid())
	}
	return out
}

func ArmyDeploymentRequestsFromDTOs(items []ArmyDeploymentRequestDTO) []domain.ArmyDeploymentRequest {
	out := make([]domain.ArmyDeploymentRequest, 0, len(items))
	for _, item := range items {
		out = append(out, domain.ArmyDeploymentRequest{PresentItemID: item.PresentItemID.Uuid(), Count: item.Count})
	}
	return out
}
