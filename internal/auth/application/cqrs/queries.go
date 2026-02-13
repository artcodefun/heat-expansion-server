package cqrs

import "github.com/google/uuid"

type QueryContext struct {
	AccountID uuid.UUID
}

type AccountQueries interface {
	// Add queries as needed
}
