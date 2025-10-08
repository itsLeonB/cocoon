package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mapper"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
)

type userServiceImpl struct {
	transactor             crud.Transactor
	userRepo               crud.Repository[entity.User]
	profileSvc             ProfileService
	passwordResetTokenRepo crud.Repository[entity.PasswordResetToken]
}

func NewUserService(
	transactor crud.Transactor,
	userRepo crud.Repository[entity.User],
	profileSvc ProfileService,
	passwordResetTokenRepo crud.Repository[entity.PasswordResetToken],
) UserService {
	return &userServiceImpl{
		transactor,
		userRepo,
		profileSvc,
		passwordResetTokenRepo,
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
	userSpec.PreloadRelations = []string{"Profile"}
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

func (us *userServiceImpl) GeneratePasswordResetToken(ctx context.Context, userID uuid.UUID) (string, error) {
	token, err := us.generateRandomToken(255)
	if err != nil {
		return "", err
	}
	resetToken := entity.PasswordResetToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	if _, err := us.passwordResetTokenRepo.Insert(ctx, resetToken); err != nil {
		return "", err
	}
	return token, nil
}

func (us *userServiceImpl) ResetPassword(ctx context.Context, userID uuid.UUID, email, resetToken, password string) (entity.User, error) {
	spec := crud.Specification[entity.User]{}
	spec.Model.ID = userID
	spec.Model.Email = email
	spec.PreloadRelations = []string{"PasswordResetTokens"}
	user, err := us.getBySpec(ctx, spec)
	if err != nil {
		return entity.User{}, err
	}

	if !us.validateToken(user.PasswordResetTokens, resetToken) {
		return entity.User{}, ungerr.BadRequestError("invalid or expired reset token")
	}

	user.Password = password
	updatedUser, err := us.userRepo.Update(ctx, user)
	if err != nil {
		return entity.User{}, err
	}

	if err = us.passwordResetTokenRepo.DeleteMany(ctx, user.PasswordResetTokens); err != nil {
		return entity.User{}, err
	}

	return updatedUser, nil
}

func (us *userServiceImpl) validateToken(resetTokens []entity.PasswordResetToken, resetToken string) bool {
	if len(resetTokens) < 1 {
		return false
	}
	if len(resetTokens) == 1 {
		return resetTokens[0].IsValid() && resetTokens[0].Token == resetToken
	}
	sort.Slice(resetTokens, func(i, j int) bool {
		return resetTokens[i].CreatedAt.After(resetTokens[j].CreatedAt)
	})
	return resetTokens[0].IsValid() && resetTokens[0].Token == resetToken
}

func (us *userServiceImpl) generateRandomToken(length int) (string, error) {
	tokenBytes := make([]byte, length)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", eris.Wrap(err, "error generating random token")
	}
	return base64.URLEncoding.EncodeToString(tokenBytes)[:length], nil
}

func (us *userServiceImpl) getByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	spec := crud.Specification[entity.User]{}
	spec.Model.ID = id
	spec.PreloadRelations = []string{"Profile"}
	return us.getBySpec(ctx, spec)
}

func (us *userServiceImpl) getBySpec(ctx context.Context, spec crud.Specification[entity.User]) (entity.User, error) {
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
