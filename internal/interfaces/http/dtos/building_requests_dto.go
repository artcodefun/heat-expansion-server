package dtos

import "strings"

// buildingListQuery contains query params for BuildingListRequest.
type buildingListQuery struct {
	Category string `form:"category,parser=encoding.TextUnmarshaler" binding:"required,build_category"`
}

// BuildingListRequest captures the path and query values for building list endpoints.
type BuildingListRequest = Request[BaseURI, buildingListQuery, None]

// buildingQueueBody represents the JSON payload required to queue a building.
type buildingQueueBody struct {
	PrototypeID int `json:"prototype_id" binding:"required,min=1"`
}

// BuildingQueueRequest binds the path and body for enqueueing new buildings.
type BuildingQueueRequest = Request[BaseURI, None, buildingQueueBody]

// buildingSpeedUpURI contains the URI params for speeding up production.
type buildingSpeedUpURI struct {
	BaseURI
	ItemID UuidStr `uri:"itemId" binding:"required,uuid"`
}

// BuildingSpeedUpRequest binds the path parameters for speeding up production.
type BuildingSpeedUpRequest = Request[buildingSpeedUpURI, None, None]

// buildingItemURI holds shared URI params for cancel/delete endpoints.
type buildingItemURI struct {
	BaseURI
	ItemID UuidStr `uri:"itemId" binding:"required,uuid"`
}

// BuildingCancelRequest binds the path parameters to cancel pending buildings.
type BuildingCancelRequest = Request[buildingItemURI, None, None]

// BuildingDeleteRequest binds the path parameters for deleting present buildings.
type BuildingDeleteRequest = Request[buildingItemURI, None, None]

// IsValidBuildCategory returns true when value matches one of the known BuildCategory constants.
func IsValidBuildCategory(value string) bool {
	upper := strings.ToUpper(value)
	switch BuildCategory(upper) {
	case Control, Resources, Defense, Military, Intelligence:
		return true
	default:
		return false
	}
}
