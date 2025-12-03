package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type BuildingHandler struct {
	queries  cqrs.BuildingQueries
	commands cqrs.BuildingCommands
}

func NewBuildingHandler(queries cqrs.BuildingQueries, commands cqrs.BuildingCommands) *BuildingHandler {
	return &BuildingHandler{queries: queries, commands: commands}
}

// ListNew handles GET /bases/:baseId/buildings/new.
func (h *BuildingHandler) ListNew(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	category := c.Query("category")
	ctx := queryCtx(c)

	items, err := h.queries.ListNewBuildItems(ctx, baseID, category)
	if handleCQRS(c, err) {
		return
	}

	resp := dtos.BuildItemsNewFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListPending handles GET /bases/:baseId/buildings/pending.
func (h *BuildingHandler) ListPending(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	category := c.Query("category")
	ctx := queryCtx(c)

	items, err := h.queries.ListPendingBuildItems(ctx, baseID, category)
	if handleCQRS(c, err) {
		return
	}

	resp := dtos.BuildItemsPendingFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListInProduction handles GET /bases/:baseId/buildings/in-production.
func (h *BuildingHandler) ListInProduction(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	category := c.Query("category")
	ctx := queryCtx(c)

	items, err := h.queries.ListInProductionBuildItems(ctx, baseID, category)
	if handleCQRS(c, err) {
		return
	}

	resp := dtos.BuildItemsInProductionFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListPresent handles GET /bases/:baseId/buildings/present.
func (h *BuildingHandler) ListPresent(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	category := c.Query("category")
	ctx := queryCtx(c)

	items, err := h.queries.ListPresentBuildItems(ctx, baseID, category)
	if handleCQRS(c, err) {
		return
	}

	resp := dtos.BuildItemsPresentFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// Queue handles POST /bases/:baseId/buildings/queue.
func (h *BuildingHandler) Queue(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	var body dtos.QueueBuildingRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.QueueBuilding(ctx, baseID, body.PrototypeID); handleCQRS(c, err) {
		return
	}

	c.Status(http.StatusAccepted)
}

// SpeedUpProduction handles POST /bases/:baseId/buildings/production/:taskId/speed-up.
func (h *BuildingHandler) SpeedUpProduction(c *gin.Context) {
	var uri dtos.BuildingTaskURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameters"})
		return
	}
	itemID, err := uuid.Parse(uri.TaskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid taskId"})
		return
	}
	var body dtos.SpeedUpBuildingRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.SpeedUpProductionWithCrystals(ctx, uri.BaseID, ctx.UserID, itemID); handleCQRS(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

// CancelPending handles POST /bases/:baseId/buildings/pending/:itemId/cancel.
func (h *BuildingHandler) CancelPending(c *gin.Context) {
	var uri dtos.BuildingItemURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameters"})
		return
	}
	itemID, err := uuid.Parse(uri.ItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid itemId"})
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.CancelPendingBuilding(ctx, uri.BaseID, itemID); handleCQRS(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

// DeletePresent handles DELETE /bases/:baseId/buildings/present/:itemId.
func (h *BuildingHandler) DeletePresent(c *gin.Context) {
	var uri dtos.BuildingItemURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameters"})
		return
	}
	itemID, err := uuid.Parse(uri.ItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid itemId"})
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.DeletePresentBuilding(ctx, uri.BaseID, itemID); handleCQRS(c, err) {
		return
	}

	c.Status(http.StatusOK)
}
