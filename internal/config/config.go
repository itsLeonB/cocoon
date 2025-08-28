package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App
	Auth
	DB
}

type App struct {
	Name    string        `default:"Cocoon"`
	Env     string        `default:"debug"`
	Port    string        `default:"50051"`
	Timeout time.Duration `default:"10s"`
}

type Auth struct {
	SecretKey     string        `split_words:"true" default:"thisissecret"`
	TokenDuration time.Duration `split_words:"true" default:"24h"`
	Issuer        string        `default:"cocoon"`
	HashCost      int           `split_words:"true" default:"10"`
}

func Load() Config {
	var app App
	envconfig.MustProcess("APP", &app)

	var auth Auth
	envconfig.MustProcess("AUTH", &auth)

	var db DB
	envconfig.MustProcess("DB", &db)

	return Config{
		App:  app,
		Auth: auth,
		DB:   db,
	}
}
