package dtos

// ActivityListURI captures baseId from the URL.
type ActivityListURI struct {
	BaseID int `uri:"baseId" binding:"required,min=1"`
}

// ActivityListQuery captures query params for listing activities.
type ActivityListQuery struct {
	Limit int `form:"limit" binding:"omitempty,min=1"`
}
