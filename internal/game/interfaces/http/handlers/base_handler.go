package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
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

	ctx := queryCtx(c)
	stats, err := h.queries.GetBaseStats(ctx, req.Uri.BaseID)
	if handleCoreErr(c, err) {
		return
	}

	resp := dtos.BaseResourcesFromReadModel(stats)
	c.JSON(http.StatusOK, resp)
}

// CreateBase handles POST /bases.
func (h *BaseHandler) CreateBase(c *gin.Context) {
	ctx := commandCtx(c)
	if err := h.commands.CreateBase(ctx, ctx.UserID); handleCoreErr(c, err) {
		return
	}

	c.Status(http.StatusCreated)
}

// ListUserBases handles GET /bases (list bases owned by the authenticated user).
func (h *BaseHandler) ListUserBases(c *gin.Context) {
	ctx := queryCtx(c)
	bases, err := h.queries.ListUserBases(ctx)
	if handleCoreErr(c, err) {
		return
	}
	resp := dtos.UserBasesFromReadModels(bases)
	c.JSON(http.StatusOK, resp)
}
