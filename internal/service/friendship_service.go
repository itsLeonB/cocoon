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
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
)

type friendshipServiceImpl struct {
	transactor           crud.Transactor
	friendshipRepository repository.FriendshipRepository
	profileService       ProfileService
}

func NewFriendshipService(
	transactor crud.Transactor,
	friendshipRepository repository.FriendshipRepository,
	profileService ProfileService,
) FriendshipService {
	return &friendshipServiceImpl{
		transactor,
		friendshipRepository,
		profileService,
	}
}

func (fs *friendshipServiceImpl) CreateAnonymous(
	ctx context.Context,
	request dto.NewAnonymousFriendshipRequest,
) (dto.FriendshipResponse, error) {
	var response dto.FriendshipResponse

	err := fs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		profile, err := fs.profileService.GetByID(ctx, request.ProfileID)
		if err != nil {
			return err
		}

		if err = fs.validateExistingAnonymousFriendship(ctx, profile.ID, request.Name); err != nil {
			return err
		}

		response, err = fs.insertAnonymousFriendship(ctx, profile, request.Name)
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

func (fs *friendshipServiceImpl) GetAll(ctx context.Context, profileID uuid.UUID) ([]dto.FriendshipResponse, error) {
	profile, err := fs.profileService.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	spec := entity.FriendshipSpecification{}
	spec.Model.ProfileID1 = profile.ID
	spec.PreloadRelations = []string{"Profile1", "Profile2"}

	friendships, err := fs.friendshipRepository.FindAllBySpec(ctx, spec)
	if err != nil {
		return nil, err
	}

	mapperFunc := func(friendship entity.Friendship) (dto.FriendshipResponse, error) {
		return mapper.FriendshipToResponse(profile.ID, friendship)
	}

	return ezutil.MapSliceWithError(friendships, mapperFunc)
}

func (fs *friendshipServiceImpl) GetDetails(ctx context.Context, profileID, friendshipID uuid.UUID) (dto.FriendDetails, error) {
	profile, err := fs.profileService.GetByID(ctx, profileID)
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
		return dto.FriendDetails{}, ungerr.NotFoundError(appconstant.ErrFriendshipNotFound)
	}

	return mapper.MapToFriendDetails(profile.ID, friendship)
}

func (fs *friendshipServiceImpl) IsFriends(ctx context.Context, profileID1, profileID2 uuid.UUID) (bool, bool, error) {
	friendship, err := fs.friendshipRepository.FindByProfileIDs(ctx, profileID1, profileID2)
	if err != nil {
		return false, false, err
	}

	if friendship.IsZero() || friendship.IsDeleted() {
		return false, false, nil
	}

	return true, friendship.Type == appconstant.Anonymous, nil
}

func (fs *friendshipServiceImpl) CreateReal(ctx context.Context, userProfileID, friendProfileID uuid.UUID) (dto.FriendshipResponse, error) {
	var response dto.FriendshipResponse
	err := fs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		profiles, err := fs.profileService.GetByIDs(ctx, []uuid.UUID{userProfileID, friendProfileID})
		if err != nil {
			return err
		}

		var userProfile dto.ProfileResponse
		var friendProfile dto.ProfileResponse
		for _, profile := range profiles {
			if profile.ID == userProfileID {
				userProfile = profile
			}
			if profile.ID == friendProfileID {
				friendProfile = profile
			}
		}

		newFriendship, err := mapper.OrderProfilesToFriendship(userProfile, friendProfile)
		if err != nil {
			return err
		}

		newFriendship.Type = appconstant.Real

		insertedFriendship, err := fs.friendshipRepository.Insert(ctx, newFriendship)
		if err != nil {
			return err
		}

		response, err = mapper.FriendshipToResponse(userProfile.ID, insertedFriendship)
		return err
	})
	return response, err
}

func (fs *friendshipServiceImpl) RemoveAnonymous(ctx context.Context, userProfileID, friendProfileID uuid.UUID) error {
	return fs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		friendProfile, err := fs.profileService.GetByID(ctx, friendProfileID)
		if err != nil {
			return err
		}
		if friendProfile.UserID != uuid.Nil {
			return ungerr.UnprocessableEntityError("cannot remove friend, is not anonymous")
		}

		friendship, err := fs.friendshipRepository.FindByProfileIDs(ctx, userProfileID, friendProfileID)
		if err != nil {
			return err
		}
		if friendship.IsZero() || friendship.IsDeleted() {
			return nil
		}
		if friendship.Type != appconstant.Anonymous {
			return ungerr.UnprocessableEntityError("cannot remove friend, is not anonymous")
		}

		if err = fs.friendshipRepository.Delete(ctx, friendship); err != nil {
			return err
		}

		return fs.profileService.Delete(ctx, friendProfileID)
	})
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
		return ungerr.ConflictError(fmt.Sprintf("anonymous friend named %s already exists", friendName))
	}

	return nil
}

func (fs *friendshipServiceImpl) insertAnonymousFriendship(
	ctx context.Context,
	userProfile dto.ProfileResponse,
	friendName string,
) (dto.FriendshipResponse, error) {
	newProfile := dto.NewProfileRequest{Name: friendName}

	insertedProfile, err := fs.profileService.Create(ctx, newProfile)
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
