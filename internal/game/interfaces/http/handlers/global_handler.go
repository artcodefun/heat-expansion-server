package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type GlobalHandler struct{}

func NewGlobalHandler() *GlobalHandler {
	return &GlobalHandler{}
}

// GetMinClientVersion handles GET /min-client-version.
func (h *GlobalHandler) GetMinClientVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"version": "0.2.0"})
}
