package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/cocoon-protos/gen/go/auth"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile"
	"github.com/itsLeonB/cocoon/internal/provider"
)

type Servers struct {
	Auth       auth.AuthServiceServer
	Profile    profile.ProfileServiceServer
	Friendship friendship.FriendshipServiceServer
}

func ProvideServers(services *provider.Services) *Servers {
	validate := validator.New()

	return &Servers{
		Auth:       NewAuthServer(validate, services.Auth),
		Profile:    NewProfileServer(validate, services.Profile),
		Friendship: NewFriendshipServer(validate, services.Friendship),
	}
}
