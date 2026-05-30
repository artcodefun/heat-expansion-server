package domain

import (
	"github.com/google/uuid"
)

type OrderStatus string

const (
	OrderStatusPending OrderStatus = "PENDING"
	OrderStatusPaid    OrderStatus = "PAID"
	OrderStatusFailed  OrderStatus = "FAILED"
)

type PaymentProvider string

const PaymentProviderYooKassa PaymentProvider = "YOOKASSA"

type PurchaseOrder struct {
	EventProducer
	ID               uuid.UUID
	UserID           uuid.UUID
	PackageID        uuid.UUID
	Crystals         int
	AmountMinorUnits int64
	Currency         string
	Provider         PaymentProvider
	Status           OrderStatus
	ProviderOrderID  string
	ConfirmationURL  string
	CreatedAt        int64
	UpdatedAt        int64
}

func NewPendingOrder(userID, packageID uuid.UUID, crystals int, amountMinorUnits int64, currency string, provider PaymentProvider) *PurchaseOrder {
	id := uuid.Must(uuid.NewV7())
	now := NowUnix()
	return &PurchaseOrder{
		ID:               id,
		UserID:           userID,
		PackageID:        packageID,
		Crystals:         crystals,
		AmountMinorUnits: amountMinorUnits,
		Currency:         currency,
		Provider:         provider,
		Status:           OrderStatusPending,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func (o *PurchaseOrder) AttachProviderData(providerOrderID, confirmationURL string) {
	o.ProviderOrderID = providerOrderID
	o.ConfirmationURL = confirmationURL
	o.UpdatedAt = NowUnix()
}

func (o *PurchaseOrder) MarkPaid() error {
	if o.Status == OrderStatusPaid {
		return nil // already paid – idempotent
	}
	if o.Status != OrderStatusPending {
		return NewError("error.domain.order.invalid_status_transition", H{"status": string(o.Status)})
	}
	o.Status = OrderStatusPaid
	o.UpdatedAt = NowUnix()
	o.AddEvent(NewOrderPaidEvent(o.ID, o.UserID, o.PackageID, o.Crystals))
	return nil
}

func (o *PurchaseOrder) MarkFailed() error {
	if o.Status == OrderStatusFailed {
		return nil
	}
	if o.Status != OrderStatusPending {
		return NewError("error.domain.order.invalid_status_transition", H{"status": string(o.Status)})
	}
	o.Status = OrderStatusFailed
	o.UpdatedAt = NowUnix()
	o.AddEvent(NewOrderFailedEvent(o.ID))
	return nil
}
