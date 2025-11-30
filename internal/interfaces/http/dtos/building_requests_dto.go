package dtos

// QueueBuildingRequest represents the payload to queue a building.
type QueueBuildingRequest struct {
	PrototypeID int `json:"prototype_id" binding:"required,min=1"`
}

// SpeedUpBuildingRequest represents the payload to speed up building production.
type SpeedUpBuildingRequest struct {
	UserID int `json:"user_id" binding:"required,min=1"`
}

// BuildingItemURI binds /bases/:baseId/buildings/(pending|present)/:itemId style routes.
type BuildingItemURI struct {
	BaseID int    `uri:"baseId" binding:"required,min=1"`
	ItemID string `uri:"itemId" binding:"required"`
}

// BuildingTaskURI binds /bases/:baseId/buildings/production/:taskId routes.
type BuildingTaskURI struct {
	BaseID int    `uri:"baseId" binding:"required,min=1"`
	TaskID string `uri:"taskId" binding:"required"`
}
