package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/game/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func actor(c *gin.Context) cqrs.Actor {
	if v, ok := c.Get("userID"); ok {
		if id, ok2 := v.(uuid.UUID); ok2 {
			return cqrs.Actor{UserID: id}
		}
	}
	return cqrs.Actor{UserID: uuid.Nil}
}

func getLocale(c *gin.Context) string {
	// Simple locale extraction from Accept-Language header
	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		return "en"
	}
	// Take first language if multiple are provided
	return strings.Split(lang, ",")[0]
}

// handleCoreErr handles common core layer errors and writes an appropriate HTTP response.
// It returns true if a response was written and the caller should return.
func handleCoreErr(c *gin.Context, tr ports.Translator, err error) bool {
	if err == nil {
		return false
	}

	locale := getLocale(c)

	// AppError: high-level application errors following auth service pattern
	var appErr cqrs.AppError
	if errors.As(err, &appErr) {
		status := http.StatusInternalServerError
		switch appErr.Kind {
		case cqrs.KindNotFound:
			status = http.StatusNotFound
		case cqrs.KindForbidden:
			status = http.StatusForbidden
		case cqrs.KindConflict:
			status = http.StatusConflict
		case cqrs.KindInvalidInput:
			status = http.StatusUnprocessableEntity
		}

		c.JSON(status, gin.H{"error": tr.T(locale, appErr.Code, appErr.Params)})
		return true
	}

	// Domain errors: 422 with domain-provided message.
	var domErr domain.Error
	if errors.As(err, &domErr) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": tr.T(locale, domErr.Key, domErr.Params)})
		return true
	}

	// Fallback: 500 with generic message.
	c.JSON(http.StatusInternalServerError, gin.H{"error": tr.T(locale, "error.application.internal_server_error", nil)})
	slog.ErrorContext(c.Request.Context(), "internal error occured", "request", c.Request.URL.Path, "error", err.Error())
	return true
}
