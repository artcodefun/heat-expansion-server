package dtos

import "strings"

type operationURI struct {
	BaseURI
	OperationID int `uri:"operationId" binding:"required,min=1"`
}

// OperationGetRequest binds the operation ID path parameter.
type OperationGetRequest = Request[operationURI, None, None]

// OperationSpeedUpRequest binds the operation ID path parameter for speed-up.
type OperationSpeedUpRequest = Request[operationURI, None, None]

// OperationCancelRequest binds the operation ID path parameter for cancellation.
type OperationCancelRequest = Request[operationURI, None, None]

// OperationByBaseRequest binds a baseId path parameter for listing operations.
// Used for GET /bases/:baseId/operations.
type OperationByBaseRequest = Request[BaseURI, None, None]

// OperationActiveByBaseRequest binds a baseId path parameter for listing active operations.
// Used for GET /bases/:baseId/operations/active.
type OperationActiveByBaseRequest = Request[BaseURI, None, None]

// ArmyDeploymentRequestDTO represents a deployed unit in the create operation payload.
// It references an existing present army stack by its ID and a count to send.
type ArmyDeploymentRequestDTO struct {
	PresentItemID UuidStr `json:"present_item_id" binding:"required,uuid"`
	Count         int     `json:"count" binding:"required,min=1"`
}

// operationCreateBody is the JSON payload for creating a military operation.
type operationCreateBody struct {
	Type    OperationType `json:"type" binding:"required,operation_type"`
	TargetX *int          `json:"target_x" binding:"required"`
	TargetY *int          `json:"target_y" binding:"required"`
	// Deployed contains army stacks to send, identified by prototype IDs.
	Deployed []ArmyDeploymentRequestDTO `json:"deployed" binding:"required,dive"`
}

// OperationCreateRequest binds the create operation body payload and source base URI.
type OperationCreateRequest = Request[BaseURI, None, operationCreateBody]

// IsValidOperationType returns true when value matches one of the known
// OperationType constants. Comparison is case-insensitive.
func IsValidOperationType(value string) bool {
	upper := strings.ToUpper(value)
	switch OperationType(upper) {
	case OperationTypeAttack, OperationTypeSpy, OperationTypeOccupation:
		return true
	default:
		return false
	}
}
