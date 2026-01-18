package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	queries  cqrs.AlertQueries
	commands cqrs.AlertCommands
}

func NewAlertHandler(queries cqrs.AlertQueries, commands cqrs.AlertCommands) *AlertHandler {
	return &AlertHandler{
		queries:  queries,
		commands: commands,
	}
}

// ListActive handles GET /bases/:baseId/alerts.
func (h *AlertHandler) ListActive(c *gin.Context) {
	var req dtos.AlertListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	alerts, err := h.queries.ListActiveAlerts(ctx, req.Uri.BaseID)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.AlertItemsFromReadModels(alerts))
}

// GetUnreadCount handles GET /bases/:baseId/alerts/unread-count.
func (h *AlertHandler) GetUnreadCount(c *gin.Context) {
	var req dtos.AlertListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	count, err := h.queries.GetUnreadAlertsCount(ctx, req.Uri.BaseID)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

// MarkAllAsRead handles POST /bases/:baseId/alerts/read-all.
func (h *AlertHandler) MarkAllAsRead(c *gin.Context) {
	var req dtos.AlertListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := commandCtx(c)
	err := h.commands.MarkAllAsRead(req.Uri.BaseID, ctx.UserID)
	if handleCoreErr(c, err) {
		return
	}
	c.Status(http.StatusNoContent)
}
