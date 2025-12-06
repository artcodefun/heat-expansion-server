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

func (h *SectorHandler) GetSector(c *gin.Context) {
	var req dtos.SectorGetRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	sector, err := h.queries.GetSector(ctx, req.Uri.X, req.Uri.Y)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.SectorFromReadModel(sector))
}

func (h *SectorHandler) GetLatestScans(c *gin.Context) {
	var req dtos.SectorLatestScansRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	reports, err := h.queries.GetLatestScans(ctx, req.Uri.BaseID)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.SectorScanReportsFromReadModels(reports))
}

func (h *SectorHandler) GetScansNear(c *gin.Context) {
	var req dtos.SectorScansNearRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	reports, err := h.queries.GetScansNear(ctx, req.Uri.BaseID, req.Query.CenterX, req.Query.CenterY, req.Query.Radius)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.SectorScanReportsFromReadModels(reports))
}

func (h *SectorHandler) ListOccupiedCoordinates(c *gin.Context) {
	ctx := queryCtx(c)
	coords, err := h.queries.ListOccupiedCoordinates(ctx)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.Vector2iListFromReadModels(coords))
}

func (h *SectorHandler) ListSectorsInRadius(c *gin.Context) {
	var req dtos.SectorRadiusOnlyRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	sectors, err := h.queries.ListSectorsInRadius(ctx, req.Query.CenterX, req.Query.CenterY, req.Query.Radius)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.SectorModelsFromReadModels(sectors))
}
