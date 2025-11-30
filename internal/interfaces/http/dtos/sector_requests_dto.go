package dtos

// SectorCoordinatesURI captures sector coordinates from the URL.
type SectorCoordinatesURI struct {
	X int `uri:"x" binding:"required"`
	Y int `uri:"y" binding:"required"`
}

// SectorBaseURI captures a baseId from the URL.
type SectorBaseURI struct {
	BaseID int `uri:"baseId" binding:"required,min=1"`
}

// SectorRadiusQuery captures common radius-based query params.
type SectorRadiusQuery struct {
	CenterX int `form:"centerX" binding:"omitempty"`
	CenterY int `form:"centerY" binding:"omitempty"`
	Radius  int `form:"radius" binding:"omitempty,min=0"`
}
