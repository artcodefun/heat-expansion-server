package dtos

// sectorCoordinatesURI captures sector coordinates from the URL.
type sectorCoordinatesURI struct {
	X int `uri:"x" binding:"required"`
	Y int `uri:"y" binding:"required"`
}

// SectorGetRequest bundles URI params for sector lookups.
type SectorGetRequest = Request[sectorCoordinatesURI, None, None]

// SectorLatestScansRequest bundles baseId URI params for sector-specific endpoints.
type SectorLatestScansRequest = Request[BaseURI, None, None]

// sectorRadiusQuery captures common radius-based query params.
type sectorRadiusQuery struct {
	CenterX int `form:"centerX" binding:"omitempty"`
	CenterY int `form:"centerY" binding:"omitempty"`
	Radius  int `form:"radius" binding:"omitempty,min=0"`
}

// SectorRadiusOnlyRequest binds only the radius query params.
type SectorRadiusOnlyRequest = Request[None, sectorRadiusQuery, None]

// SectorScansNearRequest binds both base URI and radius query params.
type SectorScansNearRequest = Request[BaseURI, sectorRadiusQuery, None]
