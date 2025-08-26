package grpc

import (
	"github.com/itsLeonB/cocoon-protos/gen/go/auth/v1"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile/v1"
	"google.golang.org/grpc"
)

func (s *Server) register(grpcServer *grpc.Server) {
	auth.RegisterAuthServiceServer(grpcServer, s.servers.Auth)
	profile.RegisterProfileServiceServer(grpcServer, s.servers.Profile)
	friendship.RegisterFriendshipServiceServer(grpcServer, s.servers.Friendship)
}
