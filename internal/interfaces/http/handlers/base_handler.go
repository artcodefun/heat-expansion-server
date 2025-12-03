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
	var uri dtos.GetBaseStatusURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid baseId"})
		return
	}

	// TODO: derive user from auth middleware
	ctx := queryCtx(c)
	stats, err := h.queries.GetBaseStats(ctx, uri.BaseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := dtos.BaseResourcesFromReadModel(stats)
	c.JSON(http.StatusOK, resp)
}

// CreateBase handles POST /bases.
func (h *BaseHandler) CreateBase(c *gin.Context) {
	// TODO: derive user from auth middleware
	ctx := commandCtx(c)
	if err := h.commands.CreateBase(ctx, ctx.UserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// ListUserBases handles GET /bases (list bases owned by the authenticated user).
func (h *BaseHandler) ListUserBases(c *gin.Context) {
	ctx := queryCtx(c)
	bases, err := h.queries.ListUserBases(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := dtos.UserBasesFromReadModels(bases)
	c.JSON(http.StatusOK, resp)
}
