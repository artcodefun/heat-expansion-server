package security

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SimpleTokenProvider struct {
	privateKey *ecdsa.PrivateKey
}

func NewSimpleTokenProvider(privateKeyPEM string) (*SimpleTokenProvider, error) {
	privateKeyPEM = strings.ReplaceAll(privateKeyPEM, `\n`, "\n")
	block, _ := pem.Decode([]byte(privateKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode PEM block for EC private key")
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	if key.Curve == nil || key.Curve.Params().Name != "P-256" {
		return nil, errors.New("ECDSA private key must use P-256 curve for ES256")
	}
	return &SimpleTokenProvider{privateKey: key}, nil
}

func (s *SimpleTokenProvider) Generate(accountID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": accountID.String(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString(s.privateKey)
}

func (s *SimpleTokenProvider) Validate(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return &s.privateKey.PublicKey, nil
	}, jwt.WithValidMethods([]string{"ES256"}))

	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		sub, ok := claims["sub"].(string)
		if !ok {
			return uuid.Nil, fmt.Errorf("invalid subject claim")
		}
		uid, err := uuid.Parse(sub)
		if err != nil {
			return uuid.Nil, fmt.Errorf("invalid user id in token: %w", err)
		}
		return uid, nil
	}

	return uuid.Nil, fmt.Errorf("invalid token")
}
