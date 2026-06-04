package readmodels

import (
	"github.com/google/uuid"
)

// AdminProfile is returned from GetProfile queries.
type AdminProfile struct {
	ID        uuid.UUID
	Username  string
	Active    bool
	CreatedAt int64 // unix seconds
}
