package dtos

import "github.com/google/uuid"

type UuidStr string

func (str UuidStr) Uuid() uuid.UUID {
	return uuid.MustParse(string(str))
}

type Request[U, Q, B any] struct {
	Uri   U
	Query Q
	Body  B
}

type None struct{}
