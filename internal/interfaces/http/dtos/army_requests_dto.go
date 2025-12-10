package dtos

import "strings"

// armyListQuery contains query params for ArmyListRequest.
type armyListQuery struct {
	Category string `form:"category" binding:"required,army_category"`
}

// ArmyListRequest aggregates URI and query params for army list endpoints.
type ArmyListRequest = Request[BaseURI, armyListQuery, None]

// armyQueueBody represents the JSON payload for queuing armies.
type armyQueueBody struct {
	PrototypeID int `json:"prototype_id" binding:"required,min=1"`
	Count       int `json:"count" binding:"required,min=1"`
}

// ArmyQueueRequest bundles URI params with the queue request body.
type ArmyQueueRequest = Request[BaseURI, None, armyQueueBody]

// armySpeedUpURI contains URI params for speeding up an army item in production.
type armySpeedUpURI struct {
	BaseURI
	ItemID UuidStr `uri:"itemId" binding:"required,uuid"`
}

// ArmySpeedUpRequest bundles URI params for the speed-up endpoint.
type ArmySpeedUpRequest = Request[armySpeedUpURI, None, None]

// armyCancelURI contains the URI parts of the cancel endpoint.
type armyCancelURI struct {
	BaseURI
	ItemID UuidStr `uri:"itemId" binding:"required,uuid"`
}

// armyCancelBody represents the JSON payload for cancelling pending army items.
type armyCancelBody struct {
	Count int `json:"count" binding:"required,min=1"`
}

// ArmyCancelRequest bundles URI params with a JSON body containing count.
type ArmyCancelRequest = Request[armyCancelURI, None, armyCancelBody]

// armyDeleteURI contains the URI parts of the delete endpoint.
type armyDeleteURI struct {
	BaseURI
	ItemID UuidStr `uri:"itemId" binding:"required,uuid"`
}

// armyDeleteBody represents the JSON payload for deleting present army items.
type armyDeleteBody struct {
	Count int `json:"count" binding:"required,min=1"`
}

// ArmyDeleteRequest bundles URI params with a JSON body containing count.
type ArmyDeleteRequest = Request[armyDeleteURI, None, armyDeleteBody]

// IsValidArmyCategory returns true if value matches one of the predefined ArmyCategory constants.
func IsValidArmyCategory(value string) bool {
	upper := strings.ToUpper(value)
	switch ArmyCategory(upper) {
	case Infantry, Armored, Artillery, Aviation, Spy, Special:
		return true
	default:
		return false
	}
}
