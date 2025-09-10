package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/cocoon/internal/test/service_test/mocks"
	"github.com/itsLeonB/sekure"
	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Register_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHashService := mocks.NewMockHashService(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockProfileService := mocks.NewMockProfileService(ctrl)

	authService := service.NewAuthService(
		mockHashService,
		mockJWTService,
		mockUserRepo,
		mockTransactor,
		mockProfileService,
	)

	ctx := context.Background()
	request := dto.RegisterRequest{
		Email:                "test@example.com",
		Password:             "password123",
		PasswordConfirmation: "password123",
	}

	// Mock expectations
	mockTransactor.EXPECT().WithinTransaction(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	mockUserRepo.EXPECT().FindFirst(ctx, gomock.Any()).Return(entity.User{}, nil)
	mockHashService.EXPECT().Hash("password123").Return("hashedpassword", nil)

	newUser := entity.User{Email: "test@example.com", Password: "hashedpassword"}
	newUser.ID = uuid.New()
	mockUserRepo.EXPECT().Insert(ctx, gomock.Any()).Return(newUser, nil)

	profileResp := dto.ProfileResponse{ID: uuid.New()}
	mockProfileService.EXPECT().Create(ctx, gomock.Any()).Return(profileResp, nil)

	err := authService.Register(ctx, request)

	assert.NoError(t, err)
}

func TestAuthService_Login_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHashService := mocks.NewMockHashService(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockProfileService := mocks.NewMockProfileService(ctrl)

	authService := service.NewAuthService(
		mockHashService,
		mockJWTService,
		mockUserRepo,
		mockTransactor,
		mockProfileService,
	)

	ctx := context.Background()
	request := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	user := entity.User{Email: "test@example.com", Password: "hashedpassword"}
	user.ID = uuid.New()

	mockUserRepo.EXPECT().FindFirst(ctx, gomock.Any()).Return(user, nil)
	mockHashService.EXPECT().CheckHash("hashedpassword", "password123").Return(true, nil)
	mockJWTService.EXPECT().CreateToken(gomock.Any()).Return("jwt.token.here", nil)

	response, err := authService.Login(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, "Bearer", response.Type)
	assert.Equal(t, "jwt.token.here", response.Token)
}

func TestAuthService_VerifyToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHashService := mocks.NewMockHashService(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockProfileService := mocks.NewMockProfileService(ctrl)

	authService := service.NewAuthService(
		mockHashService,
		mockJWTService,
		mockUserRepo,
		mockTransactor,
		mockProfileService,
	)

	ctx := context.Background()
	token := "valid.jwt.token"
	userID := uuid.New()
	profileID := uuid.New()

	claims := sekure.JWTClaims{
		Data: map[string]interface{}{
			appconstant.ContextUserID: userID.String(),
		},
	}

	user := entity.User{Email: "test@example.com"}
	user.ID = userID
	user.Profile = entity.UserProfile{UserID: userID}
	user.Profile.ID = profileID

	mockJWTService.EXPECT().VerifyToken(token).Return(claims, nil)
	mockUserRepo.EXPECT().FindFirst(ctx, gomock.Any()).Return(user, nil)

	authData, err := authService.VerifyToken(ctx, token)

	assert.NoError(t, err)
	assert.Equal(t, profileID, authData.ProfileID)
}

func TestAuthService_VerifyToken_InvalidToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHashService := mocks.NewMockHashService(ctrl)
	mockJWTService := mocks.NewMockJWTService(ctrl)
	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockProfileService := mocks.NewMockProfileService(ctrl)

	authService := service.NewAuthService(
		mockHashService,
		mockJWTService,
		mockUserRepo,
		mockTransactor,
		mockProfileService,
	)

	ctx := context.Background()
	token := "invalid.jwt.token"

	mockJWTService.EXPECT().VerifyToken(token).Return(sekure.JWTClaims{}, eris.New("invalid token"))

	_, err := authService.VerifyToken(ctx, token)

	assert.Error(t, err)
}
