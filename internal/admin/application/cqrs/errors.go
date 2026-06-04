package cqrs

import "fmt"

// ErrorKind classifies application-layer errors for HTTP status mapping.
type ErrorKind int

const (
	KindInternal     ErrorKind = iota
	KindNotFound               // 404
	KindForbidden              // 403
	KindConflict               // 409
	KindInvalidInput           // 422
)

// AppError is a translatable application-layer error.
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
	ErrNotFound  = NewAppError(KindNotFound, "error.application.not_found")
	ErrForbidden = NewAppError(KindForbidden, "error.application.forbidden")
)
