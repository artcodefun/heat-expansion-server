package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type ArmyHandler struct {
	queries  cqrs.ArmyQueries
	commands cqrs.ArmyCommands
}

func NewArmyHandler(queries cqrs.ArmyQueries, commands cqrs.ArmyCommands) *ArmyHandler {
	return &ArmyHandler{queries: queries, commands: commands}
}

// ListNew handles GET /bases/:baseId/army/new.
func (h *ArmyHandler) ListNew(c *gin.Context) {
	var req dtos.ArmyListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)

	category := dtos.ArmyCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListNewArmyItems(ctx, req.Uri.BaseID, category)
	if handleCoreErr(c, err) {
		return
	}

	resp := dtos.ArmyItemsNewFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListPending handles GET /bases/:baseId/army/pending.
func (h *ArmyHandler) ListPending(c *gin.Context) {
	var req dtos.ArmyListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)

	category := dtos.ArmyCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListPendingArmyItems(ctx, req.Uri.BaseID, category)
	if handleCoreErr(c, err) {
		return
	}

	resp := dtos.ArmyItemsPendingFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListInProduction handles GET /bases/:baseId/army/in-production.
func (h *ArmyHandler) ListInProduction(c *gin.Context) {
	var req dtos.ArmyListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)

	category := dtos.ArmyCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListInProductionArmyItems(ctx, req.Uri.BaseID, category)
	if handleCoreErr(c, err) {
		return
	}

	resp := dtos.ArmyItemsInProductionFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListPresent handles GET /bases/:baseId/army/present.
func (h *ArmyHandler) ListPresent(c *gin.Context) {
	var req dtos.ArmyListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)

	category := dtos.ArmyCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListPresentArmyItems(ctx, req.Uri.BaseID, category)
	if handleCoreErr(c, err) {
		return
	}

	resp := dtos.ArmyItemsPresentFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// Queue handles POST /bases/:baseId/army/queue.
func (h *ArmyHandler) Queue(c *gin.Context) {
	var req dtos.ArmyQueueRequest
	if !bindRequest(c, &req) {
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.QueueArmy(ctx, req.Uri.BaseID, req.Body.PrototypeID, req.Body.Count); handleCoreErr(c, err) {
		return
	}

	c.Status(http.StatusAccepted)
}

// SpeedUpProduction handles POST /bases/:baseId/army/production/:itemId/speed-up.
func (h *ArmyHandler) SpeedUpProduction(c *gin.Context) {
	var req dtos.ArmySpeedUpRequest
	if !bindRequest(c, &req) {
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.SpeedUpArmyProductionWithCrystals(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

// CancelPending handles POST /bases/:baseId/army/pending/:itemId/cancel.
func (h *ArmyHandler) CancelPending(c *gin.Context) {
	var req dtos.ArmyCancelRequest
	if !bindRequest(c, &req) {
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.CancelPendingArmy(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid(), req.Body.Count); handleCoreErr(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

// DeletePresent handles DELETE /bases/:baseId/army/present/:itemId.
func (h *ArmyHandler) DeletePresent(c *gin.Context) {
	var req dtos.ArmyDeleteRequest
	if !bindRequest(c, &req) {
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.DeletePresentArmy(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid(), req.Body.Count); handleCoreErr(c, err) {
		return
	}

	c.Status(http.StatusOK)
}
