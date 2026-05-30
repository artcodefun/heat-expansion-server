package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
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
	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		return "en"
	}
	return strings.Split(lang, ",")[0]
}

func handleCoreErr(c *gin.Context, tr ports.Translator, err error) bool {
	if err == nil {
		return false
	}
	locale := getLocale(c)
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
		case cqrs.KindUnavailable:
			status = http.StatusServiceUnavailable
		}
		c.JSON(status, gin.H{"error": tr.T(locale, appErr.Code, appErr.Params)})
		return true
	}
	var domErr domain.Error
	if errors.As(err, &domErr) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": tr.T(locale, domErr.Key, domErr.Params)})
		return true
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": tr.T(locale, "error.application.internal_server_error", nil)})
	slog.ErrorContext(c.Request.Context(), "internal error occurred", "request", c.Request.URL.Path, "error", err.Error())
	return true
}

func bindRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		slog.WarnContext(c.Request.Context(), "request rejected; invalid input", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}

func bindURI(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindUri(req); err != nil {
		slog.WarnContext(c.Request.Context(), "request rejected; invalid uri", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}
