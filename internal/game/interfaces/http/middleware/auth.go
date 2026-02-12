package middleware

import (
	"net/http"
	"strings"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/gin-gonic/gin"
)

// Auth attaches authenticated userID to the context if token is valid.
// Expects Authorization: Bearer <token> header.
func Auth(provider ports.TokenProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
			return
		}
		token := strings.TrimSpace(auth[len("Bearer "):])
		userID, err := provider.Validate(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}
