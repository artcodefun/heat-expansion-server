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

func (h *TechHandler) ListNew(c *gin.Context) {
	var req dtos.TechListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	items, err := h.queries.ListNewTechItems(ctx, req.Uri.BaseID)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsNewFromReadModels(items))
}

func (h *TechHandler) ListInProgress(c *gin.Context) {
	var req dtos.TechListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	items, err := h.queries.ListInResearchTechItems(ctx, req.Uri.BaseID)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsInProgressFromReadModels(items))
}

func (h *TechHandler) ListDone(c *gin.Context) {
	var req dtos.TechListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	items, err := h.queries.ListDoneTechItems(ctx, req.Uri.BaseID)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsDoneFromReadModels(items))
}

func (h *TechHandler) Queue(c *gin.Context) {
	var req dtos.TechQueueRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.StartTechResearch(ctx, req.Uri.BaseID, req.Body.PrototypeID); handleCQRS(c, err) {
		return
	}
	c.Status(http.StatusAccepted)
}

func (h *TechHandler) SpeedUpProduction(c *gin.Context) {
	var req dtos.TechSpeedUpRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.SpeedUpTechResearchWithCrystals(ctx, req.Uri.BaseID, req.Uri.TaskID.Uuid()); handleCQRS(c, err) {
		return
	}
	c.Status(http.StatusOK)
}
