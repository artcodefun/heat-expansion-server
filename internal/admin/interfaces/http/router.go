package http

import (
	"net/http"

	contract "github.com/artcodefun/heat-expansion-server/contracts/admin/http/v1"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/admin/interfaces/http/handlers"
	"github.com/artcodefun/heat-expansion-server/internal/admin/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Commands groups CQRS command interfaces needed by HTTP handlers.
type Commands struct {
	Admin cqrs.AdminCommands
}

// Queries groups CQRS query interfaces needed by HTTP handlers.
type Queries struct {
	Admin cqrs.AdminQueries
}

// NewRouter constructs the Gin engine, registers middleware and routes.
func NewRouter(cmd Commands, qry Queries, sessionValidator ports.SessionValidator, tr ports.Translator) *gin.Engine {
	r := gin.Default()
	r.Use(otelgin.Middleware("heat-expansion-admin"))

	handler := handlers.NewAdminHandler(cmd.Admin, qry.Admin, tr)

	// Global routes
	r.GET("/health", HealthHandler)

	// Public routes (no session required)
	publicApi := r.Group("/api/v1")
	{
		publicApi.GET("/openapi.yaml", func(c *gin.Context) {
			c.Data(http.StatusOK, "application/yaml; charset=utf-8", contract.OpenAPI())
		})
		publicApi.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("../openapi.yaml")))

		// Identity: establish a session.
		auth := publicApi.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}
	}

	// Private routes (valid session required)
	api := r.Group("/api/v1")
	api.Use(middleware.Auth(sessionValidator))
	{
		// Identity: manage the current session.
		auth := api.Group("/auth")
		{
			auth.POST("/logout", handler.Logout)
			auth.GET("/me", handler.Me)
		}
	}

	return r
}
