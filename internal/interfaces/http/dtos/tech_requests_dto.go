package dtos

import "strings"

// techListQuery contains query params for TechListRequest.
type techListQuery struct {
	Category string `form:"category" binding:"required,tech_category"`
}

// TechListRequest represents a base-scoped tech list request.
type TechListRequest = Request[BaseURI, techListQuery, None]

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

// IsValidTechCategory returns true if value matches one of the predefined TechCategory constants.
func IsValidTechCategory(value string) bool {
	upper := strings.ToUpper(value)
	switch TechCategory(upper) {
	case Army, Build, Base, Politics:
		return true
	default:
		return false
	}
}
