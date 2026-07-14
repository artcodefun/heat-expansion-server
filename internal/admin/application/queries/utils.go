package queries

import (
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
)

// repoErr maps port-level repository sentinel errors to application CQRS errors.
func repoErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ports.ErrNotFound) {
		return cqrs.ErrNotFound
	}
	return err
}

// clientErr maps port-level sentinel errors from downstream service clients
// to application CQRS errors. Infrastructure adapters (e.g. gRPC clients)
// return sentinels so the application layer owns the mapping.
func clientErr(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, ports.ErrClientNotFound):
		return cqrs.ErrNotFound
	case errors.Is(err, ports.ErrClientInvalidInput):
		return cqrs.ErrInvalidInput
	case errors.Is(err, ports.ErrClientForbidden):
		return cqrs.ErrForbidden
	default:
		return err
	}
}
