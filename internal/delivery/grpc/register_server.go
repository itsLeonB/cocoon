package grpc

import (
	"github.com/itsLeonB/cocoon-protos/gen/go/auth"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile"
	"google.golang.org/grpc"
)

func (s *Server) register(grpcServer *grpc.Server) {
	auth.RegisterAuthServiceServer(grpcServer, s.servers.Auth)
	profile.RegisterProfileServiceServer(grpcServer, s.servers.Profile)
	friendship.RegisterFriendshipServiceServer(grpcServer, s.servers.Friendship)
}
