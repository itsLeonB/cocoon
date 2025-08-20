package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/cocoon/gen/go/auth"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc/mapper"
	"github.com/itsLeonB/cocoon/internal/service"
)

type AuthServer struct {
	auth.UnimplementedAuthServiceServer
	validate    *validator.Validate
	authService service.AuthService
}

func NewAuthServer(
	validate *validator.Validate,
	authService service.AuthService,
) *AuthServer {
	return &AuthServer{
		validate:    validate,
		authService: authService,
	}
}

func (as *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	request := mapper.FromRegisterRequestProto(req)

	if err := as.validate.Struct(request); err != nil {
		return nil, err
	}

	if err := as.authService.Register(ctx, request); err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		Message: "success",
	}, nil
}

func (as *AuthServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	request := mapper.FromLoginRequestProto(req)

	if err := as.validate.Struct(request); err != nil {
		return nil, err
	}

	response, err := as.authService.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	return mapper.ToLoginResponseProto(response), nil
}
