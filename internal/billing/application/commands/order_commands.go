package commands

import (
	"context"
	"errors"
	"log/slog"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/artcodefun/heat-expansion-server/internal/billing/domain"
	"github.com/google/uuid"
)

type OrderCommands struct {
	OrderRepo   ports.PurchaseOrderRepository
	PackageRepo ports.CrystalPackageRepository
	UserRepo    ports.UserRepository
	Gateway     ports.PaymentGateway
	Outbox      ports.OutboxEventRepository
	TxMgr       ports.TransactionManager
}

func NewOrderCommands(
	orderRepo ports.PurchaseOrderRepository,
	packageRepo ports.CrystalPackageRepository,
	userRepo ports.UserRepository,
	gateway ports.PaymentGateway,
	outbox ports.OutboxEventRepository,
	txMgr ports.TransactionManager,
) *OrderCommands {
	return &OrderCommands{
		OrderRepo:   orderRepo,
		PackageRepo: packageRepo,
		UserRepo:    userRepo,
		Gateway:     gateway,
		Outbox:      outbox,
		TxMgr:       txMgr,
	}
}

// paymentsEnabled gates order creation. It is temporarily false while the
// YooKassa account is pending moderation; flip it to true (and remove
// cqrs.ErrPaymentsUnavailable) once payments go live.
var paymentsEnabled = false

func (c *OrderCommands) CreateOrder(ctx context.Context, actor cqrs.Actor, packageID uuid.UUID, returnURL string) (uuid.UUID, string, error) {
	if !paymentsEnabled {
		// Reject up front: no order is persisted and the gateway is never
		// called. Package listing is unaffected.
		return uuid.Nil, "", cqrs.ErrPaymentsUnavailable
	}

	pkg, err := c.PackageRepo.FindByID(ctx, packageID)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return uuid.Nil, "", cqrs.ErrPackageNotFound
		}
		return uuid.Nil, "", err
	}
	if !pkg.IsActive {
		return uuid.Nil, "", cqrs.ErrPackageNotFound
	}

	// The fiscal receipt (54-FZ) must be issued to the buyer's email, which we
	// project locally from auth registration events. Resolve it before
	// persisting anything so we fail fast if the projection has not caught up.
	user, err := c.UserRepo.FindByID(ctx, actor.UserID)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			slog.WarnContext(ctx, "cannot create order: buyer email not yet projected", "user_id", actor.UserID.String())
			return uuid.Nil, "", cqrs.ErrCustomerEmailUnavailable
		}
		return uuid.Nil, "", err
	}
	if user.Email == "" {
		slog.WarnContext(ctx, "cannot create order: buyer email is empty", "user_id", actor.UserID.String())
		return uuid.Nil, "", cqrs.ErrCustomerEmailUnavailable
	}

	order := domain.NewPendingOrder(
		actor.UserID,
		pkg.ID,
		pkg.Crystals,
		pkg.PriceMinorUnits,
		pkg.Currency,
		domain.PaymentProviderYooKassa,
	)

	// Persist the pending order before calling the gateway so it durably exists
	// with the ID we pass to YooKassa as the idempotence key and metadata
	// order_id, even if the process crashes during or right after CreatePayment.
	// Webhook matching is keyed on provider_order_id, which is only stored by the
	// Update below; a webhook racing this method therefore relies on YooKassa's
	// delivery retries rather than on this initial save. The two writes need not
	// be atomic, so no transaction is used.
	if err := c.OrderRepo.Save(ctx, order); err != nil {
		return uuid.Nil, "", err
	}

	providerOrderID, confirmationURL, err := c.Gateway.CreatePayment(ctx, order, pkg, user.Email, returnURL)
	if err != nil {
		slog.ErrorContext(ctx, "payment gateway failed to create payment", "error", err)
		return uuid.Nil, "", cqrs.ErrPaymentGatewayFailed
	}

	order.AttachProviderData(providerOrderID, confirmationURL)

	if err := c.OrderRepo.Update(ctx, order); err != nil {
		return uuid.Nil, "", err
	}

	return order.ID, confirmationURL, nil
}

func (c *OrderCommands) ConfirmPayment(ctx context.Context, rawBody []byte) error {
	providerOrderID, paid, err := c.Gateway.VerifyWebhook(ctx, rawBody)
	if err != nil {
		if errors.Is(err, ports.ErrMalformedWebhook) {
			return cqrs.ErrInvalidWebhookPayload
		}
		// Transient failure (e.g. re-query to the provider failed): surface as
		// internal so the provider retries delivery rather than giving up.
		slog.ErrorContext(ctx, "failed to verify payment webhook", "error", err)
		return cqrs.ErrPaymentGatewayFailed
	}

	return c.TxMgr.WithTx(ctx, func(tx ports.Transaction) error {
		oRepo := c.OrderRepo.Tx(tx)

		order, err := oRepo.FindByProviderOrderIDForUpdate(ctx, providerOrderID)
		if err != nil {
			if errors.Is(err, ports.ErrNotFound) {
				return nil // Unknown order – ignore
			}
			return err
		}

		if paid {
			if err := order.MarkPaid(); err != nil {
				return err
			}
		} else {
			if err := order.MarkFailed(); err != nil {
				return err
			}
		}

		if err := oRepo.Update(ctx, order); err != nil {
			return err
		}

		return c.Outbox.Tx(tx).Save(ctx, order.EventProducer.PullEvents())
	})
}
