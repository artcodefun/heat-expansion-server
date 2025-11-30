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
	var uri dtos.SectorCoordinatesURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid coordinates"})
		return
	}
	ctx := queryCtx(c)
	sector, err := h.queries.GetSector(ctx, uri.X, uri.Y)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dtos.SectorFromReadModel(sector))
}

func (h *SectorHandler) GetLatestScans(c *gin.Context) {
	var uri dtos.BaseURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid baseId"})
		return
	}
	ctx := queryCtx(c)
	reports, err := h.queries.GetLatestScans(ctx, uri.BaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dtos.SectorScanReportsFromReadModels(reports))
}

func (h *SectorHandler) GetScansNear(c *gin.Context) {
	var uri dtos.BaseURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid baseId"})
		return
	}
	var query dtos.SectorRadiusQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
		return
	}
	ctx := queryCtx(c)
	reports, err := h.queries.GetScansNear(ctx, uri.BaseID, query.CenterX, query.CenterY, query.Radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dtos.SectorScanReportsFromReadModels(reports))
}

func (h *SectorHandler) ListOccupiedCoordinates(c *gin.Context) {
	ctx := queryCtx(c)
	coords, err := h.queries.ListOccupiedCoordinates(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dtos.Vector2iListFromReadModels(coords))
}

func (h *SectorHandler) ListSectorsInRadius(c *gin.Context) {
	var query dtos.SectorRadiusQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query"})
		return
	}
	ctx := queryCtx(c)
	sectors, err := h.queries.ListSectorsInRadius(ctx, query.CenterX, query.CenterY, query.Radius)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dtos.SectorModelsFromReadModels(sectors))
}
