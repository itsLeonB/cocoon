package mapper_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mapper"
	"github.com/stretchr/testify/assert"
)

func TestUserToAuthData(t *testing.T) {
	userID := uuid.New()
	user := entity.User{
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	user.ID = userID

	result := mapper.UserToAuthData(user)

	assert.NotNil(t, result)
	assert.Equal(t, userID, result[appconstant.ContextUserID])
	assert.Len(t, result, 1)
}
