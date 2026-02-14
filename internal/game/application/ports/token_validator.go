package ports

import "github.com/google/uuid"

// TokenValidator defines validation of authentication tokens.
// Validate returns the embedded userID if the token is valid.
type TokenValidator interface {
	Validate(token string) (uuid.UUID, error)
}
