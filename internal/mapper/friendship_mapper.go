package mapper

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/helper"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
)

func FriendshipToResponse(userProfileID uuid.UUID, friendship entity.Friendship) (dto.FriendshipResponse, error) {
	_, friendProfile, err := helper.SelectProfiles(userProfileID, friendship)
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	return dto.FriendshipResponse{
		ID:          friendship.ID,
		Type:        friendship.Type,
		ProfileID:   friendProfile.ID,
		ProfileName: friendProfile.Name,
		CreatedAt:   friendship.CreatedAt,
		UpdatedAt:   friendship.UpdatedAt,
		DeletedAt:   friendship.DeletedAt.Time,
	}, nil
}

func OrderProfilesToFriendship(userProfile, friendProfile dto.ProfileResponse) (entity.Friendship, error) {
	switch ezutil.CompareUUID(userProfile.ID, friendProfile.ID) {
	case 1:
		return entity.Friendship{
			ProfileID1: friendProfile.ID,
			ProfileID2: userProfile.ID,
		}, nil
	case -1:
		return entity.Friendship{
			ProfileID1: userProfile.ID,
			ProfileID2: friendProfile.ID,
		}, nil
	default:
		return entity.Friendship{}, eris.New("both IDs are equal, cannot create friendship")
	}
}

func MapToFriendshipWithProfile(userProfileID uuid.UUID, friendship entity.Friendship) (dto.FriendshipWithProfile, error) {
	friendshipResponse, err := FriendshipToResponse(userProfileID, friendship)
	if err != nil {
		return dto.FriendshipWithProfile{}, err
	}

	userProfile, friendProfile, err := helper.SelectProfiles(userProfileID, friendship)
	if err != nil {
		return dto.FriendshipWithProfile{}, err
	}

	return dto.FriendshipWithProfile{
		Friendship:    friendshipResponse,
		UserProfile:   ProfileToResponse(userProfile),
		FriendProfile: ProfileToResponse(friendProfile),
	}, nil
}

func MapToFriendDetails(userProfileID uuid.UUID, friendship entity.Friendship) (dto.FriendDetails, error) {
	friendshipWithProfile, err := MapToFriendshipWithProfile(userProfileID, friendship)
	if err != nil {
		return dto.FriendDetails{}, err
	}

	friendProfile := friendshipWithProfile.FriendProfile

	return dto.FriendDetails{
		ID:         friendship.ID,
		ProfileID:  friendProfile.ID,
		Name:       friendProfile.Name,
		Type:       friendship.Type,
		CreatedAt:  friendship.CreatedAt,
		UpdatedAt:  friendship.UpdatedAt,
		DeletedAt:  friendship.DeletedAt.Time,
		ProfileID1: friendship.ProfileID1,
		ProfileID2: friendship.ProfileID2,
	}, nil
}
