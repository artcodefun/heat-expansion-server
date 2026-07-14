package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// actor extracts the authenticated admin from the gin context.
// Unauthenticated handlers (Register, Login) receive a zero-value Actor.
func actor(c *gin.Context) cqrs.Actor {
	if v, ok := c.Get("adminID"); ok {
		if id, ok2 := v.(uuid.UUID); ok2 {
			return cqrs.Actor{AdminID: id}
		}
	}
	return cqrs.Actor{AdminID: uuid.Nil}
}

func getLocale(c *gin.Context) string {
	lang := c.GetHeader("Accept-Language")
	if lang == "" {
		return "en"
	}
	return strings.Split(lang, ",")[0]
}

// handleCoreErr maps a domain or application error to an HTTP response and
// returns true so callers can do: if handleCoreErr(c, tr, err) { return }.
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

// bindRequest deserialises the JSON body and writes a 400 on failure.
// Returns false after writing the error response so callers can just return.
func bindRequest(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		slog.WarnContext(c.Request.Context(), "request rejected; invalid input",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"error", err.Error(),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}

// bindURI deserialises URI parameters and writes a 400 on failure.
func bindURI(c *gin.Context, req interface{}) bool {
	if err := c.ShouldBindUri(req); err != nil {
		slog.WarnContext(c.Request.Context(), "request rejected; invalid uri", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}
