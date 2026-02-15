package cqrs

import (
	"fmt"
)

type ErrorKind int

const (
	KindInternal ErrorKind = iota
	KindNotFound
	KindForbidden
	KindUnauthenticated
	KindConflict
	KindInvalidInput
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
	ErrNotFound           = NewAppError(KindNotFound, "error.application.not_found")
	ErrForbidden          = NewAppError(KindForbidden, "error.application.forbidden")
	ErrEmailAlreadyInUse  = NewAppError(KindConflict, "error.application.auth.email_taken")
	ErrInvalidCredentials = NewAppError(KindUnauthenticated, "error.application.auth.invalid_creds")
)
