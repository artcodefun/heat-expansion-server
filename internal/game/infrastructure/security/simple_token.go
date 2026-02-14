package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// SimpleTokenValidator validates JWT tokens (HS256).
type SimpleTokenValidator struct {
	secret string
	ttl    time.Duration
}

// NewSimpleTokenValidator creates a validator with default TTL of 1h.
func NewSimpleTokenValidator(secret string) *SimpleTokenValidator {
	return &SimpleTokenValidator{secret: secret, ttl: time.Hour}
}

// Validate verifies signature and expiry and returns the subject userID.
func (p *SimpleTokenValidator) Validate(tokenString string) (uuid.UUID, error) {
	if tokenString == "" {
		return uuid.Nil, errors.New("empty token")
	}
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		// Enforce HS256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(p.secret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}
	// Subject must be userID
	uid, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("invalid subject")
	}
	return uid, nil
}
