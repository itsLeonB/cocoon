package service_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/cocoon/internal/test/service_test/mocks"
	"github.com/rotisserie/eris"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFriendshipService_IsFriends_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockFriendshipRepo := mocks.NewMockFriendshipRepository(ctrl)
	mockProfileService := mocks.NewMockProfileService(ctrl)

	friendshipService := service.NewFriendshipService(
		mockTransactor,
		mockFriendshipRepo,
		mockProfileService,
	)

	ctx := context.Background()
	profileID1 := uuid.New()
	profileID2 := uuid.New()

	friendship := entity.Friendship{
		ProfileID1: profileID1,
		ProfileID2: profileID2,
	}
	friendship.ID = uuid.New()

	mockFriendshipRepo.EXPECT().FindByProfileIDs(ctx, profileID1, profileID2).Return(friendship, nil)

	isFriends, isDeleted, err := friendshipService.IsFriends(ctx, profileID1, profileID2)

	assert.NoError(t, err)
	assert.True(t, isFriends)
	assert.False(t, isDeleted)
}

func TestFriendshipService_IsFriends_NotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockFriendshipRepo := mocks.NewMockFriendshipRepository(ctrl)
	mockProfileService := mocks.NewMockProfileService(ctrl)

	friendshipService := service.NewFriendshipService(
		mockTransactor,
		mockFriendshipRepo,
		mockProfileService,
	)

	ctx := context.Background()
	profileID1 := uuid.New()
	profileID2 := uuid.New()

	mockFriendshipRepo.EXPECT().FindByProfileIDs(ctx, profileID1, profileID2).Return(entity.Friendship{}, nil)

	isFriends, isDeleted, err := friendshipService.IsFriends(ctx, profileID1, profileID2)

	assert.NoError(t, err)
	assert.False(t, isFriends)
	assert.False(t, isDeleted)
}

func TestFriendshipService_IsFriends_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTransactor := mocks.NewMockTransactor(ctrl)
	mockFriendshipRepo := mocks.NewMockFriendshipRepository(ctrl)
	mockProfileService := mocks.NewMockProfileService(ctrl)

	friendshipService := service.NewFriendshipService(
		mockTransactor,
		mockFriendshipRepo,
		mockProfileService,
	)

	ctx := context.Background()
	profileID1 := uuid.New()
	profileID2 := uuid.New()

	mockFriendshipRepo.EXPECT().FindByProfileIDs(ctx, profileID1, profileID2).Return(entity.Friendship{}, eris.New("database error"))

	_, _, err := friendshipService.IsFriends(ctx, profileID1, profileID2)

	assert.Error(t, err)
}
