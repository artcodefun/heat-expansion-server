package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	queries    cqrs.ActivityQueries
	translator ports.Translator
}

func NewActivityHandler(queries cqrs.ActivityQueries, translator ports.Translator) *ActivityHandler {
	return &ActivityHandler{
		queries:    queries,
		translator: translator,
	}
}

// ListOffense handles GET /bases/:baseId/activities/offense.
func (h *ActivityHandler) ListOffense(c *gin.Context) {
	var req dtos.OffenseActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	activities, err := h.queries.ListOffenseActivities(ctx, req.Uri.BaseID, req.Query.Subtype, req.Query.Limit)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities, h.translator, getLocale(c)))
}

// ListDefense handles GET /bases/:baseId/activities/defense.
func (h *ActivityHandler) ListDefense(c *gin.Context) {
	var req dtos.DefenseActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	activities, err := h.queries.ListDefenseActivities(ctx, req.Uri.BaseID, req.Query.Subtype, req.Query.Limit)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities, h.translator, getLocale(c)))
}

// ListScan handles GET /bases/:baseId/activities/scan.
func (h *ActivityHandler) ListScan(c *gin.Context) {
	var req dtos.ScanActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	activities, err := h.queries.ListScanActivities(ctx, req.Uri.BaseID, req.Query.Subtype, req.Query.Limit)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities, h.translator, getLocale(c)))
}

// ListRadar handles GET /bases/:baseId/activities/radar.
func (h *ActivityHandler) ListRadar(c *gin.Context) {
	var req dtos.ActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	activities, err := h.queries.ListRadarActivities(ctx, req.Uri.BaseID, req.Query.Limit)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities, h.translator, getLocale(c)))
}

// ListTrade handles GET /bases/:baseId/activities/trade.
func (h *ActivityHandler) ListTrade(c *gin.Context) {
	var req dtos.ActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	activities, err := h.queries.ListTradeActivities(ctx, req.Uri.BaseID, req.Query.Limit)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities, h.translator, getLocale(c)))
}
