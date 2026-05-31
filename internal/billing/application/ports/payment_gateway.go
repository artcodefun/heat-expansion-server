package ports

import (
	"context"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
)

// ErrMalformedWebhook indicates the webhook body could not be parsed or was
// missing required fields. It is a permanent, client-side failure and must be
// distinguished from transient errors (e.g. a failed re-query of the provider),
// which should be retried.
var ErrMalformedWebhook = errors.New("malformed webhook payload")

// PaymentGateway abstracts a payment provider (YooKassa, etc.)
type PaymentGateway interface {
	// CreatePayment creates a payment and returns the provider order ID and redirect URL.
	// customerEmail is the address the fiscal receipt (54-FZ) is issued to.
	CreatePayment(ctx context.Context, order *domain.PurchaseOrder, pkg *domain.CrystalPackage, customerEmail, returnURL string) (providerOrderID, confirmationURL string, err error)
	// VerifyWebhook validates an incoming webhook notification by re-querying
	// the provider for the canonical payment state.
	// Returns the providerOrderID and whether the payment succeeded.
	VerifyWebhook(ctx context.Context, rawBody []byte) (providerOrderID string, paid bool, err error)
}
