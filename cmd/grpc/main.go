package main

import (
	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc"
	"github.com/itsLeonB/cocoon/internal/logging"
	"github.com/itsLeonB/ezutil"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	defaultConfigs := config.DefaultConfigs()
	configs := ezutil.LoadConfig(defaultConfigs)
	logging.InitLogger(configs.App)
	s := grpc.Setup(configs)
	s.Run()
}
