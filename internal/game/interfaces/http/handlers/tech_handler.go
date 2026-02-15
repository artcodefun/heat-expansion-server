package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type TechHandler struct {
	queries    cqrs.TechQueries
	commands   cqrs.TechCommands
	translator ports.Translator
}

func NewTechHandler(queries cqrs.TechQueries, commands cqrs.TechCommands, translator ports.Translator) *TechHandler {
	return &TechHandler{
		queries:    queries,
		commands:   commands,
		translator: translator,
	}
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
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsNewFromReadModels(items, h.translator, getLocale(c)))
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
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsInProgressFromReadModels(items, h.translator, getLocale(c)))
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
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsDoneFromReadModels(items, h.translator, getLocale(c)))
}

// Queue handles POST /bases/:baseId/tech/queue.
func (h *TechHandler) Queue(c *gin.Context) {
	var req dtos.TechQueueRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.StartTechResearch(ctx, req.Uri.BaseID, req.Body.PrototypeID); handleCoreErr(c, h.translator, err) {
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
	if err := h.commands.SpeedUpTechResearchWithCrystals(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}
