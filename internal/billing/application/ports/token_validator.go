package ports

import "github.com/google/uuid"

type TokenValidator interface {
	Validate(token string) (uuid.UUID, error)
}
