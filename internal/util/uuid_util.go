package util

import "github.com/google/uuid"

func NewValidNullUUID(u uuid.UUID) uuid.NullUUID {
	return uuid.NullUUID{
		UUID:  u,
		Valid: true,
	}
}
