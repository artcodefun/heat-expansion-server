package ports

// PasswordHasher defines the interface for password hashing and verification.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}
