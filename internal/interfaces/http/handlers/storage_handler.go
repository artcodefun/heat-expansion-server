package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	queries  cqrs.StorageQueries
	commands cqrs.StorageCommands
}

func NewStorageHandler(queries cqrs.StorageQueries, commands cqrs.StorageCommands) *StorageHandler {
	return &StorageHandler{queries: queries, commands: commands}
}

// ListPresent handles GET /bases/:baseId/storage/present.
func (h *StorageHandler) ListPresent(c *gin.Context) {
	var req dtos.StorageListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	items, err := h.queries.ListPresentStorageItems(ctx, req.Uri.BaseID)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.StorageItemsPresentFromReadModels(items))
}

// DeleteItem handles DELETE /bases/:baseId/storage/items/:itemId.
func (h *StorageHandler) DeleteItem(c *gin.Context) {
	var req dtos.StorageItemRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.DeletePresentStorageItem(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusOK)
}

// ActivateBuff handles POST /bases/:baseId/storage/items/:itemId/activate.
func (h *StorageHandler) ActivateBuff(c *gin.Context) {
	var req dtos.StorageItemRequest
	if !bindRequest(c, &req) {
		return
	}

	ctx := commandCtx(c)
	if err := h.commands.ActivateBuff(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusOK)
}
