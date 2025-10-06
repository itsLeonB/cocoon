package main

import (
	"log"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc"
	_ "github.com/joho/godotenv/autoload"
	"github.com/rotisserie/eris"
)

func main() {
	configs, err := config.Load()
	if err != nil {
		log.Fatal(eris.ToString(err, true))
	}
	s, err := grpc.Setup(configs)
	if err != nil {
		log.Fatal(eris.ToString(err, true))
	}
	s.Run()
}
