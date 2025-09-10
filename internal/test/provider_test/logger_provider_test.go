package provider_test

import (
	"testing"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/provider"
	"github.com/stretchr/testify/assert"
)

func TestProvideLogger_DebugMode(t *testing.T) {
	appConfig := config.App{
		Name: "TestApp",
		Env:  "debug",
	}

	logger := provider.ProvideLogger(appConfig)

	assert.NotNil(t, logger)
}

func TestProvideLogger_ReleaseMode(t *testing.T) {
	appConfig := config.App{
		Name: "TestApp",
		Env:  "release",
	}

	logger := provider.ProvideLogger(appConfig)

	assert.NotNil(t, logger)
}

func TestProvideLogger_DifferentNames(t *testing.T) {
	appConfig1 := config.App{
		Name: "App1",
		Env:  "debug",
	}
	appConfig2 := config.App{
		Name: "App2",
		Env:  "debug",
	}

	logger1 := provider.ProvideLogger(appConfig1)
	logger2 := provider.ProvideLogger(appConfig2)

	assert.NotNil(t, logger1)
	assert.NotNil(t, logger2)
}
