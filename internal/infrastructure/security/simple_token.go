package security

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SimpleTokenProvider issues and validates JWT tokens (HS256).
type SimpleTokenProvider struct {
	secret string
	ttl    time.Duration
}

// NewSimpleTokenProvider creates a provider with default TTL of 1h.
func NewSimpleTokenProvider(secret string) *SimpleTokenProvider {
	return &SimpleTokenProvider{secret: secret, ttl: time.Hour}
}

// Generate creates a signed JWT with subject=userID and standard time claims.
func (p *SimpleTokenProvider) Generate(userID int) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(userID),
		Issuer:    "heat-expansion",
		Audience:  []string{"user"},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(p.ttl)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.secret))
}

// Validate verifies signature and expiry and returns the subject userID.
func (p *SimpleTokenProvider) Validate(tokenString string) (int, error) {
	if tokenString == "" {
		return 0, errors.New("empty token")
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
		return 0, err
	}
	if !token.Valid {
		return 0, errors.New("invalid token")
	}
	// Subject must be userID
	uid, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}
	return uid, nil
}
