package cqrs

import (
	"context"

	authv1 "github.com/artcodefun/heat-expansion-server/contracts/auth/events/v1"
	"github.com/google/uuid"
)

type OrderCommands interface {
	CreateOrder(ctx context.Context, actor Actor, packageID uuid.UUID, returnURL string) (orderID uuid.UUID, confirmationURL string, err error)
	ConfirmPayment(ctx context.Context, rawBody []byte) error
}

type UserCommands interface {
	HandleAccountRegisteredV1Event(ctx context.Context, ev authv1.AccountRegisteredV1) error
}
