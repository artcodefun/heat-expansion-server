package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs/readmodels"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type ActivityHandler struct {
	queries cqrs.ActivityQueries
}

func NewActivityHandler(queries cqrs.ActivityQueries) *ActivityHandler {
	return &ActivityHandler{queries: queries}
}

// List handles GET /bases/:baseId/activities.
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
	var activities []*readmodels.ActivityItem
	var err error
	if req.Query.Kind == "" {
		activities, err = h.queries.ListActivities(ctx, req.Uri.BaseID, limit)
	} else {
		kind := dtos.ActivityKindFromDTO(req.Query.Kind)
		activities, err = h.queries.ListActivitiesByKind(ctx, req.Uri.BaseID, kind, limit)
	}
	if handleCoreErr(c, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.ActivityItemsFromReadModels(activities))
}
