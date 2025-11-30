package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StorageHandler struct {
	queries  cqrs.StorageQueries
	commands cqrs.StorageCommands
}

func NewStorageHandler(queries cqrs.StorageQueries, commands cqrs.StorageCommands) *StorageHandler {
	return &StorageHandler{queries: queries, commands: commands}
}

func (h *StorageHandler) ListPresent(c *gin.Context) {
	var uri dtos.BaseURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid baseId"})
		return
	}
	baseID := uri.BaseID
	ctx := queryCtx(c)
	items, err := h.queries.ListPresentStorageItems(ctx, baseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dtos.StorageItemsPresentFromReadModels(items))
}

func (h *StorageHandler) DeleteItem(c *gin.Context) {
	var uri dtos.StorageItemURI
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
	if err := h.commands.DeletePresentStorageItem(ctx, uri.BaseID, itemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *StorageHandler) ActivateBuff(c *gin.Context) {
	var uri dtos.StorageItemURI
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
	if err := h.commands.ActivateBuff(ctx, uri.BaseID, itemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
