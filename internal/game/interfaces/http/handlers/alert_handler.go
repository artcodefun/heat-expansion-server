package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type AlertHandler struct {
	queries    cqrs.AlertQueries
	commands   cqrs.AlertCommands
	translator ports.Translator
}

func NewAlertHandler(queries cqrs.AlertQueries, commands cqrs.AlertCommands, translator ports.Translator) *AlertHandler {
	return &AlertHandler{
		queries:    queries,
		commands:   commands,
		translator: translator,
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
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.AlertItemsFromReadModels(alerts, h.translator, getLocale(c)))
}

// GetUnreadCount handles GET /bases/:baseId/alerts/unread-count.
func (h *AlertHandler) GetUnreadCount(c *gin.Context) {
	var req dtos.AlertListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	count, err := h.queries.GetUnreadAlertsCount(ctx, req.Uri.BaseID)
	if handleCoreErr(c, h.translator, err) {
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
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusNoContent)
}
