package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/auth/application/cqrs"
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

func handleCoreErr(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	var appErr cqrs.AppError
	if !errors.As(err, &appErr) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
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
		"code":    appErr.Code,
		"message": appErr.Message,
		"params":  appErr.Params,
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
