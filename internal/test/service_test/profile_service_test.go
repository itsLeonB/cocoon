package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mocks"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/cocoon/internal/util"
	"github.com/itsLeonB/go-crud"
	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestProfileService_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := crud.NewMockTransactor(ctrl)
	mockProfileRepo := mocks.NewMockUserProfileRepository(ctrl)
	mockUserRepo := crud.NewMockRepository[entity.User](ctrl)
	mockFriendshipRepo := mocks.NewMockFriendshipRepository(ctrl)
	mockRelatedRepo := crud.NewMockRepository[entity.RelatedProfile](ctrl)

	profileService := service.NewProfileService(
		mockTransactor,
		mockProfileRepo,
		mockUserRepo,
		mockFriendshipRepo,
		mockRelatedRepo,
	)

	ctx := context.Background()
	userID := uuid.New()
	request := dto.NewProfileRequest{
		UserID: userID,
		Name:   "John Doe",
	}

	insertedProfile := entity.UserProfile{
		UserID: util.NewValidNullUUID(userID),
		Name:   "John Doe",
	}
	insertedProfile.ID = uuid.New()

	mockTransactor.EXPECT().WithinTransaction(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	mockProfileRepo.EXPECT().Insert(ctx, gomock.Any()).Return(insertedProfile, nil)

	response, err := profileService.Create(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, insertedProfile.ID, response.ID)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, "John Doe", response.Name)
}

func TestProfileService_GetByID_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := crud.NewMockTransactor(ctrl)
	mockProfileRepo := mocks.NewMockUserProfileRepository(ctrl)
	mockUserRepo := crud.NewMockRepository[entity.User](ctrl)
	mockFriendshipRepo := mocks.NewMockFriendshipRepository(ctrl)
	mockRelatedRepo := crud.NewMockRepository[entity.RelatedProfile](ctrl)

	profileService := service.NewProfileService(
		mockTransactor,
		mockProfileRepo,
		mockUserRepo,
		mockFriendshipRepo,
		mockRelatedRepo,
	)

	ctx := context.Background()
	profileID := uuid.New()
	userID := uuid.New()

	profile := entity.UserProfile{
		UserID: util.NewValidNullUUID(userID),
		Name:   "John Doe",
	}
	profile.ID = profileID

	mockProfileRepo.EXPECT().FindFirst(ctx, gomock.Any()).Return(profile, nil)
	mockUserRepo.EXPECT().FindFirst(ctx, gomock.Any()).Return(entity.User{}, nil)
	mockRelatedRepo.EXPECT().FindAll(ctx, gomock.Any()).Return([]entity.RelatedProfile{}, nil)

	response, err := profileService.GetByID(ctx, profileID)

	assert.NoError(t, err)
	assert.Equal(t, profileID, response.ID)
	assert.Equal(t, userID, response.UserID)
	assert.Equal(t, "John Doe", response.Name)
}

func TestProfileService_GetByIDs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := crud.NewMockTransactor(ctrl)
	mockProfileRepo := mocks.NewMockUserProfileRepository(ctrl)
	mockUserRepo := crud.NewMockRepository[entity.User](ctrl)
	mockFriendshipRepo := mocks.NewMockFriendshipRepository(ctrl)
	mockRelatedRepo := crud.NewMockRepository[entity.RelatedProfile](ctrl)

	profileService := service.NewProfileService(
		mockTransactor,
		mockProfileRepo,
		mockUserRepo,
		mockFriendshipRepo,
		mockRelatedRepo,
	)

	ctx := context.Background()
	profileID1 := uuid.New()
	profileID2 := uuid.New()
	ids := []uuid.UUID{profileID1, profileID2}

	profiles := []entity.UserProfile{
		{UserID: util.NewValidNullUUID(uuid.New()), Name: "User 1"},
		{UserID: util.NewValidNullUUID(uuid.New()), Name: "User 2"},
	}
	profiles[0].ID = profileID1
	profiles[1].ID = profileID2

	mockProfileRepo.EXPECT().FindByIDs(ctx, ids).Return(profiles, nil)

	responses, err := profileService.GetByIDs(ctx, ids)

	assert.NoError(t, err)
	assert.Len(t, responses, 2)
	assert.Equal(t, profileID1, responses[0].ID)
	assert.Equal(t, profileID2, responses[1].ID)
}

func TestProfileService_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := crud.NewMockTransactor(ctrl)
	mockProfileRepo := mocks.NewMockUserProfileRepository(ctrl)
	mockUserRepo := crud.NewMockRepository[entity.User](ctrl)
	mockFriendshipRepo := mocks.NewMockFriendshipRepository(ctrl)
	mockRelatedRepo := crud.NewMockRepository[entity.RelatedProfile](ctrl)

	profileService := service.NewProfileService(
		mockTransactor,
		mockProfileRepo,
		mockUserRepo,
		mockFriendshipRepo,
		mockRelatedRepo,
	)

	ctx := context.Background()
	request := dto.NewProfileRequest{
		UserID: uuid.New(),
		Name:   "John Doe",
	}

	mockTransactor.EXPECT().WithinTransaction(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	mockProfileRepo.EXPECT().Insert(ctx, gomock.Any()).Return(entity.UserProfile{}, eris.New("database error"))

	_, err := profileService.Create(ctx, request)

	assert.Error(t, err)
}

func TestProfileService_Associate_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := crud.NewMockTransactor(ctrl)
	mockProfileRepo := mocks.NewMockUserProfileRepository(ctrl)
	mockFriendshipRepo := mocks.NewMockFriendshipRepository(ctrl)
	mockRelatedRepo := crud.NewMockRepository[entity.RelatedProfile](ctrl)
	mockUserRepo := crud.NewMockRepository[entity.User](ctrl)

	profileService := service.NewProfileService(
		mockTransactor,
		mockProfileRepo,
		mockUserRepo,
		mockFriendshipRepo,
		mockRelatedRepo,
	)

	ctx := context.Background()
	userProfileID := uuid.New()
	realProfileID := uuid.New()
	anonProfileID := uuid.New()

	request := dto.AssociateProfileRequest{
		UserProfileID: userProfileID,
		RealProfileID: realProfileID,
		AnonProfileID: anonProfileID,
	}

	// Transactor
	mockTransactor.EXPECT().WithinTransaction(ctx, gomock.Any()).DoAndReturn(
		func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		})

	// GetByID calls (first for real, second for anon)
	// Note: getByID calls FindFirst with specific ID
	// Since arguments are gomock.Any(), order matters.
	mockProfileRepo.EXPECT().FindFirst(ctx, gomock.Any()).Return(entity.UserProfile{BaseEntity: crud.BaseEntity{ID: realProfileID}}, nil)
	mockProfileRepo.EXPECT().FindFirst(ctx, gomock.Any()).Return(entity.UserProfile{BaseEntity: crud.BaseEntity{ID: anonProfileID}}, nil)

	// Check if related profile exists (should be empty)
	mockRelatedRepo.EXPECT().FindFirst(ctx, gomock.Any()).Return(entity.RelatedProfile{}, nil)

	// Check friendship with real profile
	mockFriendshipRepo.EXPECT().FindByProfileIDs(ctx, userProfileID, realProfileID).Return(entity.Friendship{BaseEntity: crud.BaseEntity{ID: uuid.New()}}, nil)

	// Check friendship with anon profile
	mockFriendshipRepo.EXPECT().FindByProfileIDs(ctx, userProfileID, anonProfileID).Return(entity.Friendship{BaseEntity: crud.BaseEntity{ID: uuid.New()}}, nil)

	// Insert association
	mockRelatedRepo.EXPECT().Insert(ctx, gomock.Any()).Return(entity.RelatedProfile{}, nil)

	err := profileService.Associate(ctx, request)

	assert.NoError(t, err)
}
