package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
)

const yookassaBaseURL = "https://api.yookassa.ru/v3"

type YooKassaGateway struct {
	shopID    string
	secretKey string
	client    *http.Client
}

func NewYooKassaGateway(shopID, secretKey string) *YooKassaGateway {
	return &YooKassaGateway{
		shopID:    shopID,
		secretKey: secretKey,
		client:    &http.Client{Timeout: 30 * time.Second},
	}
}

type yookassaAmount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type yookassaConfirmation struct {
	Type      string `json:"type"`
	ReturnURL string `json:"return_url"`
}

type yookassaCreatePaymentRequest struct {
	Amount       yookassaAmount       `json:"amount"`
	Confirmation yookassaConfirmation `json:"confirmation"`
	Capture      bool                 `json:"capture"`
	Description  string               `json:"description"`
	Metadata     map[string]string    `json:"metadata"`
}

type yookassaConfirmationResponse struct {
	Type            string `json:"type"`
	ConfirmationURL string `json:"confirmation_url"`
}

type yookassaPaymentResponse struct {
	ID           string                       `json:"id"`
	Status       string                       `json:"status"`
	Confirmation yookassaConfirmationResponse `json:"confirmation"`
}

type yookassaWebhookObject struct {
	ID string `json:"id"`
}

type yookassaWebhookNotification struct {
	Type   string                `json:"type"`
	Event  string                `json:"event"`
	Object yookassaWebhookObject `json:"object"`
}

func (g *YooKassaGateway) CreatePayment(ctx context.Context, order *domain.PurchaseOrder, pkg *domain.CrystalPackage, returnURL string) (string, string, error) {
	// Format amount as decimal string (e.g. kopecks -> rubles: 9900 -> "99.00")
	rubles := order.AmountMinorUnits / 100
	kopecks := order.AmountMinorUnits % 100
	amountValue := fmt.Sprintf("%d.%02d", rubles, kopecks)

	reqBody := yookassaCreatePaymentRequest{
		Amount: yookassaAmount{
			Value:    amountValue,
			Currency: order.Currency,
		},
		Confirmation: yookassaConfirmation{
			Type:      "redirect",
			ReturnURL: returnURL,
		},
		Capture:     true,
		Description: fmt.Sprintf("Purchase of %d crystals (%s)", pkg.Crystals, pkg.Name),
		Metadata: map[string]string{
			"order_id":   order.ID.String(),
			"package_id": pkg.ID.String(),
			"user_id":    order.UserID.String(),
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", "", fmt.Errorf("yookassa: marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, yookassaBaseURL+"/payments", bytes.NewReader(body))
	if err != nil {
		return "", "", fmt.Errorf("yookassa: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotence-Key", order.ID.String())
	req.SetBasicAuth(g.shopID, g.secretKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("yookassa: http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("yookassa: unexpected status %d", resp.StatusCode)
	}

	var result yookassaPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", fmt.Errorf("yookassa: decode response: %w", err)
	}

	return result.ID, result.Confirmation.ConfirmationURL, nil
}

// VerifyWebhook parses the incoming notification to extract the payment ID,
// then re-queries YooKassa for the canonical payment state. This avoids
// trusting the webhook body and removes reliance on proxy-level IP filtering.
func (g *YooKassaGateway) VerifyWebhook(ctx context.Context, rawBody []byte) (string, bool, error) {
	var notification yookassaWebhookNotification
	if err := json.Unmarshal(rawBody, &notification); err != nil {
		return "", false, fmt.Errorf("yookassa webhook: unmarshal: %w", errors.Join(ports.ErrMalformedWebhook, err))
	}
	if notification.Event == "" || notification.Object.ID == "" {
		return "", false, fmt.Errorf("yookassa webhook: %w: missing required fields", ports.ErrMalformedWebhook)
	}

	payment, err := g.getPayment(ctx, notification.Object.ID)
	if err != nil {
		return "", false, fmt.Errorf("yookassa webhook: re-query failed: %w", err)
	}

	return payment.ID, payment.Status == "succeeded", nil
}

func (g *YooKassaGateway) getPayment(ctx context.Context, paymentID string) (*yookassaPaymentResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, yookassaBaseURL+"/payments/"+paymentID, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.SetBasicAuth(g.shopID, g.secretKey)

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var payment yookassaPaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&payment); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &payment, nil
}
