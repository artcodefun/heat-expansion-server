package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/auth/application/ports"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func queryCtx(c *gin.Context) cqrs.QueryContext {
	if v, ok := c.Get("accountID"); ok {
		if id, ok2 := v.(uuid.UUID); ok2 {
			return cqrs.QueryContext{AccountID: id}
		}
	}
	return cqrs.QueryContext{AccountID: uuid.Nil}
}

func commandCtx(c *gin.Context) cqrs.CommandContext {
	if v, ok := c.Get("accountID"); ok {
		if id, ok2 := v.(uuid.UUID); ok2 {
			return cqrs.CommandContext{AccountID: id}
		}
	}
	return cqrs.CommandContext{AccountID: uuid.Nil}
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

	var appErr cqrs.AppError
	if !errors.As(err, &appErr) {
		status := http.StatusInternalServerError
		c.JSON(status, gin.H{"error": tr.T(locale, "error.application.internal_server_error", nil)})
		slog.Error("internal error occurred", "request", c.Request.URL.Path, "error", err.Error())
		return true
	}

	status := http.StatusInternalServerError
	switch appErr.Kind {
	case cqrs.KindNotFound:
		status = http.StatusNotFound
	case cqrs.KindForbidden:
		status = http.StatusForbidden
	case cqrs.KindConflict:
		status = http.StatusConflict
	case cqrs.KindUnauthenticated:
		status = http.StatusUnauthorized
	case cqrs.KindInvalidInput:
		status = http.StatusBadRequest
	}

	c.JSON(status, gin.H{
		"error": tr.T(locale, appErr.Code, appErr.Params),
	})
	return true
}

func bindRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}
