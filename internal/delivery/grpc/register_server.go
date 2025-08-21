package grpc

import (
	"github.com/itsLeonB/cocoon-protos/gen/go/auth"
	"google.golang.org/grpc"
)

func (s *Server) register(grpcServer *grpc.Server) {
	auth.RegisterAuthServiceServer(grpcServer, s.servers.Auth)
}
