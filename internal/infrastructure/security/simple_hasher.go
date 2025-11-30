package security

import "golang.org/x/crypto/bcrypt"

// SimpleHasher is a dev-friendly PasswordHasher using bcrypt.
type SimpleHasher struct{}

func NewSimpleHasher() *SimpleHasher { return &SimpleHasher{} }

func (h *SimpleHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (h *SimpleHasher) Verify(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
