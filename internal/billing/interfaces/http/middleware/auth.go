package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/gin-gonic/gin"
)

// Auth attaches authenticated userID to the context if token is valid.
// Expects Authorization: Bearer <token> header.
func Auth(validator ports.TokenValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			slog.WarnContext(c.Request.Context(), "request rejected; missing bearer token", "method", c.Request.Method, "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimSpace(auth[len("Bearer "):])
		userID, err := validator.Validate(token)
		if err != nil {
			slog.WarnContext(c.Request.Context(), "request rejected; invalid bearer token", "method", c.Request.Method, "path", c.Request.URL.Path, "error", err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
