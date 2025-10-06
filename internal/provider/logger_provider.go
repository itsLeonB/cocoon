package provider

import (
	"github.com/itsLeonB/ezutil/v2"
)

func ProvideLogger(appName, env string) ezutil.Logger {
	minLevel := 0
	if env == "release" {
		minLevel = 1
	}

	return ezutil.NewSimpleLogger(appName, true, minLevel)
}
