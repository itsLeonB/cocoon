package dto

import "github.com/google/uuid"

type RegisterRequest struct {
	Email                string `validate:"required,email,min=3"`
	Password             string `validate:"required,eqfield=PasswordConfirmation"`
	PasswordConfirmation string `validate:"required"`
	VerificationURL      string
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

func NewBearerTokenResp(token string) LoginResponse {
	return LoginResponse{
		Type:  "Bearer",
		Token: token,
	}
}

type OAuthCallbackData struct {
	Provider string `validate:"required,min=1"`
	Code     string `validate:"required,min=1"`
	State    string `validate:"required,min=1"`
}
