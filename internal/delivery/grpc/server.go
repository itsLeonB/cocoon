package grpc

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc/interceptor"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc/server"
	"github.com/itsLeonB/cocoon/internal/provider"
	"github.com/itsLeonB/ezutil/v2"
	"google.golang.org/grpc"
)

type Server struct {
	logger       ezutil.Logger
	address      string
	opts         []grpc.ServerOption
	servers      *server.Servers
	shutdownFunc func() error
}

func Setup(configs config.Config) *Server {
	providers := provider.All(configs)
	servers := server.ProvideServers(providers.Services)

	// Middlewares/Interceptors
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.NewLoggerInterceptor(providers.Logger),
			interceptor.NewErrorInterceptor(providers.Logger),
		),
	}

	return &Server{
		logger:       providers.Logger,
		address:      ":" + configs.App.Port,
		opts:         opts,
		servers:      servers,
		shutdownFunc: providers.Shutdown,
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		s.logger.Fatalf("error listening to %s: %v", s.address, err)
	}

	grpcServer := grpc.NewServer(s.opts...)
	s.register(grpcServer)

	go func() {
		s.logger.Infof("server started at: %s", s.address)
		if err := grpcServer.Serve(listener); err != nil {
			s.logger.Fatalf("failed to serve: %v", err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	s.logger.Info("shutting down server...")
	grpcServer.GracefulStop()

	s.logger.Info("initating cleanup")
	if err := s.shutdownFunc(); err != nil {
		s.logger.Errorf("error during cleanup: %v", err)
	}

	s.logger.Info("server successfully shutdown")
}
