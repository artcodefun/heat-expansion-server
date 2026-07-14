package dtos

import "github.com/google/uuid"

type UuidStr string

func (s UuidStr) Uuid() uuid.UUID {
	return uuid.MustParse(string(s))
}
