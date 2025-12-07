package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
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
	var req dtos.BuildingListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)

	category := dtos.BuildCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListNewBuildItems(ctx, req.Uri.BaseID, category)
	if handleCoreErr(c, err) {
		return
	}

	resp := dtos.BuildItemsNewFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListPending handles GET /bases/:baseId/buildings/pending.
func (h *BuildingHandler) ListPending(c *gin.Context) {
	var req dtos.BuildingListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)

	category := dtos.BuildCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListPendingBuildItems(ctx, req.Uri.BaseID, category)
	if handleCoreErr(c, err) {
		return
	}

	resp := dtos.BuildItemsPendingFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListInProduction handles GET /bases/:baseId/buildings/in-production.
func (h *BuildingHandler) ListInProduction(c *gin.Context) {
	var req dtos.BuildingListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)

	category := dtos.BuildCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListInProductionBuildItems(ctx, req.Uri.BaseID, category)
	if handleCoreErr(c, err) {
		return
	}

	resp := dtos.BuildItemsInProductionFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListPresent handles GET /bases/:baseId/buildings/present.
func (h *BuildingHandler) ListPresent(c *gin.Context) {
	var req dtos.BuildingListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)

	category := dtos.BuildCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListPresentBuildItems(ctx, req.Uri.BaseID, category)
	if handleCoreErr(c, err) {
		return
	}

	resp := dtos.BuildItemsPresentFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// Queue handles POST /bases/:baseId/buildings/queue.
func (h *BuildingHandler) Queue(c *gin.Context) {
	var req dtos.BuildingQueueRequest
	if !bindRequest(c, &req) {
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.QueueBuilding(ctx, req.Uri.BaseID, req.Body.PrototypeID); handleCoreErr(c, err) {
		return
	}

	c.Status(http.StatusAccepted)
}

// SpeedUpProduction handles POST /bases/:baseId/buildings/production/:taskId/speed-up.
func (h *BuildingHandler) SpeedUpProduction(c *gin.Context) {
	var req dtos.BuildingSpeedUpRequest
	if !bindRequest(c, &req) {
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.SpeedUpProductionWithCrystals(ctx, req.Uri.BaseID, req.Uri.TaskID.Uuid()); handleCoreErr(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

// CancelPending handles POST /bases/:baseId/buildings/pending/:itemId/cancel.
func (h *BuildingHandler) CancelPending(c *gin.Context) {
	var req dtos.BuildingCancelRequest
	if !bindRequest(c, &req) {
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.CancelPendingBuilding(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

// DeletePresent handles DELETE /bases/:baseId/buildings/present/:itemId.
func (h *BuildingHandler) DeletePresent(c *gin.Context) {
	var req dtos.BuildingDeleteRequest
	if !bindRequest(c, &req) {
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.DeletePresentBuilding(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, err) {
		return
	}

	c.Status(http.StatusOK)
}
