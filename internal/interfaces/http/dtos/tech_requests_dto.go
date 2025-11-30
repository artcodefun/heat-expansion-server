package dtos

// QueueTechRequest represents the payload to queue a technology research.
type QueueTechRequest struct {
	PrototypeID int `json:"prototype_id" binding:"required,min=1"`
}

// SpeedUpTechRequest represents the payload to speed up tech research.
type SpeedUpTechRequest struct {
	UserID int `json:"user_id" binding:"required,min=1"`
}

// TechTaskURI binds /bases/:baseId/tech/production/:taskId routes.
type TechTaskURI struct {
	BaseID int    `uri:"baseId" binding:"required,min=1"`
	TaskID string `uri:"taskId" binding:"required"`
}
