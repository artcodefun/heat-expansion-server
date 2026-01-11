package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type TechHandler struct {
	queries  cqrs.TechQueries
	commands cqrs.TechCommands
}

func NewTechHandler(queries cqrs.TechQueries, commands cqrs.TechCommands) *TechHandler {
	return &TechHandler{queries: queries, commands: commands}
}

// ListNew handles GET /bases/:baseId/tech/new.
func (h *TechHandler) ListNew(c *gin.Context) {
	var req dtos.TechListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	category := dtos.TechCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListNewTechItems(ctx, req.Uri.BaseID, category)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsNewFromReadModels(items))
}

// ListInProgress handles GET /bases/:baseId/tech/in-progress.
func (h *TechHandler) ListInProgress(c *gin.Context) {
	var req dtos.TechListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	category := dtos.TechCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListInResearchTechItems(ctx, req.Uri.BaseID, category)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsInProgressFromReadModels(items))
}

// ListDone handles GET /bases/:baseId/tech/done.
func (h *TechHandler) ListDone(c *gin.Context) {
	var req dtos.TechListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	category := dtos.TechCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListDoneTechItems(ctx, req.Uri.BaseID, category)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsDoneFromReadModels(items))
}

// Queue handles POST /bases/:baseId/tech/queue.
func (h *TechHandler) Queue(c *gin.Context) {
	var req dtos.TechQueueRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.StartTechResearch(ctx, req.Uri.BaseID, req.Body.PrototypeID); handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusAccepted)
}

// SpeedUpProduction handles POST /bases/:baseId/tech/production/:itemId/speed-up.
func (h *TechHandler) SpeedUpProduction(c *gin.Context) {
	var req dtos.TechSpeedUpRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.SpeedUpTechResearchWithCrystals(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusOK)
}
