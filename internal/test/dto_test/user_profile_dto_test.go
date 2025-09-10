package dto_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestProfileResponse_Structure(t *testing.T) {
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	resp := dto.ProfileResponse{
		ID:          id,
		UserID:      userID,
		Name:        "John Doe",
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   time.Time{},
		IsAnonymous: false,
	}

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, userID, resp.UserID)
	assert.Equal(t, "John Doe", resp.Name)
	assert.Equal(t, now, resp.CreatedAt)
	assert.Equal(t, now, resp.UpdatedAt)
	assert.Equal(t, time.Time{}, resp.DeletedAt)
	assert.False(t, resp.IsAnonymous)
}

func TestNewProfileRequest_Structure(t *testing.T) {
	userID := uuid.New()
	req := dto.NewProfileRequest{
		UserID: userID,
		Name:   "Jane Doe",
	}

	assert.Equal(t, userID, req.UserID)
	assert.Equal(t, "Jane Doe", req.Name)
}
