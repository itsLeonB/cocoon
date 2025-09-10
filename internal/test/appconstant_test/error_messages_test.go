package appconstant_test

import (
	"fmt"
	"testing"

	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/stretchr/testify/assert"
)

func TestErrorMessages_Constants(t *testing.T) {
	assert.Equal(t, "error retrieving data", appconstant.ErrDataSelect)
	assert.Equal(t, "error inserting new data", appconstant.ErrDataInsert)
	assert.Equal(t, "user is not found", appconstant.ErrAuthUserNotFound)
	assert.Equal(t, "user with email %s is already registered", appconstant.ErrAuthDuplicateUser)
	assert.Equal(t, "unknown credentials, please check your email/password", appconstant.ErrAuthUnknownCredentials)
	assert.Equal(t, "user with ID: %s is not found", appconstant.ErrUserNotFound)
	assert.Equal(t, "user with ID: %s is deleted", appconstant.ErrUserDeleted)
	assert.Equal(t, "friendship not found", appconstant.ErrFriendshipNotFound)
	assert.Equal(t, "friendship is deleted", appconstant.ErrFriendshipDeleted)
}

func TestErrorMessages_Formatting(t *testing.T) {
	email := "test@example.com"
	userID := "123e4567-e89b-12d3-a456-426614174000"

	formattedDuplicateUser := fmt.Sprintf(appconstant.ErrAuthDuplicateUser, email)
	formattedUserNotFound := fmt.Sprintf(appconstant.ErrUserNotFound, userID)
	formattedUserDeleted := fmt.Sprintf(appconstant.ErrUserDeleted, userID)

	assert.Equal(t, "user with email test@example.com is already registered", formattedDuplicateUser)
	assert.Equal(t, "user with ID: 123e4567-e89b-12d3-a456-426614174000 is not found", formattedUserNotFound)
	assert.Equal(t, "user with ID: 123e4567-e89b-12d3-a456-426614174000 is deleted", formattedUserDeleted)
}
