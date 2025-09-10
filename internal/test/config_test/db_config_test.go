package config_test

import (
	"testing"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestDB_Structure(t *testing.T) {
	db := config.DB{
		Driver:   "postgres",
		Host:     "localhost",
		Port:     "5432",
		User:     "testuser",
		Password: "testpass",
		Name:     "testdb",
	}

	assert.Equal(t, "postgres", db.Driver)
	assert.Equal(t, "localhost", db.Host)
	assert.Equal(t, "5432", db.Port)
	assert.Equal(t, "testuser", db.User)
	assert.Equal(t, "testpass", db.Password)
	assert.Equal(t, "testdb", db.Name)
}
