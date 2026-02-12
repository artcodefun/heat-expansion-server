package handlers

import (
	"log/slog"
	"net/http"

	"github.com/artcodefun/heat-expansion-api/internal/game/core/cqrs"
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

// handleCoreErr handles common core layer errors and writes an appropriate HTTP response.
// It returns true if a response was written and the caller should return.
func handleCoreErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}
	// Domain errors: 400 with domain-provided message.
	if de, ok := err.(cqrs.DomainError); ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": de.Error()})
		return true
	}
	// Authorization failures.
	if err == cqrs.ErrForbidden {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return true
	}
	// Resource not found.
	if err == cqrs.ErrNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return true
	}
	// Fallback: 500 with generic message.
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	slog.Error("internal error occured", "request", c.Request, "error", err.Error())
	return true
}
