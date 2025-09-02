package server

import (
	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/cocoon-protos/gen/go/auth/v1"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile/v1"
	"github.com/itsLeonB/cocoon/internal/provider"
	"google.golang.org/grpc"
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

func (s *Servers) Register(grpcServer *grpc.Server) error {
	auth.RegisterAuthServiceServer(grpcServer, s.Auth)
	profile.RegisterProfileServiceServer(grpcServer, s.Profile)
	friendship.RegisterFriendshipServiceServer(grpcServer, s.Friendship)
	return nil
}
