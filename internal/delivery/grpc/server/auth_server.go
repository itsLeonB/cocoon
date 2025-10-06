package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/cocoon-protos/gen/go/auth/v1"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
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
	if req == nil {
		return nil, eris.New("request is nil")
	}
	switch request := req.GetLoginMethod().(type) {
	case *auth.LoginRequest_InternalRequest:
		return as.handleInternalLogin(ctx, request.InternalRequest)
	case *auth.LoginRequest_Oauth2Request:
		return as.handleOAuth2Login(ctx, request.Oauth2Request)
	default:
		return nil, eris.Errorf("unsupported login method: %T", request)
	}
}

func (as *AuthServer) handleOAuth2Login(ctx context.Context, req *auth.OAuth2LoginRequest) (*auth.LoginResponse, error) {
	if req.GetProvider() == "" {
		return nil, ungerr.BadRequestError("provider is empty")
	}
	if req.GetCode() == "" {
		return nil, ungerr.BadRequestError("code is empty")
	}
	if req.GetState() == "" {
		return nil, ungerr.BadRequestError("state is empty")
	}

	response, err := as.authService.HandleOAuthCallback(ctx, req.GetProvider(), req.GetCode(), req.GetState())
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Type:  response.Type,
		Token: response.Token,
	}, nil
}

func (as *AuthServer) handleInternalLogin(ctx context.Context, req *auth.InternalLoginRequest) (*auth.LoginResponse, error) {
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

func (as *AuthServer) GetOAuth2Url(ctx context.Context, req *auth.GetOAuth2UrlRequest) (*auth.GetOAuth2UrlResponse, error) {
	if req == nil {
		return nil, eris.New("request is nil")
	}
	if req.GetProvider() == "" {
		return nil, ungerr.BadRequestError("provider is empty")
	}

	url, err := as.authService.GetOAuthURL(ctx, req.GetProvider())
	if err != nil {
		return nil, err
	}

	return &auth.GetOAuth2UrlResponse{
		Url: url,
	}, nil
}
