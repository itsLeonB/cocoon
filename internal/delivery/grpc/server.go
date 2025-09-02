package grpc

import (
	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc/interceptor"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc/server"
	"github.com/itsLeonB/cocoon/internal/provider"
	"github.com/itsLeonB/gerpc"
	"google.golang.org/grpc"
)

func Setup(configs config.Config) *gerpc.GrpcServer {
	providers := provider.All(configs)
	servers := server.ProvideServers(providers.Services)

	// Middlewares/Interceptors
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.NewLoggerInterceptor(providers.Logger),
			interceptor.NewErrorInterceptor(providers.Logger),
		),
	}

	return gerpc.NewGrpcServer().
		WithLogger(providers.Logger).
		WithAddress(":" + configs.App.Port).
		WithOpts(opts...).
		WithRegisterSrvFunc(servers.Register).
		WithShutdownFunc(providers.Shutdown)
}
