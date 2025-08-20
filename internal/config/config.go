package config

import (
	"log"
	"time"

	"github.com/itsLeonB/ezutil"
)

func DefaultConfigs() ezutil.Config {
	timeout, _ := time.ParseDuration("10s")
	tokenDuration, _ := time.ParseDuration("24h")
	cookieDuration, _ := time.ParseDuration("24h")
	secretKey, err := ezutil.GenerateRandomString(32)
	if err != nil {
		log.Fatal("error generating secret key: %w", err)
	}

	appConfig := ezutil.App{
		Env:        "debug",
		Port:       "8080",
		Timeout:    timeout,
		ClientUrls: []string{"http://localhost:3000"},
		Timezone:   "Asia/Jakarta",
	}

	authConfig := ezutil.Auth{
		SecretKey:      secretKey,
		TokenDuration:  tokenDuration,
		CookieDuration: cookieDuration,
		Issuer:         "cocoon",
		URL:            "http://localhost:8000",
	}

	return ezutil.Config{
		App:  &appConfig,
		Auth: &authConfig,
	}
}
