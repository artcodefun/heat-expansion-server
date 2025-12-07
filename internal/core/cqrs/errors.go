package cqrs

import "errors"

var (
	ErrForbidden = errors.New("forbidden")
	ErrNotFound  = errors.New("not found")
)

// DomainError wraps a domain-level error so adapters can distinguish it
// from infrastructure or unexpected failures.
type DomainError struct {
	Err error
}

func (e DomainError) Error() string {
	if e.Err == nil {
		return "domain error"
	}
	return e.Err.Error()
}

func (e DomainError) Unwrap() error { return e.Err }

// NewDomainError wraps the given error in a DomainError. If err is nil, it returns nil.
func NewDomainError(err error) error {
	if err == nil {
		return nil
	}
	return DomainError{Err: err}
}
