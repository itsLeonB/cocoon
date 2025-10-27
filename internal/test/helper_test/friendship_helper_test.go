package helper_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/helper"
	"github.com/itsLeonB/cocoon/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestSelectProfiles(t *testing.T) {
	profileID1 := uuid.New()
	profileID2 := uuid.New()

	profile1 := entity.UserProfile{UserID: util.NewValidNullUUID(profileID1), Name: "User 1"}
	profile2 := entity.UserProfile{UserID: util.NewValidNullUUID(profileID2), Name: "User 2"}

	friendship := entity.Friendship{
		ProfileID1: profileID1,
		ProfileID2: profileID2,
		Profile1:   profile1,
		Profile2:   profile2,
	}

	tests := []struct {
		name            string
		userProfileID   uuid.UUID
		expectedProfile entity.UserProfile
		expectedFriend  entity.UserProfile
		expectError     bool
	}{
		{
			name:            "select profile1 as user",
			userProfileID:   profileID1,
			expectedProfile: profile1,
			expectedFriend:  profile2,
			expectError:     false,
		},
		{
			name:            "select profile2 as user",
			userProfileID:   profileID2,
			expectedProfile: profile2,
			expectedFriend:  profile1,
			expectError:     false,
		},
		{
			name:          "mismatched profile ID",
			userProfileID: uuid.New(),
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userProfile, friendProfile, err := helper.SelectProfiles(tt.userProfileID, friendship)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, userProfile)
				assert.Empty(t, friendProfile)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedProfile, userProfile)
				assert.Equal(t, tt.expectedFriend, friendProfile)
			}
		})
	}
}
