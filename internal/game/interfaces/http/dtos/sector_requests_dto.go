package dtos

// sectorRadiusQuery captures common radius-based query params.
type sectorRadiusQuery struct {
	CenterX int `form:"centerX" binding:"omitempty"`
	CenterY int `form:"centerY" binding:"omitempty"`
	Radius  int `form:"radius,default=10" binding:"omitempty,min=1,max=100"`
}

// SectorScansNearRequest binds both base URI and radius query params.
// Used for GET /bases/:baseId/sectors/scans/near.
type SectorScansNearRequest = Request[BaseURI, sectorRadiusQuery, None]

// sectorScanIDURI captures baseId and scan report ID from the URL.
type sectorScanIDURI struct {
	BaseID int `uri:"baseId" binding:"required"`
	ID     int `uri:"id" binding:"required"`
}

// SectorScanGetRequest bundles URI params for fetching a scan report by ID.
// Used for GET /bases/:baseId/sectors/scans/:id.
type SectorScanGetRequest = Request[sectorScanIDURI, None, None]

// sectorScanBeforeQuery captures query params for latest-before queries.
type sectorScanBeforeQuery struct {
	X      int   `form:"x" binding:"omitempty"`
	Y      int   `form:"y" binding:"omitempty"`
	Before int64 `form:"date" binding:"required"`
}

// SectorScanBeforeRequest bundles URI and query params for latest-before lookups.
// Used for GET /bases/:baseId/sectors/scans/before.
type SectorScanBeforeRequest = Request[BaseURI, sectorScanBeforeQuery, None]
