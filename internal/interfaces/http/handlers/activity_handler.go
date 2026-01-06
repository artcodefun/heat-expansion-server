package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	queries cqrs.ActivityQueries
}

func NewActivityHandler(queries cqrs.ActivityQueries) *ActivityHandler {
	return &ActivityHandler{queries: queries}
}

// ListOffense handles GET /bases/:baseId/activities/offense.
func (h *ActivityHandler) ListOffense(c *gin.Context) {
	var req dtos.OffenseActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	activities, err := h.queries.ListOffenseActivities(ctx, req.Uri.BaseID, req.Query.Subtype, req.Query.Limit)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities))
}

// ListDefense handles GET /bases/:baseId/activities/defense.
func (h *ActivityHandler) ListDefense(c *gin.Context) {
	var req dtos.DefenseActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	activities, err := h.queries.ListDefenseActivities(ctx, req.Uri.BaseID, req.Query.Subtype, req.Query.Limit)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities))
}

// ListScan handles GET /bases/:baseId/activities/scan.
func (h *ActivityHandler) ListScan(c *gin.Context) {
	var req dtos.ActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	activities, err := h.queries.ListScanActivities(ctx, req.Uri.BaseID, req.Query.Limit)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities))
}

// ListRadar handles GET /bases/:baseId/activities/radar.
func (h *ActivityHandler) ListRadar(c *gin.Context) {
	var req dtos.ActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	activities, err := h.queries.ListRadarActivities(ctx, req.Uri.BaseID, req.Query.Limit)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities))
}

// ListTrade handles GET /bases/:baseId/activities/trade.
func (h *ActivityHandler) ListTrade(c *gin.Context) {
	var req dtos.ActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	activities, err := h.queries.ListTradeActivities(ctx, req.Uri.BaseID, req.Query.Limit)
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities))
}
