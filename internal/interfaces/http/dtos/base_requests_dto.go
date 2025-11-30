package dtos

// GetBaseStatusURI captures the base ID from the URL.
type GetBaseStatusURI struct {
	BaseID int `uri:"baseId" binding:"required,min=1"`
}

// BaseURI is a generic URI binding for routes that include a base identifier.
// It is used across multiple handlers (army, buildings, tech, storage, etc.).
type BaseURI struct {
	BaseID int `uri:"baseId" binding:"required,min=1"`
}
