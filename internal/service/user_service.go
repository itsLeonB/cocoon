package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mapper"
	"github.com/itsLeonB/cocoon/internal/repository"
	"github.com/itsLeonB/ezutil"
)

type userServiceImpl struct {
	transactor            ezutil.Transactor
	userRepository        repository.UserRepository
	userProfileRepository repository.UserProfileRepository
}

func NewUserService(
	transactor ezutil.Transactor,
	userRepository repository.UserRepository,
	userProfileRepository repository.UserProfileRepository,
) UserService {
	return &userServiceImpl{
		transactor,
		userRepository,
		userProfileRepository,
	}
}

func (us *userServiceImpl) GetProfile(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error) {
	spec := ezutil.Specification[entity.User]{}
	spec.Model.ID = id

	user, err := us.userRepository.FindFirst(ctx, spec)
	if err != nil {
		return dto.ProfileResponse{}, err
	}
	if user.IsZero() {
		return dto.ProfileResponse{}, nil
	}

	return mapper.UserToProfileResponse(user), nil
}

func (us *userServiceImpl) GetEntityByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	userSpec := ezutil.Specification[entity.User]{}
	userSpec.Model.ID = id

	user, err := us.userRepository.FindFirst(ctx, userSpec)
	if err != nil {
		return entity.User{}, err
	}
	if user.IsZero() {
		return entity.User{}, ezutil.NotFoundError(fmt.Sprintf(appconstant.ErrUserNotFound, id))
	}
	if user.IsDeleted() {
		return entity.User{}, ezutil.UnprocessableEntityError(fmt.Sprintf(appconstant.ErrUserDeleted, id))
	}

	return user, nil
}
