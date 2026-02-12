package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
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
	category := dtos.StorageCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListPresentStorageItems(ctx, req.Uri.BaseID, category)
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

// StartIntelDecryption handles POST /bases/:baseId/storage/items/:itemId/decrypt.
func (h *StorageHandler) StartIntelDecryption(c *gin.Context) {
	var req dtos.StorageItemRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.StartIntelDecryption(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusOK)
}

// StartDamagedItemRestoration handles POST /bases/:baseId/storage/items/:itemId/restore.
func (h *StorageHandler) StartDamagedItemRestoration(c *gin.Context) {
	var req dtos.StorageItemRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.StartDamagedItemRestoration(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusOK)
}

// ActivateArtifact handles POST /bases/:baseId/storage/items/:itemId/enable.
func (h *StorageHandler) ActivateArtifact(c *gin.Context) {
	var req dtos.StorageItemRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.ActivateArtifact(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusOK)
}

// DeactivateArtifact handles POST /bases/:baseId/storage/items/:itemId/disable.
func (h *StorageHandler) DeactivateArtifact(c *gin.Context) {
	var req dtos.StorageItemRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.DeactivateArtifact(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusOK)
}

// OpenBox handles POST /bases/:baseId/storage/items/:itemId/open.
func (h *StorageHandler) OpenBox(c *gin.Context) {
	var req dtos.StorageItemRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	if err := h.commands.OpenConsumableBox(ctx, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusOK)
}
