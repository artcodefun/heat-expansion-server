package domain

import (
	"time"

	"github.com/google/uuid"
)

const resetTokenTTL = time.Hour

type PasswordResetToken struct {
	ID        uuid.UUID
	AccountID uuid.UUID
	TokenHash string
	ExpiresAt int64
	UsedAt    *int64
}

func NewPasswordResetToken(accountID uuid.UUID, tokenHash string) *PasswordResetToken {
	id, _ := uuid.NewV7()
	return &PasswordResetToken{
		ID:        id,
		AccountID: accountID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(resetTokenTTL).Unix(),
	}
}

func (t *PasswordResetToken) IsExpired() bool {
	return time.Now().Unix() > t.ExpiresAt
}

func (t *PasswordResetToken) IsUsed() bool {
	return t.UsedAt != nil
}
