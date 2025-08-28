package provider

import (
	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/ezutil/v2"
)

func ProvideLogger(configs config.App) ezutil.Logger {
	minLevel := 0
	if configs.Env == "release" {
		minLevel = 1
	}

	return ezutil.NewSimpleLogger(configs.Name, true, minLevel)
}
