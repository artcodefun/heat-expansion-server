package handlers

import (
	"net/http"
	"strconv"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	category := c.Query("category")
	ctx := queryCtx(c)

	items, err := h.queries.ListNewArmyItems(ctx, baseID, category)
	if handleCQRS(c, err) {
		return
	}

	resp := dtos.ArmyItemsNewFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListPending handles GET /bases/:baseId/army/pending.
func (h *ArmyHandler) ListPending(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	category := c.Query("category")
	ctx := queryCtx(c)

	items, err := h.queries.ListPendingArmyItems(ctx, baseID, category)
	if handleCQRS(c, err) {
		return
	}

	resp := dtos.ArmyItemsPendingFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListInProduction handles GET /bases/:baseId/army/in-production.
func (h *ArmyHandler) ListInProduction(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	category := c.Query("category")
	ctx := queryCtx(c)

	items, err := h.queries.ListInProductionArmyItems(ctx, baseID, category)
	if handleCQRS(c, err) {
		return
	}

	resp := dtos.ArmyItemsInProductionFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// ListPresent handles GET /bases/:baseId/army/present.
func (h *ArmyHandler) ListPresent(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	category := c.Query("category")
	ctx := queryCtx(c)

	items, err := h.queries.ListPresentArmyItems(ctx, baseID, category)
	if handleCQRS(c, err) {
		return
	}

	resp := dtos.ArmyItemsPresentFromReadModels(items)
	c.JSON(http.StatusOK, resp)
}

// Queue handles POST /bases/:baseId/army/queue.
func (h *ArmyHandler) Queue(c *gin.Context) {
	baseID, ok := baseIDFromCtx(c)
	if !ok {
		return
	}
	var body dtos.QueueArmyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.QueueArmy(ctx, baseID, body.PrototypeID, body.Count); handleCQRS(c, err) {
		return
	}

	c.Status(http.StatusAccepted)
}

// SpeedUpProduction handles POST /bases/:baseId/army/production/:taskId/speed-up.
func (h *ArmyHandler) SpeedUpProduction(c *gin.Context) {
	var uri dtos.ArmyTaskURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameters"})
		return
	}
	itemID, err := uuid.Parse(uri.TaskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid taskId"})
		return
	}
	var body dtos.SpeedUpArmyRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.SpeedUpArmyProductionWithCrystals(ctx, uri.BaseID, body.UserID, itemID); handleCQRS(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

// CancelPending handles POST /bases/:baseId/army/pending/:itemId/cancel.
func (h *ArmyHandler) CancelPending(c *gin.Context) {
	var uri dtos.ArmyItemURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameters"})
		return
	}
	itemID, err := uuid.Parse(uri.ItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid itemId"})
		return
	}

	countStr := c.Query("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid count"})
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.CancelPendingArmy(ctx, uri.BaseID, itemID, count); handleCQRS(c, err) {
		return
	}

	c.Status(http.StatusOK)
}

// DeletePresent handles DELETE /bases/:baseId/army/present/:itemId.
func (h *ArmyHandler) DeletePresent(c *gin.Context) {
	var uri dtos.ArmyItemURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid path parameters"})
		return
	}
	itemID, err := uuid.Parse(uri.ItemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid itemId"})
		return
	}

	countStr := c.Query("count")
	count, err := strconv.Atoi(countStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid count"})
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.DeletePresentArmy(ctx, uri.BaseID, itemID, count); handleCQRS(c, err) {
		return
	}

	c.Status(http.StatusOK)
}
