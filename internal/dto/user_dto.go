package dto

import (
	"time"

	"github.com/google/uuid"
)

type NewUserRequest struct {
	Email     string
	Password  string
	Name      string
	Avatar    string
	VerifyNow bool
}

type UserResponse struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Profile   ProfileResponse
}
