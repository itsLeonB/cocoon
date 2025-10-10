package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/itsLeonB/cocoon-protos/gen/go/auth/v1"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
	"google.golang.org/protobuf/types/known/emptypb"
)

type authServer struct {
	auth.UnimplementedAuthServiceServer
	validate    *validator.Validate
	authService service.AuthService
	oAuthSvc    service.OAuthService
}

func newAuthServer(
	validate *validator.Validate,
	authService service.AuthService,
	oAuthSvc service.OAuthService,
) auth.AuthServiceServer {
	return &authServer{
		validate:    validate,
		authService: authService,
		oAuthSvc:    oAuthSvc,
	}
}

func (as *authServer) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	request := dto.RegisterRequest{
		Email:                req.GetEmail(),
		Password:             req.GetPassword(),
		PasswordConfirmation: req.GetPasswordConfirmation(),
		VerificationURL:      req.GetVerificationUrl(),
	}

	if err := as.validate.Struct(request); err != nil {
		return nil, err
	}

	isVerified, err := as.authService.Register(ctx, request)
	if err != nil {
		return nil, err
	}

	return &auth.RegisterResponse{
		IsVerified: isVerified,
	}, nil
}

func (as *authServer) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
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

func (as *authServer) handleOAuth2Login(ctx context.Context, req *auth.OAuth2LoginRequest) (*auth.LoginResponse, error) {
	data := dto.OAuthCallbackData{
		Provider: req.GetProvider(),
		Code:     req.GetCode(),
		State:    req.GetState(),
	}
	if err := as.validate.Struct(data); err != nil {
		return nil, eris.Wrap(err, appconstant.ErrStructValidation)
	}

	response, err := as.oAuthSvc.HandleOAuthCallback(ctx, data)
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Type:  response.Type,
		Token: response.Token,
	}, nil
}

func (as *authServer) handleInternalLogin(ctx context.Context, req *auth.InternalLoginRequest) (*auth.LoginResponse, error) {
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

func (as *authServer) VerifyToken(ctx context.Context, req *auth.VerifyTokenRequest) (*auth.VerifyTokenResponse, error) {
	authData, err := as.authService.VerifyToken(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return &auth.VerifyTokenResponse{
		ProfileId: authData.ProfileID.String(),
	}, nil
}

func (as *authServer) GetOAuth2Url(ctx context.Context, req *auth.GetOAuth2UrlRequest) (*auth.GetOAuth2UrlResponse, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	if req.GetProvider() == "" {
		return nil, ungerr.BadRequestError("provider is empty")
	}

	url, err := as.oAuthSvc.GetOAuthURL(ctx, req.GetProvider())
	if err != nil {
		return nil, err
	}

	return &auth.GetOAuth2UrlResponse{
		Url: url,
	}, nil
}

func (as *authServer) VerifyRegistration(ctx context.Context, req *auth.VerifyRegistrationRequest) (*auth.LoginResponse, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	if req.GetToken() == "" {
		return nil, ungerr.BadRequestError("token is empty")
	}

	response, err := as.authService.VerifyRegistration(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Type:  response.Type,
		Token: response.Token,
	}, nil
}

func (as *authServer) SendResetPassword(ctx context.Context, req *auth.SendResetPasswordRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	if req.GetEmail() == "" {
		return nil, ungerr.BadRequestError("email is empty")
	}
	if req.GetResetUrl() == "" {
		return nil, ungerr.BadRequestError("resetUrl is empty")
	}
	return nil, as.authService.SendResetPassword(ctx, req.GetResetUrl(), req.GetEmail())
}

func (as *authServer) ResetPassword(ctx context.Context, req *auth.ResetPasswordRequest) (*auth.LoginResponse, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	if req.GetToken() == "" {
		return nil, ungerr.BadRequestError("token is empty")
	}
	if req.GetNewPassword() == "" {
		return nil, ungerr.BadRequestError("newPassword is empty")
	}

	response, err := as.authService.ResetPassword(ctx, req.GetToken(), req.GetNewPassword())
	if err != nil {
		return nil, err
	}

	return &auth.LoginResponse{
		Type:  response.Type,
		Token: response.Token,
	}, nil
}
