package appconstant_test

import (
	"testing"

	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/stretchr/testify/assert"
)

func TestFriendshipType_Constants(t *testing.T) {
	assert.Equal(t, appconstant.FriendshipType("REAL"), appconstant.Real)
	assert.Equal(t, appconstant.FriendshipType("ANON"), appconstant.Anonymous)
}

func TestFriendshipType_Values(t *testing.T) {
	assert.Equal(t, "REAL", string(appconstant.Real))
	assert.Equal(t, "ANON", string(appconstant.Anonymous))
}
