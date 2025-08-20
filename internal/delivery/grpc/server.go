package grpc

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsLeonB/cocoon/internal/delivery/grpc/interceptor"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc/server"
	"github.com/itsLeonB/cocoon/internal/logging"
	"github.com/itsLeonB/cocoon/internal/provider"
	"github.com/itsLeonB/ezutil"
	"google.golang.org/grpc"
)

type Server struct {
	address string
	opts    []grpc.ServerOption
	servers *server.Servers
}

func Setup(configs *ezutil.Config) *Server {
	repos := provider.ProvideRepositories(configs.GORM)
	services := provider.ProvideServices(configs, repos)
	servers := server.ProvideServers(services)

	// Middlewares/Interceptors
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.NewLoggerInterceptor(),
			interceptor.NewErrorInterceptor(),
		),
	}

	return &Server{
		address: ":" + configs.App.Port,
		opts:    opts,
		servers: servers,
	}
}

func (s *Server) Run() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		logging.Logger.Fatalf("error listening to %s: %v", s.address, err)
	}

	grpcServer := grpc.NewServer(s.opts...)
	s.register(grpcServer)

	go func() {
		logging.Logger.Infof("server started at: %s", s.address)
		if err := grpcServer.Serve(listener); err != nil {
			logging.Logger.Fatalf("failed to serve: %v", err)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	<-exit
	logging.Logger.Info("shutting down server...")
	grpcServer.GracefulStop()
	logging.Logger.Info("server successfully shutdown")
}
