package repository_test

import (
	"testing"

	"github.com/itsLeonB/cocoon/internal/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestProfileRepository_Constructor(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	repo := repository.NewProfileRepository(db)
	assert.NotNil(t, repo)
}
