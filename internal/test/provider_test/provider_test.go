package provider_test

import (
	"testing"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/provider"
	"github.com/stretchr/testify/assert"
)

func TestProviderStructure(t *testing.T) {
	// Test that provider structure can be created
	// This is a basic structure test without database dependencies
	appConfig := config.App{
		Name: "TestApp",
		Env:  "debug",
	}

	logger := provider.ProvideLogger(appConfig)
	assert.NotNil(t, logger)
}
