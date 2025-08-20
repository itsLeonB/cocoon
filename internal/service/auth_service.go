package service

import (
	"context"
	"fmt"

	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mapper"
	"github.com/itsLeonB/cocoon/internal/repository"
	"github.com/itsLeonB/cocoon/internal/util"
	"github.com/itsLeonB/ezutil"
)

type authServiceImpl struct {
	hashService           ezutil.HashService
	jwtService            ezutil.JWTService
	userRepository        repository.UserRepository
	transactor            ezutil.Transactor
	userProfileRepository repository.UserProfileRepository
}

func NewAuthService(
	hashService ezutil.HashService,
	jwtService ezutil.JWTService,
	userRepository repository.UserRepository,
	transactor ezutil.Transactor,
	userProfileRepository repository.UserProfileRepository,
) AuthService {
	return &authServiceImpl{
		hashService,
		jwtService,
		userRepository,
		transactor,
		userProfileRepository,
	}
}

func (as *authServiceImpl) Register(ctx context.Context, request dto.RegisterRequest) error {
	return as.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := ezutil.Specification[entity.User]{}
		spec.Model.Email = request.Email

		existingUser, err := as.userRepository.FindFirst(ctx, spec)
		if err != nil {
			return err
		}
		if !existingUser.IsZero() {
			return ezutil.ConflictError(fmt.Sprintf(appconstant.ErrAuthDuplicateUser, request.Email))
		}

		hash, err := as.hashService.Hash(request.Password)
		if err != nil {
			return err
		}

		spec.Model.Password = hash

		user, err := as.userRepository.Insert(ctx, spec.Model)
		if err != nil {
			return err
		}

		profile := entity.UserProfile{
			UserID: user.ID,
			Name:   util.GetNameFromEmail(request.Email),
		}

		if _, err = as.userProfileRepository.Insert(ctx, profile); err != nil {
			return err
		}

		return nil
	})
}

func (as *authServiceImpl) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	spec := ezutil.Specification[entity.User]{}
	spec.Model.Email = request.Email

	user, err := as.userRepository.FindFirst(ctx, spec)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if user.IsZero() {
		return dto.LoginResponse{}, ezutil.NotFoundError(appconstant.ErrAuthUnknownCredentials)
	}

	ok, err := as.hashService.CheckHash(user.Password, request.Password)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if !ok {
		return dto.LoginResponse{}, ezutil.NotFoundError(appconstant.ErrAuthUnknownCredentials)
	}

	token, err := as.jwtService.CreateToken(mapper.UserToAuthData(user))
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.LoginResponse{
		Type:  "Bearer",
		Token: token,
	}, nil
}
