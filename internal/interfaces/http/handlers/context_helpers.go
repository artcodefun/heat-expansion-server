package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/core/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

func queryCtx(c *gin.Context) cqrs.QueryContext {
	if v, ok := c.Get("userID"); ok {
		if id, ok2 := v.(int); ok2 {
			return cqrs.QueryContext{UserID: id}
		}
	}
	return cqrs.QueryContext{UserID: 0}
}

func commandCtx(c *gin.Context) cqrs.CommandContext {
	if v, ok := c.Get("userID"); ok {
		if id, ok2 := v.(int); ok2 {
			return cqrs.CommandContext{UserID: id}
		}
	}
	return cqrs.CommandContext{UserID: 0}
}

// baseIDFromCtx extracts a required baseId path parameter into an integer using the shared BaseURI DTO.
// It returns the parsed base ID and a boolean indicating success; on error it writes a 400 response.
func baseIDFromCtx(c *gin.Context) (int, bool) {
	var uri dtos.BaseURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid baseId"})
		return 0, false
	}
	return uri.BaseID, true
}
