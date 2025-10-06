package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rotisserie/eris"
)

type Config struct {
	App
	Auth
	DB
	OAuthProviders
	Valkey
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

func Load() (Config, error) {
	errMsg := "error loading config"

	var app App
	if err := envconfig.Process("APP", &app); err != nil {
		return Config{}, eris.Wrap(err, errMsg)
	}

	var auth Auth
	if err := envconfig.Process("AUTH", &auth); err != nil {
		return Config{}, eris.Wrap(err, errMsg)
	}

	var db DB
	if err := envconfig.Process("DB", &db); err != nil {
		return Config{}, eris.Wrap(err, errMsg)
	}

	oAuthProviders, err := loadOAuthProviderConfig()
	if err != nil {
		return Config{}, eris.Wrap(err, errMsg)
	}

	var valkey Valkey
	if err = envconfig.Process("VALKEY", &valkey); err != nil {
		return Config{}, eris.Wrap(err, errMsg)
	}

	return Config{
		App:            app,
		Auth:           auth,
		DB:             db,
		OAuthProviders: oAuthProviders,
		Valkey:         valkey,
	}, nil
}
