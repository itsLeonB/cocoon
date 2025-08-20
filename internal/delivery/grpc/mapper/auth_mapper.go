package mapper

import (
	"strings"

	"github.com/itsLeonB/cocoon/gen/go/auth"
	"github.com/itsLeonB/cocoon/internal/dto"
)

func FromRegisterRequestProto(req *auth.RegisterRequest) dto.RegisterRequest {
	return dto.RegisterRequest{
		Email:                strings.ToLower(req.GetEmail()),
		Password:             req.GetPassword(),
		PasswordConfirmation: req.GetPasswordConfirmation(),
	}
}

func FromLoginRequestProto(req *auth.LoginRequest) dto.LoginRequest {
	return dto.LoginRequest{
		Email:    strings.ToLower(req.GetEmail()),
		Password: req.GetPassword(),
	}
}

func ToLoginResponseProto(resp dto.LoginResponse) *auth.LoginResponse {
	return &auth.LoginResponse{
		Type:  resp.Type,
		Token: resp.Token,
	}
}
