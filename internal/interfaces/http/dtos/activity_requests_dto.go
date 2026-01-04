package dtos

import (
	"strings"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
)

// activityListQuery captures query params for listing activities.
type activityListQuery struct {
	Limit int    `form:"limit,default=10" binding:"omitempty,min=1"`
	Kind  string `form:"kind" binding:"omitempty,activity_kind"`
}

// ActivityListRequest bundles the URI and query params for activity listing endpoints.
type ActivityListRequest = Request[BaseURI, activityListQuery, None]

// ActivityKindFromDTO normalizes a request kind string to the read-model type.
func ActivityKindFromDTO(value string) readmodels.ActivityKind {
	return readmodels.ActivityKind(strings.ToUpper(strings.TrimSpace(value)))
}

// IsValidActivityKind returns true if value matches one of the predefined ActivityKind constants.
func IsValidActivityKind(value string) bool {
	upper := strings.ToUpper(strings.TrimSpace(value))
	switch readmodels.ActivityKind(upper) {
	case readmodels.ActivityKindMilitary, readmodels.ActivityKindScan, readmodels.ActivityKindRadar, readmodels.ActivityKindTrade:
		return true
	case "":
		return true
	default:
		return false
	}
}
