package services

import (
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-server/internal/game/application/ports"
	"github.com/google/uuid"
)

// AccessControlService centralizes authorization checks for aggregates.
// Keeps authorization logic outside domain & closer to application layer.
type AccessControlService struct {
	Bases ports.UserBaseRepository
}

func NewAccessControlService(bases ports.UserBaseRepository) *AccessControlService {
	return &AccessControlService{Bases: bases}
}

// EnsureBaseOwnership verifies that the provided userID owns the base.
// userID == uuid.Nil -> ErrForbidden (unauthenticated)
// base not found -> ErrNotFound
// mismatch owner -> ErrForbidden
func (s *AccessControlService) EnsureBaseOwnership(userID uuid.UUID, baseID int) error {
	if userID == uuid.Nil {
		return cqrs.ErrForbidden
	}
	ownerID, err := s.Bases.GetOwnerID(baseID)
	if err != nil {
		if errors.Is(err, ports.ErrNotFound) {
			return cqrs.ErrNotFound
		}
		return err
	}
	if ownerID != userID {
		return cqrs.ErrForbidden
	}
	return nil
}
