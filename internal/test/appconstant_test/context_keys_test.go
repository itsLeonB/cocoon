package appconstant_test

import (
	"testing"

	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/stretchr/testify/assert"
)

func TestContextKeys_Constants(t *testing.T) {
	assert.Equal(t, "userID", appconstant.ContextUserID)
	assert.Equal(t, "profileID", appconstant.ContextProfileID)
}
