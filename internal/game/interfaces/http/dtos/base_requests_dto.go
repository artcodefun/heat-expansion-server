package dtos

// BaseURI is a generic URI binding for routes that include a base identifier.
// It is used across multiple handlers (army, buildings, tech, storage, etc.).
type BaseURI struct {
	BaseID int `uri:"baseId" binding:"required,min=1"`
}

type GetBaseStatusRequest = Request[BaseURI, None, None]
