package security

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// SimpleTokenValidator validates JWT tokens (ES256). Token expiry is enforced
// by the JWT `exp` claim set by the issuing service.
type SimpleTokenValidator struct {
	publicKey *ecdsa.PublicKey
}

// NewSimpleTokenValidator parses a PEM-encoded ES256 public key and returns a validator.
func NewSimpleTokenValidator(publicKeyPEM string) (*SimpleTokenValidator, error) {
	publicKeyPEM = strings.ReplaceAll(publicKeyPEM, `\n`, "\n")
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode PEM block for EC public key")
	}
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ecKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("key is not an ECDSA public key")
	}
	if ecKey.Curve == nil || ecKey.Curve.Params().Name != "P-256" {
		return nil, errors.New("ECDSA public key must use P-256 curve for ES256")
	}
	return &SimpleTokenValidator{publicKey: ecKey}, nil
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
