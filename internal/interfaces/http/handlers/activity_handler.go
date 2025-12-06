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

func (h *ActivityHandler) List(c *gin.Context) {
	var req dtos.ActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	limit := req.Query.Limit
	if limit <= 0 {
		limit = 20
	}
	activities, err := h.queries.ListActivities(ctx, req.Uri.BaseID, limit)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities))
}

func (h *ActivityHandler) ListMilitary(c *gin.Context) {
	var req dtos.ActivityListRequest
	if !bindRequest(c, &req) {
		return
	}
	ctx := queryCtx(c)
	limit := req.Query.Limit
	if limit <= 0 {
		limit = 20
	}
	activities, err := h.queries.ListMilitaryActivities(ctx, req.Uri.BaseID, limit)
	if handleCQRS(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities))
}
