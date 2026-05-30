package handlers

import (
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/interfaces/http/dtos"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	commands   cqrs.OrderCommands
	queries    cqrs.OrderQueries
	translator ports.Translator
}

func NewOrderHandler(commands cqrs.OrderCommands, queries cqrs.OrderQueries, translator ports.Translator) *OrderHandler {
	return &OrderHandler{commands: commands, queries: queries, translator: translator}
}

// CreateOrder handles POST /orders.
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dtos.CreateOrderRequest
	if !bindRequest(c, &req) {
		return
	}
	act := actor(c)
	orderID, confirmationURL, err := h.commands.CreateOrder(c.Request.Context(), act, req.PackageID, req.ReturnURL)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusCreated, dtos.CreateOrderResponse{
		OrderID:         orderID,
		ConfirmationURL: confirmationURL,
	})
}

// GetOrder handles GET /orders/:id.
func (h *OrderHandler) GetOrder(c *gin.Context) {
	var uri dtos.GetOrderURI
	if !bindURI(c, &uri) {
		return
	}
	act := actor(c)
	order, err := h.queries.GetOrder(c.Request.Context(), act, uri.ID)
	if handleCoreErr(c, h.translator, err) {
		return
	}
	c.JSON(http.StatusOK, dtos.OrderStatusResponseFromReadModel(*order))
}
