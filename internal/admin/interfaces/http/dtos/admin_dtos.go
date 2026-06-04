package dtos

import (
	"github.com/artcodefun/heat-expansion-server/internal/admin/application/cqrs/readmodels"
	"github.com/google/uuid"
)

// RegisterRequest is the body for POST /api/v1/auth/register.
type RegisterRequest struct {
	Username    string `json:"username"     binding:"required"`
	InviteToken string `json:"invite_token" binding:"required"`
	Password    string `json:"password"     binding:"required,min=8"`
}

// LoginRequest is the body for POST /api/v1/auth/login.
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SessionResponse is returned after a successful Register or Login.
type SessionResponse struct {
	Token string `json:"token"`
}

// ProfileResponse is returned from GET /api/v1/auth/me.
type ProfileResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Active    bool      `json:"active"`
	CreatedAt int64     `json:"created_at"` // unix seconds
}

// ProfileResponseFromModel maps a readmodel to a ProfileResponse.
func ProfileResponseFromModel(p *readmodels.AdminProfile) ProfileResponse {
	return ProfileResponse{
		ID:        p.ID,
		Username:  p.Username,
		Active:    p.Active,
		CreatedAt: p.CreatedAt,
	}
}
