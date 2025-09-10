package entity_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestFriendship_Structure(t *testing.T) {
	profileID1 := uuid.New()
	profileID2 := uuid.New()

	friendship := entity.Friendship{
		ProfileID1: profileID1,
		ProfileID2: profileID2,
		Type:       appconstant.Real,
	}

	assert.Equal(t, profileID1, friendship.ProfileID1)
	assert.Equal(t, profileID2, friendship.ProfileID2)
	assert.Equal(t, appconstant.Real, friendship.Type)
	assert.NotNil(t, friendship.Profile1)
	assert.NotNil(t, friendship.Profile2)
}

func TestFriendshipSpecification_Structure(t *testing.T) {
	spec := entity.FriendshipSpecification{
		Name: "test friendship",
	}

	assert.Equal(t, "test friendship", spec.Name)
	assert.NotNil(t, spec.Specification)
}
