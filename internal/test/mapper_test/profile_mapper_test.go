package mapper_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mapper"
	"github.com/stretchr/testify/assert"
)

func TestProfileToResponse(t *testing.T) {
	userID := uuid.New()
	profileID := uuid.New()
	now := time.Now()

	profile := entity.UserProfile{
		UserID: userID,
		Name:   "John Doe",
		Avatar: "avatar.jpeg",
	}
	profile.ID = profileID
	profile.CreatedAt = now
	profile.UpdatedAt = now

	response := mapper.ProfileToResponse(profile)

	assert.Equal(t, profileID, response.ID)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, "John Doe", response.Name)
	assert.Equal(t, now, response.CreatedAt)
	assert.Equal(t, now, response.UpdatedAt)
	assert.Equal(t, "avatar.jpeg", response.Avatar)
}

func TestProfileToResponse_WithDeletedAt(t *testing.T) {
	userID := uuid.New()
	now := time.Now()

	profile := entity.UserProfile{
		UserID: userID,
		Name:   "Deleted User",
	}
	profile.DeletedAt.Time = now
	profile.DeletedAt.Valid = true

	response := mapper.ProfileToResponse(profile)

	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, "Deleted User", response.Name)
	assert.Equal(t, now, response.DeletedAt)
}
