package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Email                string `validate:"required,email,min=3"`
	Password             string `validate:"required,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `validate:"required"`
}

type LoginRequest struct {
	Email    string `validate:"required,email,min=3"`
	Password string `validate:"required"`
}

type LoginResponse struct {
	Type  string
	Token string
}

type AuthData struct {
	ProfileID uuid.UUID
}
