package provider

import (
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/ezutil/v2"
)

func ProvideLogger(appName string, env appconstant.Environment) ezutil.Logger {
	minLevel := 0
	if env == appconstant.EnvProd {
		minLevel = 1
	}

	return ezutil.NewSimpleLogger(appName, true, minLevel)
}
