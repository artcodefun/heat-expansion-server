package security

import (
	"crypto/ecdsa"
	"errors"

	"github.com/artcodefun/heat-expansion-server/internal/billing/application/ports"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var _ ports.TokenValidator = (*SimpleTokenValidator)(nil)

// SimpleTokenValidator validates JWT tokens (ES256). Token expiry is enforced
// by the JWT `exp` claim set by the issuing service.
type SimpleTokenValidator struct {
	publicKey *ecdsa.PublicKey
}

// NewSimpleTokenValidator creates a validator for ES256 tokens signed by the auth service.
func NewSimpleTokenValidator(publicKey *ecdsa.PublicKey) *SimpleTokenValidator {
	return &SimpleTokenValidator{publicKey: publicKey}
}

// Validate verifies signature and expiry and returns the subject userID.
func (p *SimpleTokenValidator) Validate(tokenString string) (uuid.UUID, error) {
	if tokenString == "" {
		return uuid.Nil, errors.New("empty token")
	}
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return p.publicKey, nil
	}, jwt.WithValidMethods([]string{"ES256"}))
	if err != nil {
		return uuid.Nil, err
	}
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}
	uid, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.New("invalid subject")
	}
	return uid, nil
}
