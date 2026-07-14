package security

import (
	"crypto/rand"
	"encoding/base64"

	"github.com/artcodefun/heat-expansion-server/internal/admin/application/ports"
)

// RandomSessionTokenGenerator generates a cryptographically-random, URL-safe
// base64-encoded session token (32 random bytes → 44-char string).
type RandomSessionTokenGenerator struct{}

func NewRandomSessionTokenGenerator() *RandomSessionTokenGenerator {
	return &RandomSessionTokenGenerator{}
}

func (g *RandomSessionTokenGenerator) Generate() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

var _ ports.SessionTokenGenerator = (*RandomSessionTokenGenerator)(nil)
