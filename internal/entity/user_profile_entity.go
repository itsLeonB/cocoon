package entity

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
)

type UserProfile struct {
	crud.BaseEntity
	UserID uuid.NullUUID
	Name   string
	Avatar string
}

type ProfileName struct {
	ID   uuid.UUID
	Name string
}
