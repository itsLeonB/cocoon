package main

import (
	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	configs := config.Load()
	s := grpc.Setup(configs)
	s.Run()
}
