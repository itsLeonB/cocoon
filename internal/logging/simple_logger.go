package logging

import (
	"github.com/itsLeonB/ezutil"
)

var Logger ezutil.Logger

func InitLogger(configs *ezutil.App) {
	minLevel := 0
	if configs.Env == "release" {
		minLevel = 1
	}

	Logger = ezutil.NewSimpleLogger("Cocoon", true, minLevel)
}
