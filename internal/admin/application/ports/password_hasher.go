package ports

// PasswordHasher defines hashing and verification for passwords.
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}
