package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type StorageHandler struct {
	queries    cqrs.StorageQueries
	commands   cqrs.StorageCommands
	translator ports.Translator
}

func NewStorageHandler(queries cqrs.StorageQueries, commands cqrs.StorageCommands, translator ports.Translator) *StorageHandler {
	return &StorageHandler{
		queries:    queries,
		commands:   commands,
		translator: translator,
	}
}

// ListPresent handles GET /bases/:baseId/storage/present.
func (h *StorageHandler) ListPresent(c *gin.Context) {
	var req dtos.StorageListRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	category := dtos.StorageCategoryFromDTO(req.Query.Category)
	items, err := h.queries.ListPresentStorageItems(c.Request.Context(), actor, req.Uri.BaseID, category)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.StorageItemsPresentFromReadModels(items, h.translator, getLocale(c)))
}

// DeleteItem handles DELETE /bases/:baseId/storage/items/:itemId.
func (h *StorageHandler) DeleteItem(c *gin.Context) {
	var req dtos.StorageItemRequest
	if !bindRequest(c, &req) {
		return
	}
	actor := actor(c)
	if err := h.commands.DeletePresentStorageItem(c.Request.Context(), actor, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, h.translator, err) {
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

	actor := actor(c)
	if err := h.commands.ActivateBuff(c.Request.Context(), actor, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, h.translator, err) {
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
	actor := actor(c)
	if err := h.commands.StartIntelDecryption(c.Request.Context(), actor, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, h.translator, err) {
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
	actor := actor(c)
	if err := h.commands.StartDamagedItemRestoration(c.Request.Context(), actor, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, h.translator, err) {
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
	actor := actor(c)
	if err := h.commands.ActivateArtifact(c.Request.Context(), actor, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, h.translator, err) {
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
	actor := actor(c)
	if err := h.commands.DeactivateArtifact(c.Request.Context(), actor, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, h.translator, err) {
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
	actor := actor(c)
	if err := h.commands.OpenConsumableBox(c.Request.Context(), actor, req.Uri.BaseID, req.Uri.ItemID.Uuid()); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}
