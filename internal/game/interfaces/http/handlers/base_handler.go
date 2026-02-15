package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
	queries    cqrs.BaseQueries
	commands   cqrs.BaseCommands
	translator ports.Translator
}

func NewBaseHandler(queries cqrs.BaseQueries, commands cqrs.BaseCommands, translator ports.Translator) *BaseHandler {
	return &BaseHandler{queries: queries, commands: commands, translator: translator}
}

// GetBaseStatus handles GET /bases/:baseId/status.
func (h *BaseHandler) GetBaseStatus(c *gin.Context) {
	var req dtos.GetBaseStatusRequest
	if !bindRequest(c, &req) {
		return
	}

	ctx := queryCtx(c)
	stats, err := h.queries.GetBaseStats(ctx, req.Uri.BaseID)
	if handleCoreErr(c, h.translator, err) {
		return
	}

	resp := dtos.BaseResourcesFromReadModel(stats, h.translator, getLocale(c))
	c.JSON(http.StatusOK, resp)
}

// CreateBase handles POST /bases.
func (h *BaseHandler) CreateBase(c *gin.Context) {
	ctx := commandCtx(c)
	if err := h.commands.CreateBase(ctx, ctx.UserID); handleCoreErr(c, h.translator, err) {
		return
	}

	c.Status(http.StatusCreated)
}

// ListUserBases handles GET /bases (list bases owned by the authenticated user).
func (h *BaseHandler) ListUserBases(c *gin.Context) {
	ctx := queryCtx(c)
	bases, err := h.queries.ListUserBases(ctx)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	resp := dtos.UserBasesFromReadModels(bases, h.translator, getLocale(c))
	c.JSON(http.StatusOK, resp)
}
