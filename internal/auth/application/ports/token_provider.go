package ports

import "github.com/google/uuid"

// TokenProvider defines generating and validating authentication tokens.
// Validate returns the embedded accountID if the token is valid.
type TokenProvider interface {
	Generate(accountID uuid.UUID) (string, error)
	Validate(token string) (uuid.UUID, error)
}
