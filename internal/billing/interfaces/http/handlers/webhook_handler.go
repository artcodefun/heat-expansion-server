package handlers

import (
	"io"
	"net/http"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	commands   cqrs.OrderCommands
	translator ports.Translator
}

func NewWebhookHandler(commands cqrs.OrderCommands, translator ports.Translator) *WebhookHandler {
	return &WebhookHandler{commands: commands, translator: translator}
}

// maxWebhookBodyBytes caps the webhook payload size to avoid unbounded reads
// on a public, unauthenticated endpoint.
const maxWebhookBodyBytes = 64 << 10 // 64 KiB

// HandleYooKassa handles POST /webhooks/yookassa.
func (h *WebhookHandler) HandleYooKassa(c *gin.Context) {
	body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, maxWebhookBodyBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot read body"})
		return
	}
	if err := h.commands.ConfirmPayment(c.Request.Context(), body); handleCoreErr(c, h.translator, err) {
		return
	}
	c.Status(http.StatusOK)
}
