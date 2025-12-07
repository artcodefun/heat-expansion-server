package dtos

// TechListRequest represents a base-scoped tech list request.
type TechListRequest = Request[BaseURI, None, None]

// techQueueBody represents the JSON payload to queue a technology research.
type techQueueBody struct {
	PrototypeID int `json:"prototype_id" binding:"required,min=1"`
}

// TechQueueRequest bundles URI params with the queue payload.
type TechQueueRequest = Request[BaseURI, None, techQueueBody]

// techSpeedUpURI contains URI params for the tech speed-up endpoint.
type techSpeedUpURI struct {
	BaseURI
	ItemID UuidStr `uri:"itemId" binding:"required,uuid"`
}

// TechSpeedUpRequest bundles URI params for the tech speed-up endpoint.
type TechSpeedUpRequest = Request[techSpeedUpURI, None, None]
