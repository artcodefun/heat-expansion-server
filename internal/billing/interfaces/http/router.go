package http

import (
	"net/http"

	contract "github.com/artcodefun/heat-expansion-server/contracts/billing/http/v1"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/interfaces/http/handlers"
	"github.com/artcodefun/heat-expansion-server/internal/billing/interfaces/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// Commands groups CQRS command interfaces needed by HTTP handlers.
type Commands struct {
	Order cqrs.OrderCommands
}

// Queries groups CQRS query interfaces needed by HTTP handlers.
type Queries struct {
	Package cqrs.PackageQueries
	Order   cqrs.OrderQueries
}

// NewRouter constructs the Gin engine, registers middleware and routes.
func NewRouter(cmd Commands, qry Queries, tokenValidator ports.TokenValidator, tr ports.Translator) *gin.Engine {
	r := gin.Default()
	r.Use(otelgin.Middleware("heat-expansion-billing"))

	// Initialize handlers at the top for consistency
	packageHandler := handlers.NewPackageHandler(qry.Package, tr)
	orderHandler := handlers.NewOrderHandler(cmd.Order, qry.Order, tr)
	webhookHandler := handlers.NewWebhookHandler(cmd.Order, tr)

	// Global routes
	r.GET("/health", HealthHandler)

	// Public routes (no auth)
	publicApi := r.Group("/api/v1")
	{
		publicApi.GET("/openapi.yaml", func(c *gin.Context) {
			c.Data(http.StatusOK, "application/yaml; charset=utf-8", contract.OpenAPI())
		})
		publicApi.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("../openapi.yaml")))

		// YooKassa payment webhook (verified by re-querying the provider)
		publicApi.POST("/webhooks/yookassa", webhookHandler.HandleYooKassa)
	}

	// Private routes (auth required)
	api := r.Group("/api/v1")
	api.Use(middleware.Auth(tokenValidator))
	{
		// Packages
		api.GET("/packages", packageHandler.ListPackages)

		// Orders
		orders := api.Group("/orders")
		{
			orders.POST("", orderHandler.CreateOrder)
			orders.GET("/:id", orderHandler.GetOrder)
		}
	}

	return r
}
