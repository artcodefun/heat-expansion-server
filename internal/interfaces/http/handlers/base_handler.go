package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	queries  cqrs.BaseQueries
	commands cqrs.BaseCommands
}

func NewBaseHandler(queries cqrs.BaseQueries, commands cqrs.BaseCommands) *BaseHandler {
	return &BaseHandler{queries: queries, commands: commands}
}

// GetBaseStatus handles GET /bases/:baseId/status.
func (h *BaseHandler) GetBaseStatus(c *gin.Context) {
	var req dtos.GetBaseStatusRequest
	if !bindRequest(c, &req) {
		return
	}

	// TODO: derive user from auth middleware
	ctx := queryCtx(c)
	stats, err := h.queries.GetBaseStats(ctx, req.Uri.BaseID)
	if handleCQRS(c, err) {
		return
	}

	resp := dtos.BaseResourcesFromReadModel(stats)
	c.JSON(http.StatusOK, resp)
}

// CreateBase handles POST /bases.
func (h *BaseHandler) CreateBase(c *gin.Context) {
	// TODO: derive user from auth middleware
	ctx := commandCtx(c)
	if err := h.commands.CreateBase(ctx, ctx.UserID); handleCQRS(c, err) {
		return
	}

	c.Status(http.StatusCreated)
}

// ListUserBases handles GET /bases (list bases owned by the authenticated user).
func (h *BaseHandler) ListUserBases(c *gin.Context) {
	ctx := queryCtx(c)
	bases, err := h.queries.ListUserBases(ctx)
	if handleCQRS(c, err) {
		return
	}
	resp := dtos.UserBasesFromReadModels(bases)
	c.JSON(http.StatusOK, resp)
}
