package cqrs

import (
	"fmt"
)

type ErrorKind int

const (
	KindInternal ErrorKind = iota
	KindNotFound
	KindForbidden
	KindConflict
	KindInvalidInput
	KindUnavailable
)

type AppError struct {
	Kind   ErrorKind
	Code   string
	Params map[string]any
}

func (e AppError) Error() string {
	return fmt.Sprintf("[%d] %s", e.Kind, e.Code)
}

func NewAppError(kind ErrorKind, code string) AppError {
	return AppError{Kind: kind, Code: code}
}

func NewAppErrorWithParams(kind ErrorKind, code string, params map[string]any) AppError {
	return AppError{Kind: kind, Code: code, Params: params}
}

var (
	ErrNotFound                = NewAppError(KindNotFound, "error.application.not_found")
	ErrForbidden               = NewAppError(KindForbidden, "error.application.forbidden")
	ErrOrderAlreadyPaid        = NewAppError(KindConflict, "error.application.billing.order_already_paid")
	ErrPackageNotFound         = NewAppError(KindNotFound, "error.application.billing.package_not_found")
	ErrOrderNotFound           = NewAppError(KindNotFound, "error.application.billing.order_not_found")
	ErrPaymentGatewayFailed    = NewAppError(KindInternal, "error.application.billing.payment_gateway_failed")
	ErrInvalidWebhookPayload   = NewAppError(KindInvalidInput, "error.application.billing.invalid_webhook_payload")
	// ErrPaymentsUnavailable is returned while the payment provider account is
	// pending moderation. Remove it (and its call site in CreateOrder) once
	// payments are live.
	ErrPaymentsUnavailable     = NewAppError(KindUnavailable, "error.application.billing.payments_unavailable")
)
