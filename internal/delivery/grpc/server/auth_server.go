package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/cocoon-protos/gen/go/auth/v1"
	"github.com/itsLeonB/cocoon/internal/dto"
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
) auth.AuthServiceServer {
	return &AuthServer{
		validate:    validate,
		authService: authService,
	}
}

func (as *AuthServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	request := dto.RegisterRequest{
		Email:                req.GetEmail(),
		Password:             req.GetPassword(),
		PasswordConfirmation: req.GetPasswordConfirmation(),
	}

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
	request := dto.LoginRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	if err := as.validate.Struct(request); err != nil {
		return nil, err
	}

	response, err := as.authService.Login(ctx, request)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Type:  response.Type,
		Token: response.Token,
	}, nil
}

func (as *AuthServer) VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest) (*auth.VerifyTokenResponse, error) {
	authData, err := as.authService.VerifyToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return &auth.VerifyTokenResponse{
		ProfileId: authData.ProfileID.String(),
	}, nil
}
