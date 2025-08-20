package helper

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/rotisserie/eris"
)

func SelectProfiles(userProfileID uuid.UUID, friendship entity.Friendship) (entity.UserProfile, entity.UserProfile, error) {
	switch userProfileID {
	case friendship.ProfileID1:
		return friendship.Profile1, friendship.Profile2, nil
	case friendship.ProfileID2:
		return friendship.Profile2, friendship.Profile1, nil
	default:
		return entity.UserProfile{}, entity.UserProfile{}, eris.New(fmt.Sprintf(
			"mismatched user profile ID: %s with friendship ID: %s",
			userProfileID,
			friendship.ID,
		))
	}
}
