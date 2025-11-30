package cqrs

import "errors"

var (
	ErrForbidden = errors.New("forbidden")
	ErrNotFound  = errors.New("not found")
)

type ValidationError struct {
	Fields map[string]string
}

func (e ValidationError) Error() string { return "validation failed" }
