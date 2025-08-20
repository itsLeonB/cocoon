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

type friendshipServiceImpl struct {
	transactor            ezutil.Transactor
	userProfileRepository repository.UserProfileRepository
	friendshipRepository  repository.FriendshipRepository
	userService           UserService
}

func NewFriendshipService(
	transactor ezutil.Transactor,
	userProfileRepository repository.UserProfileRepository,
	friendshipRepository repository.FriendshipRepository,
	userService UserService,
) FriendshipService {
	return &friendshipServiceImpl{
		transactor,
		userProfileRepository,
		friendshipRepository,
		userService,
	}
}

func (fs *friendshipServiceImpl) CreateAnonymous(
	ctx context.Context,
	request dto.NewAnonymousFriendshipRequest,
) (dto.FriendshipResponse, error) {
	var response dto.FriendshipResponse

	err := fs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		user, err := fs.userService.GetEntityByID(ctx, request.UserID)
		if err != nil {
			return err
		}

		if err = fs.validateExistingAnonymousFriendship(ctx, user.Profile.ID, request.Name); err != nil {
			return err
		}

		response, err = fs.insertAnonymousFriendship(ctx, user.Profile, request.Name)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	return response, nil
}

func (fs *friendshipServiceImpl) GetAll(ctx context.Context, userID uuid.UUID) ([]dto.FriendshipResponse, error) {
	user, err := fs.userService.GetEntityByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	spec := entity.FriendshipSpecification{}
	spec.Model.ProfileID1 = user.Profile.ID
	spec.PreloadRelations = []string{"Profile1", "Profile2"}

	friendships, err := fs.friendshipRepository.FindAllBySpec(ctx, spec)
	if err != nil {
		return nil, err
	}

	mapperFunc := func(friendship entity.Friendship) (dto.FriendshipResponse, error) {
		return mapper.FriendshipToResponse(user.Profile.ID, friendship)
	}

	return ezutil.MapSliceWithError(friendships, mapperFunc)
}

func (fs *friendshipServiceImpl) GetDetails(ctx context.Context, userID, friendshipID uuid.UUID) (dto.FriendDetails, error) {
	user, err := fs.userService.GetEntityByID(ctx, userID)
	if err != nil {
		return dto.FriendDetails{}, err
	}

	spec := entity.FriendshipSpecification{}
	spec.Model.ID = friendshipID
	spec.PreloadRelations = []string{"Profile1", "Profile2"}
	friendship, err := fs.friendshipRepository.FindFirstBySpec(ctx, spec)
	if err != nil {
		return dto.FriendDetails{}, err
	}
	if friendship.IsZero() {
		return dto.FriendDetails{}, ezutil.NotFoundError(appconstant.ErrFriendshipNotFound)
	}

	return mapper.MapToFriendDetails(user.Profile.ID, friendship)
}

func (fs *friendshipServiceImpl) IsFriends(ctx context.Context, profileID1, profileID2 uuid.UUID) (bool, error) {
	friendship, err := fs.friendshipRepository.FindByProfileIDs(ctx, profileID1, profileID2)
	if err != nil {
		return false, err
	}

	if friendship.IsZero() || friendship.IsDeleted() {
		return false, nil
	}

	return true, nil
}

func (fs *friendshipServiceImpl) validateExistingAnonymousFriendship(
	ctx context.Context,
	userProfileID uuid.UUID,
	friendName string,
) error {
	friendshipSpec := entity.FriendshipSpecification{}
	friendshipSpec.Model.ProfileID1 = userProfileID
	friendshipSpec.Name = friendName
	friendshipSpec.Model.Type = appconstant.Anonymous

	existingFriendship, err := fs.friendshipRepository.FindFirstBySpec(ctx, friendshipSpec)
	if err != nil {
		return err
	}
	if !existingFriendship.IsZero() && !existingFriendship.IsDeleted() {
		return ezutil.ConflictError(fmt.Sprintf("anonymous friend named %s already exists", friendName))
	}

	return nil
}

func (fs *friendshipServiceImpl) insertAnonymousFriendship(
	ctx context.Context,
	userProfile entity.UserProfile,
	friendName string,
) (dto.FriendshipResponse, error) {
	newProfile := entity.UserProfile{Name: friendName}

	insertedProfile, err := fs.userProfileRepository.Insert(ctx, newProfile)
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	newFriendship, err := mapper.OrderProfilesToFriendship(userProfile, insertedProfile)
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	newFriendship.Type = appconstant.Anonymous

	insertedFriendship, err := fs.friendshipRepository.Insert(ctx, newFriendship)
	if err != nil {
		return dto.FriendshipResponse{}, err
	}

	return mapper.FriendshipToResponse(userProfile.ID, insertedFriendship)
}
