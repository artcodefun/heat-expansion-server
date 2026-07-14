package ports

// SessionTokenGenerator generates opaque, cryptographically-random session tokens.
type SessionTokenGenerator interface {
	Generate() (string, error)
}
