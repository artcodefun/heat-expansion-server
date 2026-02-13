package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type SimpleTokenProvider struct {
	secret string
}

func NewSimpleTokenProvider(secret string) *SimpleTokenProvider {
	return &SimpleTokenProvider{secret: secret}
}

func (s *SimpleTokenProvider) Generate(accountID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": accountID.String(),
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(s.secret))
}

func (s *SimpleTokenProvider) Validate(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secret), nil
	})

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
