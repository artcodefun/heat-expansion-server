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
	Code    string                 // Machine readable code
	Message string                 // Developer-friendly message
	Kind    ErrorKind              // High-level category
	Params  map[string]interface{} // Dynamic context for i18n
}

func (e AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func NewAppError(kind ErrorKind, code, msg string) AppError {
	return AppError{Kind: kind, Code: code, Message: msg}
}

var (
	ErrNotFound           = NewAppError(KindNotFound, "common.not_found", "resource not found")
	ErrForbidden          = NewAppError(KindForbidden, "common.forbidden", "permission denied")
	ErrEmailAlreadyInUse  = NewAppError(KindConflict, "auth.email_taken", "email already in use")
	ErrInvalidCredentials = NewAppError(KindUnauthenticated, "auth.invalid_creds", "invalid email or password")
)
