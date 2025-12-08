package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type SectorHandler struct {
	queries cqrs.SectorQueries
}

func NewSectorHandler(queries cqrs.SectorQueries) *SectorHandler {
	return &SectorHandler{queries: queries}
}

// GetSector handles GET /sectors/:x/:y.
func (h *SectorHandler) GetSector(c *gin.Context) {
	var req dtos.SectorGetRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	sector, err := h.queries.GetSector(ctx, req.Uri.X, req.Uri.Y)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.SectorFromReadModel(sector))
}

// GetLatestScans handles GET /bases/:baseId/sectors/scans/latest.
func (h *SectorHandler) GetLatestScans(c *gin.Context) {
	var req dtos.SectorLatestScansRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	reports, err := h.queries.GetLatestScans(ctx, req.Uri.BaseID)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.SectorScanReportsFromReadModels(reports))
}

// GetScansNear handles GET /bases/:baseId/sectors/scans/near.
func (h *SectorHandler) GetScansNear(c *gin.Context) {
	var req dtos.SectorScansNearRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	reports, err := h.queries.GetScansNear(ctx, req.Uri.BaseID, req.Query.CenterX, req.Query.CenterY, req.Query.Radius)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.SectorScanReportsFromReadModels(reports))
}

// ListOccupiedCoordinates handles GET /map/occupied-coordinates.
func (h *SectorHandler) ListOccupiedCoordinates(c *gin.Context) {
	ctx := queryCtx(c)
	coords, err := h.queries.ListOccupiedCoordinates(ctx)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.Vector2iListFromReadModels(coords))
}

// ListSectorsInRadius handles GET /map/sectors.
func (h *SectorHandler) ListSectorsInRadius(c *gin.Context) {
	var req dtos.SectorRadiusOnlyRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	sectors, err := h.queries.ListSectorsInRadius(ctx, req.Query.CenterX, req.Query.CenterY, req.Query.Radius)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.SectorModelsFromReadModels(sectors))
}
