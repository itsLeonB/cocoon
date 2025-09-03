package entity

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
)

type UserProfile struct {
	crud.BaseEntity
	UserID uuid.UUID
	Name   string
}

func (up UserProfile) IsAnonymous() bool {
	return up.UserID == uuid.Nil
}
