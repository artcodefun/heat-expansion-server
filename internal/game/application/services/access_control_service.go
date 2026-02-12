package services

import (
	"errors"

	"github.com/artcodefun/heat-expansion-api/internal/game/application/cqrs"
	"github.com/artcodefun/heat-expansion-api/internal/game/application/ports"
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
// userID <= 0 -> ErrForbidden (unauthenticated)
// base not found -> ErrNotFound
// mismatch owner -> ErrForbidden
func (s *AccessControlService) EnsureBaseOwnership(userID int, baseID int) error {
	if userID <= 0 {
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
