package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TechHandler struct {
	queries  cqrs.TechQueries
	commands cqrs.TechCommands
}

func NewTechHandler(queries cqrs.TechQueries, commands cqrs.TechCommands) *TechHandler {
	return &TechHandler{queries: queries, commands: commands}
}

func (h *TechHandler) ListNew(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	ctx := queryCtx(c)
	items, err := h.queries.ListNewTechItems(ctx, baseID)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsNewFromReadModels(items))
}

func (h *TechHandler) ListInProgress(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	ctx := queryCtx(c)
	items, err := h.queries.ListInResearchTechItems(ctx, baseID)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsInProgressFromReadModels(items))
}

func (h *TechHandler) ListDone(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	ctx := queryCtx(c)
	items, err := h.queries.ListDoneTechItems(ctx, baseID)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.TechItemsDoneFromReadModels(items))
}

func (h *TechHandler) Queue(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	var body dtos.QueueTechRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.StartTechResearch(ctx, baseID, body.PrototypeID); handleCQRS(c, err) {
		return
	}
	c.Status(http.StatusAccepted)
}

func (h *TechHandler) SpeedUpProduction(c *gin.Context) {
	var uri dtos.TechTaskURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameters"})
		return
	}
	parsed, err := uuid.Parse(uri.TaskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid taskId"})
		return
	}
	var body dtos.SpeedUpTechRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.SpeedUpTechResearchWithCrystals(ctx, uri.BaseID, ctx.UserID, parsed); handleCQRS(c, err) {
		return
	}
	c.Status(http.StatusOK)
}
