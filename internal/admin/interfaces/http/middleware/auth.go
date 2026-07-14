package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/gin-gonic/gin"
)

const AdminIDKey = "adminID"

// Auth validates the bearer session token and attaches the admin ID to the context.
// Expects: Authorization: Bearer <token>
func Auth(validator ports.SessionValidator) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			slog.WarnContext(c.Request.Context(), "admin request rejected; missing bearer token",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimSpace(auth[len("Bearer "):])
		adminID, err := validator.ValidateSession(c.Request.Context(), token)
		if err != nil {
			slog.WarnContext(c.Request.Context(), "admin request rejected; invalid session token",
				"method", c.Request.Method,
				"path", c.Request.URL.Path,
				"error", err.Error(),
			)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired session"})
			return
		}
		c.Set(AdminIDKey, adminID)
		c.Next()
	}
}
