package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProfileResponse struct {
	ID                       uuid.UUID
	UserID                   uuid.UUID
	Name                     string
	Avatar                   string
	Email                    string
	CreatedAt                time.Time
	UpdatedAt                time.Time
	DeletedAt                time.Time
	AssociatedAnonProfileIDs []uuid.UUID
	RealProfileID            uuid.UUID
}

type NewProfileRequest struct {
	UserID uuid.UUID
	Name   string `validate:"required,min=1,max=255"`
	Avatar string
}

type UpdateProfileRequest struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Name   string
	Avatar string
}

type AssociateProfileRequest struct {
	UserProfileID uuid.UUID
	RealProfileID uuid.UUID
	AnonProfileID uuid.UUID
}
