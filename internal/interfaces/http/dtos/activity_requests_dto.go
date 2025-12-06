package dtos

// activityListQuery captures query params for listing activities.
type activityListQuery struct {
	Limit int `form:"limit" binding:"omitempty,min=1"`
}

// ActivityListRequest bundles the URI and query params for activity listing endpoints.
type ActivityListRequest = Request[BaseURI, activityListQuery, None]
