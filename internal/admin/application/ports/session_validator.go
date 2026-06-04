package ports

import (
	"context"

	"github.com/google/uuid"
)

// SessionValidator validates a session token and returns the authenticated admin ID.
type SessionValidator interface {
	ValidateSession(ctx context.Context, token string) (uuid.UUID, error)
}
