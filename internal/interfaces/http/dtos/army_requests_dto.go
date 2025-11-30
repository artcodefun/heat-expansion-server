package dtos

// QueueArmyRequest represents the payload to queue army production.
type QueueArmyRequest struct {
	PrototypeID int `json:"prototype_id" binding:"required,min=1"`
	Count       int `json:"count" binding:"required,min=1"`
}

// SpeedUpArmyRequest represents the payload to speed up army production.
type SpeedUpArmyRequest struct {
	UserID int `json:"user_id" binding:"required,min=1"`
}

// ArmyItemURI binds /bases/:baseId/army/(pending|present)/:itemId style routes.
type ArmyItemURI struct {
	BaseID int    `uri:"baseId" binding:"required,min=1"`
	ItemID string `uri:"itemId" binding:"required"`
}

// ArmyTaskURI binds /bases/:baseId/army/production/:taskId routes.
type ArmyTaskURI struct {
	BaseID int    `uri:"baseId" binding:"required,min=1"`
	TaskID string `uri:"taskId" binding:"required"`
}
