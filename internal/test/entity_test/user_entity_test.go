package entity_test

import (
	"testing"

	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestUser_Structure(t *testing.T) {
	user := entity.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}

	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "hashedpassword", user.Password)
	assert.NotNil(t, user.Profile)
}
