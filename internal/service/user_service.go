package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mapper"
	"github.com/itsLeonB/go-crud"
	"github.com/rotisserie/eris"
)

type userServiceImpl struct {
	transactor crud.Transactor
	userRepo   crud.Repository[entity.User]
	profileSvc ProfileService
}

func NewUserService(
	transactor crud.Transactor,
	userRepo crud.Repository[entity.User],
	profileSvc ProfileService,
) UserService {
	return &userServiceImpl{
		transactor,
		userRepo,
		profileSvc,
	}
}

func (us *userServiceImpl) CreateNew(ctx context.Context, request dto.NewUserRequest) (entity.User, error) {
	var response entity.User
	err := us.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		newUser := entity.User{
			Email:    request.Email,
			Password: request.Password,
		}

		if request.VerifyNow {
			newUser.VerifiedAt = sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			}
		}

		user, err := us.userRepo.Insert(ctx, newUser)
		if err != nil {
			return err
		}

		profile := dto.NewProfileRequest{
			UserID: user.ID,
			Name:   request.Name,
			Avatar: request.Avatar,
		}

		if _, err = us.profileSvc.Create(ctx, profile); err != nil {
			return err
		}

		response = user
		return nil
	})

	return response, err
}

func (us *userServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error) {
	user, err := us.getByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return mapper.UserToResponse(user), nil
}

func (us *userServiceImpl) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	userSpec := crud.Specification[entity.User]{}
	userSpec.Model.Email = email
	userSpec.DeletedFilter = crud.ExcludeDeleted
	return us.userRepo.FindFirst(ctx, userSpec)
}

func (us *userServiceImpl) Verify(ctx context.Context, id uuid.UUID, email string) (entity.User, error) {
	user, err := us.getByID(ctx, id)
	if err != nil {
		return entity.User{}, err
	}
	if user.Email != email {
		return entity.User{}, eris.New("email does not match")
	}

	user.VerifiedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	return us.userRepo.Update(ctx, user)
}

func (us *userServiceImpl) getByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	spec := crud.Specification[entity.User]{}
	spec.Model.ID = id
	spec.PreloadRelations = []string{"Profile"}
	user, err := us.userRepo.FindFirst(ctx, spec)
	if err != nil {
		return entity.User{}, err
	}
	if user.IsZero() {
		return entity.User{}, eris.New("user ID is not found")
	}
	if user.IsDeleted() {
		return entity.User{}, eris.New("user is deleted")
	}
	return user, nil
}
