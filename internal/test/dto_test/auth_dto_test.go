package dto_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestRegisterRequest_Structure(t *testing.T) {
	req := dto.RegisterRequest{
		Email:                "test@example.com",
		Password:             "password123",
		PasswordConfirmation: "password123",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
	assert.Equal(t, "password123", req.PasswordConfirmation)
}

func TestLoginRequest_Structure(t *testing.T) {
	req := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "password123", req.Password)
}

func TestLoginResponse_Structure(t *testing.T) {
	resp := dto.LoginResponse{
		Type:  "Bearer",
		Token: "jwt.token.here",
	}

	assert.Equal(t, "Bearer", resp.Type)
	assert.Equal(t, "jwt.token.here", resp.Token)
}

func TestAuthData_Structure(t *testing.T) {
	profileID := uuid.New()
	authData := dto.AuthData{
		ProfileID: profileID,
	}

	assert.Equal(t, profileID, authData.ProfileID)
}
