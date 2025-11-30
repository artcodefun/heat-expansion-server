package dtos

type CreateOperationRequest struct {
	Type         OperationType `json:"type" binding:"required,oneof=ATTACK SPY OCCUPATION"`
	SourceBaseID int           `json:"source_base_id" binding:"required,min=1"`
	TargetX      int           `json:"target_x" binding:"required"`
	TargetY      int           `json:"target_y" binding:"required"`
	// Deployed contains army stacks to send, identified by prototype IDs.
	Deployed []ArmyDeploymentRequest `json:"deployed" binding:"required,dive"`
}

// ArmyDeploymentRequest represents a deployed unit in the create operation payload.
// It references an existing present army stack by its ID and a count to send.
type ArmyDeploymentRequest struct {
	PresentItemID string `json:"present_item_id" binding:"required,uuid4"`
	Count         int    `json:"count" binding:"required,min=1"`
}
