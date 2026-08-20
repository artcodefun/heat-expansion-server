package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func actor(c *gin.Context) application.Actor {
	if v, ok := c.Get("accountID"); ok {
		if id, ok2 := v.(uuid.UUID); ok2 {
			return application.Actor{AccountID: id}
		}
	}
	return application.Actor{AccountID: uuid.Nil}
}

func getLocale(c *gin.Context) string {
	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		return "en"
	}
	// Simplified: take the first part of the header (e.g., "en-US,en;q=0.9" -> "en")
	parts := strings.Split(lang, ",")
	if len(parts) > 0 {
		localeParts := strings.Split(parts[0], ";")
		if len(localeParts) > 0 {
			fullLocale := strings.TrimSpace(localeParts[0])
			localeLang := strings.Split(fullLocale, "-")
			return strings.ToLower(localeLang[0])
		}
	}
	return "en"
}

func handleCoreErr(c *gin.Context, tr ports.Translator, err error) bool {
	if err == nil {
		return false
	}

	locale := getLocale(c)

	var appErr application.AppError
	if !errors.As(err, &appErr) {
		status := http.StatusInternalServerError
		c.JSON(status, gin.H{"error": tr.T(locale, "error.application.internal_server_error", nil)})
		slog.ErrorContext(c.Request.Context(), "internal error occurred", "request", c.Request.URL.Path, "error", err.Error())
		return true
	}

	status := http.StatusInternalServerError
	switch appErr.Kind {
	case application.KindNotFound:
		status = http.StatusNotFound
	case application.KindForbidden:
		status = http.StatusForbidden
	case application.KindConflict:
		status = http.StatusConflict
	case application.KindInvalidInput:
		status = http.StatusUnprocessableEntity
	}

	c.JSON(status, gin.H{
		"error": tr.T(locale, appErr.Code, appErr.Params),
	})
	return true
}

func bindRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		slog.WarnContext(c.Request.Context(), "request rejected; invalid input", "method", c.Request.Method, "path", c.Request.URL.Path, "source", "json", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}
