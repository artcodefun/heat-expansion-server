package dtos

import (
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
)

// ListActivitiesQuery is the base query for activity listing.
type ListActivitiesQuery struct {
	Limit int `form:"limit,default=20" binding:"omitempty,min=1"`
}

// OffenseActivityListQuery adds subtype filtering for offensive activities.
type OffenseActivityListQuery struct {
	Limit   int                               `form:"limit,default=20" binding:"omitempty,min=1"`
	Subtype readmodels.OffenseActivitySubtype `form:"subtype" binding:"omitempty,oneof=ATTACK SPY"`
}

// DefenseActivityListQuery adds subtype filtering for defensive activities.
type DefenseActivityListQuery struct {
	Limit   int                               `form:"limit,default=20" binding:"omitempty,min=1"`
	Subtype readmodels.DefenseActivitySubtype `form:"subtype" binding:"omitempty,oneof=ATTACK SPY"`
}

// ScanActivityListQuery adds subtype filtering for scan activities.
type ScanActivityListQuery struct {
	Limit   int                            `form:"limit,default=20" binding:"omitempty,min=1"`
	Subtype readmodels.ScanActivitySubtype `form:"subtype" binding:"omitempty,oneof=REPORT_PRODUCED EXTERNAL_SCAN_DETECTED"`
}

// ActivityListRequest bundles the URI and query params for activity listing endpoints.
type ActivityListRequest = Request[BaseURI, ListActivitiesQuery, None]

// OffenseActivityListRequest bundles the URI and query params for offensive activity listing.
type OffenseActivityListRequest = Request[BaseURI, OffenseActivityListQuery, None]

// DefenseActivityListRequest bundles the URI and query params for defensive activity listing.
type DefenseActivityListRequest = Request[BaseURI, DefenseActivityListQuery, None]

// ScanActivityListRequest bundles the URI and query params for scan activity listing.
type ScanActivityListRequest = Request[BaseURI, ScanActivityListQuery, None]
