package cqrs

import (
	"context"

	"github.com/google/uuid"
)

type OrderCommands interface {
	CreateOrder(ctx context.Context, actor Actor, packageID uuid.UUID, returnURL string) (orderID uuid.UUID, confirmationURL string, err error)
	ConfirmPayment(ctx context.Context, rawBody []byte) error
}
