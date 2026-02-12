package ports

// TokenProvider defines generating and validating authentication tokens.
// Validate returns the embedded userID if the token is valid.
type TokenProvider interface {
	Generate(userID int) (string, error)
	Validate(token string) (int, error)
}
