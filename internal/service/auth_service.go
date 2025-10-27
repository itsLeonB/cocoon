package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mapper"
	"github.com/itsLeonB/cocoon/internal/util"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/sekure"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
)

type authServiceImpl struct {
	hashService sekure.HashService
	jwtService  sekure.JWTService
	transactor  crud.Transactor
	userSvc     UserService
	mailSvc     MailService
}

func NewAuthService(
	transactor crud.Transactor,
	configs config.Auth,
	userSvc UserService,
	mailSvc MailService,
) AuthService {
	return &authServiceImpl{
		sekure.NewHashService(configs.HashCost),
		sekure.NewJwtService(configs.Issuer, configs.SecretKey, configs.TokenDuration),
		transactor,
		userSvc,
		mailSvc,
	}
}

func (as *authServiceImpl) Register(ctx context.Context, request dto.RegisterRequest) (bool, error) {
	isVerified := request.VerificationURL == ""
	err := as.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		existingUser, err := as.userSvc.FindByEmail(ctx, request.Email)
		if err != nil {
			return err
		}
		if !existingUser.IsZero() {
			return ungerr.ConflictError(fmt.Sprintf(appconstant.ErrAuthDuplicateUser, request.Email))
		}

		hash, err := as.hashService.Hash(request.Password)
		if err != nil {
			return err
		}

		newUserReq := dto.NewUserRequest{
			Email:     request.Email,
			Password:  hash,
			Name:      util.GetNameFromEmail(request.Email),
			VerifyNow: isVerified,
		}

		user, err := as.userSvc.CreateNew(ctx, newUserReq)
		if err != nil {
			return err
		}
		if isVerified {
			return nil
		}

		return as.sendVerificationMail(ctx, user, request.VerificationURL)
	})
	return isVerified, err
}

func (as *authServiceImpl) VerifyRegistration(ctx context.Context, token string) (dto.LoginResponse, error) {
	var response dto.LoginResponse
	err := as.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		claims, err := as.jwtService.VerifyToken(token)
		if err != nil {
			return err
		}
		id, ok := claims.Data["id"].(string)
		if !ok {
			return eris.New("error asserting id, is not a string")
		}
		userID, err := ezutil.Parse[uuid.UUID](id)
		if err != nil {
			return err
		}
		email, ok := claims.Data["email"].(string)
		if !ok {
			return eris.New("error asserting email, is not a string")
		}
		unixTime, ok := claims.Data["exp"].(int64)
		if !ok {
			return eris.New("error asserting exp, is not an int64")
		}
		if time.Now().Unix() > unixTime {
			return ungerr.UnauthorizedError("token has expired")
		}

		user, err := as.userSvc.Verify(ctx, userID, email, util.GetNameFromEmail(email), "")
		if err != nil {
			return err
		}

		response, err = as.createLoginResponse(user)
		return err
	})
	return response, err
}

func (as *authServiceImpl) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := as.userSvc.FindByEmail(ctx, request.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if user.IsZero() {
		return dto.LoginResponse{}, ungerr.NotFoundError(appconstant.ErrAuthUnknownCredentials)
	}
	if !user.IsVerified() {
		return dto.LoginResponse{}, ungerr.NotFoundError(appconstant.ErrAuthUnknownCredentials)
	}

	ok, err := as.hashService.CheckHash(user.Password, request.Password)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if !ok {
		return dto.LoginResponse{}, ungerr.NotFoundError(appconstant.ErrAuthUnknownCredentials)
	}

	return as.createLoginResponse(user)
}

func (as *authServiceImpl) VerifyToken(ctx context.Context, token string) (dto.AuthData, error) {
	claims, err := as.jwtService.VerifyToken(token)
	if err != nil {
		return dto.AuthData{}, err
	}

	tokenUserId, exists := claims.Data[appconstant.ContextUserID]
	if !exists {
		return dto.AuthData{}, eris.New("missing user ID from token")
	}
	stringUserID, ok := tokenUserId.(string)
	if !ok {
		return dto.AuthData{}, eris.New("error asserting userID, is not a string")
	}
	userID, err := ezutil.Parse[uuid.UUID](stringUserID)
	if err != nil {
		return dto.AuthData{}, err
	}

	user, err := as.userSvc.GetByID(ctx, userID)
	if err != nil {
		return dto.AuthData{}, err
	}

	return dto.AuthData{
		ProfileID: user.Profile.ID,
	}, nil
}

func (as *authServiceImpl) SendResetPassword(ctx context.Context, resetURL, email string) error {
	return as.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		user, err := as.userSvc.FindByEmail(ctx, email)
		if err != nil {
			return err
		}
		if user.IsZero() || user.IsDeleted() || !user.IsVerified() {
			return nil
		}

		resetToken, err := as.userSvc.GeneratePasswordResetToken(ctx, user.ID)
		if err != nil {
			return err
		}

		return as.sendResetPasswordMail(ctx, user, resetURL, resetToken)
	})
}

func (as *authServiceImpl) ResetPassword(ctx context.Context, token, newPassword string) (dto.LoginResponse, error) {
	var response dto.LoginResponse
	err := as.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		claims, err := as.jwtService.VerifyToken(token)
		if err != nil {
			return err
		}
		id, ok := claims.Data["id"].(string)
		if !ok {
			return eris.New("error asserting id, is not a string")
		}
		userID, err := ezutil.Parse[uuid.UUID](id)
		if err != nil {
			return err
		}
		email, ok := claims.Data["email"].(string)
		if !ok {
			return eris.New("error asserting email, is not a string")
		}
		resetToken, ok := claims.Data["reset_token"].(string)
		if !ok {
			return eris.New("error asserting reset_token, is not a string")
		}

		hashedPassword, err := as.hashService.Hash(newPassword)
		if err != nil {
			return err
		}

		user, err := as.userSvc.ResetPassword(ctx, userID, email, resetToken, hashedPassword)
		if err != nil {
			return err
		}

		response, err = as.createLoginResponse(user)
		return err
	})
	return response, err
}

func (as *authServiceImpl) createLoginResponse(user entity.User) (dto.LoginResponse, error) {
	authData := mapper.UserToAuthData(user)

	token, err := as.jwtService.CreateToken(authData)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.NewBearerTokenResp(token), nil
}

func (as *authServiceImpl) sendVerificationMail(ctx context.Context, user entity.User, verificationURL string) error {
	claims := map[string]any{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(30 * time.Minute).Unix(),
	}

	token, err := as.jwtService.CreateToken(claims)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s?token=%s", verificationURL, token)

	mailMsg := dto.MailMessage{
		RecipientMail: user.Email,
		RecipientName: util.GetNameFromEmail(user.Email),
		Subject:       "Verify your email",
		TextContent:   "Please verify your email by clicking the following link:\n\n" + url,
	}

	return as.mailSvc.Send(ctx, mailMsg)
}

func (as *authServiceImpl) sendResetPasswordMail(ctx context.Context, user entity.User, resetURL, resetToken string) error {
	claims := map[string]any{
		"id":          user.ID,
		"email":       user.Email,
		"reset_token": resetToken,
	}

	token, err := as.jwtService.CreateToken(claims)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s?token=%s", resetURL, token)

	mailMsg := dto.MailMessage{
		RecipientMail: user.Email,
		RecipientName: user.Profile.Name,
		Subject:       "Reset your password",
		TextContent:   "You have requested to reset your password.\nIf this is not you, ignore this mail.\nPlease reset your password by clicking the following link:\n\n" + url,
	}

	return as.mailSvc.Send(ctx, mailMsg)
}
