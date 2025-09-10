package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestUserProfile_IsAnonymous(t *testing.T) {
	tests := []struct {
		name     string
		userID   uuid.UUID
		expected bool
	}{
		{
			name:     "anonymous profile with nil UUID",
			userID:   uuid.Nil,
			expected: true,
		},
		{
			name:     "non-anonymous profile with valid UUID",
			userID:   uuid.New(),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := entity.UserProfile{
				UserID: tt.userID,
				Name:   "Test User",
			}

			result := profile.IsAnonymous()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestUserProfile_Structure(t *testing.T) {
	userID := uuid.New()
	profile := entity.UserProfile{
		UserID: userID,
		Name:   "John Doe",
	}

	assert.Equal(t, userID, profile.UserID)
	assert.Equal(t, "John Doe", profile.Name)
}
