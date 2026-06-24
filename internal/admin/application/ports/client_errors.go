package ports

import "errors"

// Sentinel errors returned by downstream service clients (game, billing gRPC).
var (
	ErrClientNotFound     = errors.New("downstream service: not found")
	ErrClientInvalidInput = errors.New("downstream service: invalid input")
	ErrClientForbidden    = errors.New("downstream service: forbidden")
)
