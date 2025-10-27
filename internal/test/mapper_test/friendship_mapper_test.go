package mapper_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mapper"
	"github.com/itsLeonB/cocoon/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestFriendshipToResponse(t *testing.T) {
	profileID1 := uuid.New()
	profileID2 := uuid.New()
	friendshipID := uuid.New()
	now := time.Now()

	profile1 := entity.UserProfile{UserID: util.NewValidNullUUID(profileID1), Name: "User 1"}
	profile1.ID = profileID1
	profile2 := entity.UserProfile{UserID: util.NewValidNullUUID(profileID2), Name: "User 2"}
	profile2.ID = profileID2

	friendship := entity.Friendship{
		ProfileID1: profileID1,
		ProfileID2: profileID2,
		Type:       appconstant.Real,
		Profile1:   profile1,
		Profile2:   profile2,
	}
	friendship.ID = friendshipID
	friendship.CreatedAt = now
	friendship.UpdatedAt = now

	response, err := mapper.FriendshipToResponse(profileID1, friendship)

	assert.NoError(t, err)
	assert.Equal(t, friendshipID, response.ID)
	assert.Equal(t, appconstant.Real, response.Type)
	assert.Equal(t, profileID2, response.ProfileID)
	assert.Equal(t, "User 2", response.ProfileName)
	assert.Equal(t, now, response.CreatedAt)
	assert.Equal(t, now, response.UpdatedAt)
}

func TestFriendshipToResponse_InvalidProfileID(t *testing.T) {
	profileID1 := uuid.New()
	profileID2 := uuid.New()
	invalidID := uuid.New()

	friendship := entity.Friendship{
		ProfileID1: profileID1,
		ProfileID2: profileID2,
		Type:       appconstant.Real,
	}

	_, err := mapper.FriendshipToResponse(invalidID, friendship)

	assert.Error(t, err)
}

func TestOrderProfilesToFriendship(t *testing.T) {
	profileID1 := uuid.New()
	profileID2 := uuid.New()

	userProfile := dto.ProfileResponse{ID: profileID1}
	friendProfile := dto.ProfileResponse{ID: profileID2}

	friendship, err := mapper.OrderProfilesToFriendship(userProfile, friendProfile)

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, friendship.ProfileID1)
	assert.NotEqual(t, uuid.Nil, friendship.ProfileID2)
	assert.NotEqual(t, friendship.ProfileID1, friendship.ProfileID2)
}

func TestOrderProfilesToFriendship_SameID(t *testing.T) {
	profileID := uuid.New()

	userProfile := dto.ProfileResponse{ID: profileID}
	friendProfile := dto.ProfileResponse{ID: profileID}

	_, err := mapper.OrderProfilesToFriendship(userProfile, friendProfile)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "both IDs are equal")
}

func TestMapToFriendshipWithProfile(t *testing.T) {
	profileID1 := uuid.New()
	profileID2 := uuid.New()

	profile1 := entity.UserProfile{UserID: util.NewValidNullUUID(profileID1), Name: "User 1"}
	profile1.ID = profileID1
	profile2 := entity.UserProfile{UserID: util.NewValidNullUUID(profileID2), Name: "User 2"}
	profile2.ID = profileID2

	friendship := entity.Friendship{
		ProfileID1: profileID1,
		ProfileID2: profileID2,
		Type:       appconstant.Real,
		Profile1:   profile1,
		Profile2:   profile2,
	}

	result, err := mapper.MapToFriendshipWithProfile(profileID1, friendship)

	assert.NoError(t, err)
	assert.Equal(t, profileID1, result.UserProfile.ID)
	assert.Equal(t, profileID2, result.FriendProfile.ID)
	assert.Equal(t, "User 1", result.UserProfile.Name)
	assert.Equal(t, "User 2", result.FriendProfile.Name)
}

func TestMapToFriendDetails(t *testing.T) {
	profileID1 := uuid.New()
	profileID2 := uuid.New()
	friendshipID := uuid.New()
	now := time.Now()

	profile1 := entity.UserProfile{UserID: util.NewValidNullUUID(profileID1), Name: "User 1"}
	profile1.ID = profileID1
	profile2 := entity.UserProfile{UserID: util.NewValidNullUUID(profileID2), Name: "User 2"}
	profile2.ID = profileID2

	friendship := entity.Friendship{
		ProfileID1: profileID1,
		ProfileID2: profileID2,
		Type:       appconstant.Real,
		Profile1:   profile1,
		Profile2:   profile2,
	}
	friendship.ID = friendshipID
	friendship.CreatedAt = now
	friendship.UpdatedAt = now

	result, err := mapper.MapToFriendDetails(profileID1, friendship)

	assert.NoError(t, err)
	assert.Equal(t, friendshipID, result.ID)
	assert.Equal(t, profileID2, result.ProfileID)
	assert.Equal(t, "User 2", result.Name)
	assert.Equal(t, appconstant.Real, result.Type)
	assert.Equal(t, profileID1, result.ProfileID1)
	assert.Equal(t, profileID2, result.ProfileID2)
}
