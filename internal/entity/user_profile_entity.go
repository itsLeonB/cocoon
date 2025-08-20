package entity

import "github.com/google/uuid"

type UserProfile struct {
	BaseEntity
	UserID uuid.UUID
	Name   string
}

func (up UserProfile) IsAnonymous() bool {
	return up.UserID == uuid.Nil
}
