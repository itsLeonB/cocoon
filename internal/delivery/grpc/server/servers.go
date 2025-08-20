package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/cocoon/internal/provider"
)

type Servers struct {
	Auth *AuthServer
}

func ProvideServers(services *provider.Services) *Servers {
	validate := validator.New()

	return &Servers{
		Auth: NewAuthServer(validate, services.Auth),
	}
}
