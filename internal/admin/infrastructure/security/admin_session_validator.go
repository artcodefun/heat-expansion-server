package security

import (
	"context"
	"errors"
	"fmt"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
	"github.com/google/uuid"
)

// AdminSessionValidator implements ports.SessionValidator. A token is valid only when its
// session exists and is unexpired AND the owning admin is still active, so that
// deactivating an admin immediately revokes their in-flight sessions.
type AdminSessionValidator struct {
	sessions ports.SessionRepository
	admins   ports.AdminRepository
}

func NewAdminSessionValidator(sessions ports.SessionRepository, admins ports.AdminRepository) *AdminSessionValidator {
	return &AdminSessionValidator{sessions: sessions, admins: admins}
}

func (v *AdminSessionValidator) ValidateSession(ctx context.Context, token string) (uuid.UUID, error) {
	session, err := v.sessions.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return uuid.Nil, fmt.Errorf("session not found or expired")
		}
		return uuid.Nil, fmt.Errorf("validate session: %w", err)
	}

	admin, err := v.admins.FindByID(ctx, session.AdminID)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return uuid.Nil, fmt.Errorf("admin not found")
		}
		return uuid.Nil, fmt.Errorf("validate session: %w", err)
	}
	if !admin.Active {
		return uuid.Nil, fmt.Errorf("admin is not active")
	}

	return session.AdminID, nil
}

var _ ports.SessionValidator = (*AdminSessionValidator)(nil)
