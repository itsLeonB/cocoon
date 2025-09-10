package dto_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestNewAnonymousFriendshipRequest_Structure(t *testing.T) {
	profileID := uuid.New()
	req := dto.NewAnonymousFriendshipRequest{
		ProfileID: profileID,
		Name:      "Anonymous Friend",
	}

	assert.Equal(t, profileID, req.ProfileID)
	assert.Equal(t, "Anonymous Friend", req.Name)
}

func TestFriendshipResponse_Structure(t *testing.T) {
	id := uuid.New()
	profileID := uuid.New()
	now := time.Now()

	resp := dto.FriendshipResponse{
		ID:          id,
		Type:        appconstant.Real,
		ProfileID:   profileID,
		ProfileName: "Friend Name",
		CreatedAt:   now,
		UpdatedAt:   now,
		DeletedAt:   time.Time{},
	}

	assert.Equal(t, id, resp.ID)
	assert.Equal(t, appconstant.Real, resp.Type)
	assert.Equal(t, profileID, resp.ProfileID)
	assert.Equal(t, "Friend Name", resp.ProfileName)
	assert.Equal(t, now, resp.CreatedAt)
	assert.Equal(t, now, resp.UpdatedAt)
	assert.Equal(t, time.Time{}, resp.DeletedAt)
}

func TestFriendshipWithProfile_Structure(t *testing.T) {
	friendship := dto.FriendshipResponse{ID: uuid.New()}
	userProfile := dto.ProfileResponse{ID: uuid.New()}
	friendProfile := dto.ProfileResponse{ID: uuid.New()}

	resp := dto.FriendshipWithProfile{
		Friendship:    friendship,
		UserProfile:   userProfile,
		FriendProfile: friendProfile,
	}

	assert.Equal(t, friendship, resp.Friendship)
	assert.Equal(t, userProfile, resp.UserProfile)
	assert.Equal(t, friendProfile, resp.FriendProfile)
}

func TestFriendDetails_Structure(t *testing.T) {
	id := uuid.New()
	profileID := uuid.New()
	profileID1 := uuid.New()
	profileID2 := uuid.New()
	now := time.Now()

	details := dto.FriendDetails{
		ID:         id,
		ProfileID:  profileID,
		Name:       "Friend Name",
		Type:       appconstant.Anonymous,
		Email:      "friend@example.com",
		Phone:      "+1234567890",
		Avatar:     "avatar.jpg",
		CreatedAt:  now,
		UpdatedAt:  now,
		DeletedAt:  time.Time{},
		ProfileID1: profileID1,
		ProfileID2: profileID2,
	}

	assert.Equal(t, id, details.ID)
	assert.Equal(t, profileID, details.ProfileID)
	assert.Equal(t, "Friend Name", details.Name)
	assert.Equal(t, appconstant.Anonymous, details.Type)
	assert.Equal(t, "friend@example.com", details.Email)
	assert.Equal(t, "+1234567890", details.Phone)
	assert.Equal(t, "avatar.jpg", details.Avatar)
	assert.Equal(t, now, details.CreatedAt)
	assert.Equal(t, now, details.UpdatedAt)
	assert.Equal(t, time.Time{}, details.DeletedAt)
	assert.Equal(t, profileID1, details.ProfileID1)
	assert.Equal(t, profileID2, details.ProfileID2)
}
