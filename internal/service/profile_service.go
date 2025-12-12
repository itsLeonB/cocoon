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

type profileServiceImpl struct {
	transactor         crud.Transactor
	profileRepo        repository.UserProfileRepository
	userRepo           crud.Repository[entity.User]
	friendshipRepo     repository.FriendshipRepository
	relatedProfileRepo crud.Repository[entity.RelatedProfile]
}

func NewProfileService(
	transactor crud.Transactor,
	profileRepo repository.UserProfileRepository,
	userRepo crud.Repository[entity.User],
	friendshipRepo repository.FriendshipRepository,
	relatedProfileRepo crud.Repository[entity.RelatedProfile],
) ProfileService {
	return &profileServiceImpl{
		transactor,
		profileRepo,
		userRepo,
		friendshipRepo,
		relatedProfileRepo,
	}
}

func (ps *profileServiceImpl) Create(ctx context.Context, request dto.NewProfileRequest) (dto.ProfileResponse, error) {
	var response dto.ProfileResponse

	err := ps.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		newProfile := entity.UserProfile{
			UserID: uuid.NullUUID{
				UUID:  request.UserID,
				Valid: request.UserID != uuid.Nil,
			},
			Name:   request.Name,
			Avatar: request.Avatar,
		}

		insertedProfile, err := ps.profileRepo.Insert(ctx, newProfile)
		if err != nil {
			return err
		}

		response = mapper.ProfileToResponse(insertedProfile, "", nil, uuid.Nil)

		return nil
	})

	return response, err
}

func (ps *profileServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error) {
	profile, err := ps.getByID(ctx, id)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	var email string
	if profile.IsReal() {
		userSpec := crud.Specification[entity.User]{}
		userSpec.Model.ID = profile.UserID.UUID
		user, err := ps.userRepo.FindFirst(ctx, userSpec)
		if err != nil {
			return dto.ProfileResponse{}, err
		}
		email = user.Email
	}

	anonProfileIDs, realProfileID, err := ps.getAssociations(ctx, profile)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	return mapper.ProfileToResponse(profile, email, anonProfileIDs, realProfileID), nil
}

func (ps *profileServiceImpl) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]dto.ProfileResponse, error) {
	profiles, err := ps.profileRepo.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(profiles, mapper.SimpleProfileToResponse), nil
}

func (ps *profileServiceImpl) Update(ctx context.Context, req dto.UpdateProfileRequest) (dto.ProfileResponse, error) {
	var response dto.ProfileResponse
	err := ps.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.UserProfile]{}
		spec.Model.ID = req.ID
		spec.Model.UserID = uuid.NullUUID{
			UUID:  req.UserID,
			Valid: req.UserID != uuid.Nil,
		}
		spec.DeletedFilter = crud.ExcludeDeleted
		spec.ForUpdate = true
		profile, err := ps.profileRepo.FindFirst(ctx, spec)
		if err != nil {
			return err
		}
		if profile.IsZero() {
			return ungerr.NotFoundError(fmt.Sprintf("profile ID: %s is not found", req.ID.String()))
		}

		if req.UserID != uuid.Nil {
			profile.UserID = uuid.NullUUID{
				UUID:  req.UserID,
				Valid: true,
			}
		}
		if req.Name != "" {
			profile.Name = req.Name
		}
		if req.Avatar != "" {
			profile.Avatar = req.Avatar
		}
		updatedProfile, err := ps.profileRepo.Update(ctx, profile)
		if err != nil {
			return err
		}

		response = mapper.ProfileToResponse(updatedProfile, "", nil, uuid.Nil)
		return nil
	})
	return response, err
}

func (ps *profileServiceImpl) GetByEmail(ctx context.Context, email string) (dto.ProfileResponse, error) {
	userSpec := crud.Specification[entity.User]{}
	userSpec.Model.Email = email
	userSpec.DeletedFilter = crud.ExcludeDeleted
	userSpec.PreloadRelations = []string{"Profile"}
	user, err := ps.userRepo.FindFirst(ctx, userSpec)
	if err != nil {
		return dto.ProfileResponse{}, err
	}
	if user.IsZero() || user.IsDeleted() || !user.IsVerified() {
		return dto.ProfileResponse{}, ungerr.NotFoundError(appconstant.ErrUserNotFound)
	}

	anonProfileIDs, realProfileID, err := ps.getAssociations(ctx, user.Profile)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	return mapper.ProfileToResponse(user.Profile, user.Email, anonProfileIDs, realProfileID), nil
}

func (ps *profileServiceImpl) SearchByName(ctx context.Context, query string, limit int) ([]dto.ProfileResponse, error) {
	profileNames, err := ps.profileRepo.SearchByName(ctx, query, limit)
	if err != nil {
		return nil, err
	}
	if len(profileNames) < 1 {
		return []dto.ProfileResponse{}, nil
	}

	ids := ezutil.MapSlice(profileNames, func(pn entity.ProfileName) uuid.UUID { return pn.ID })

	profiles, err := ps.GetByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	profileByID := make(map[uuid.UUID]dto.ProfileResponse, len(profiles))
	for _, profile := range profiles {
		profileByID[profile.ID] = profile
	}

	ordered := make([]dto.ProfileResponse, 0, len(profileNames))
	for _, pn := range profileNames {
		if profile, ok := profileByID[pn.ID]; ok {
			ordered = append(ordered, profile)
		}
	}

	return ordered, nil
}

func (ps *profileServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	return ps.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		profile, err := ps.getByID(ctx, id)
		if err != nil {
			return err
		}

		return ps.profileRepo.Delete(ctx, profile)
	})
}

func (ps *profileServiceImpl) getByID(ctx context.Context, id uuid.UUID) (entity.UserProfile, error) {
	spec := crud.Specification[entity.UserProfile]{}
	spec.Model.ID = id
	profile, err := ps.profileRepo.FindFirst(ctx, spec)
	if err != nil {
		return entity.UserProfile{}, err
	}
	if profile.IsZero() {
		return entity.UserProfile{}, ungerr.NotFoundError(fmt.Sprintf("profile with ID: %s is not found", id))
	}
	if profile.IsDeleted() {
		return entity.UserProfile{}, ungerr.ConflictError(fmt.Sprintf("profile with ID: %s is already deleted", id))
	}
	return profile, nil
}

func (ps *profileServiceImpl) Associate(ctx context.Context, request dto.AssociateProfileRequest) error {
	return ps.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		if request.RealProfileID == uuid.Nil || request.AnonProfileID == uuid.Nil || request.UserProfileID == uuid.Nil {
			return ungerr.BadRequestError("userProfileID / realProfileID / anonProfileID cannot be nil")
		}

		if _, err := ps.getByID(ctx, request.RealProfileID); err != nil {
			return err
		}
		if _, err := ps.getByID(ctx, request.AnonProfileID); err != nil {
			return err
		}

		if err := ps.validateAssociation(ctx, request); err != nil {
			return err
		}

		return ps.createAssociation(ctx, request)
	})
}

func (ps *profileServiceImpl) validateAssociation(ctx context.Context, request dto.AssociateProfileRequest) error {
	relatedSpec := crud.Specification[entity.RelatedProfile]{}
	relatedSpec.Model.AnonProfileID = request.AnonProfileID
	existingRelated, err := ps.relatedProfileRepo.FindFirst(ctx, relatedSpec)
	if err != nil {
		return err
	}
	if !existingRelated.IsZero() {
		return ungerr.ConflictError("anonProfileID is already associated with a real profile")
	}

	if err := ps.checkFriendship(ctx, request.UserProfileID, request.RealProfileID, "real"); err != nil {
		return err
	}
	if err := ps.checkFriendship(ctx, request.UserProfileID, request.AnonProfileID, "anonymous"); err != nil {
		return err
	}
	return nil
}

func (ps *profileServiceImpl) checkFriendship(ctx context.Context, userProfileID, friendProfileID uuid.UUID, typeStr string) error {
	f, err := ps.friendshipRepo.FindByProfileIDs(ctx, userProfileID, friendProfileID)
	if err != nil {
		return err
	}
	if f.IsZero() || f.IsDeleted() {
		return ungerr.ForbiddenError(fmt.Sprintf("user is not friends with the %s profile", typeStr))
	}
	return nil
}

func (ps *profileServiceImpl) createAssociation(ctx context.Context, request dto.AssociateProfileRequest) error {
	newRelated := entity.RelatedProfile{
		RealProfileID: request.RealProfileID,
		AnonProfileID: request.AnonProfileID,
	}
	_, err := ps.relatedProfileRepo.Insert(ctx, newRelated)
	return err
}

func (ps *profileServiceImpl) getAssociatedProfileIDs(ctx context.Context, realProfileID uuid.UUID) ([]uuid.UUID, error) {
	spec := crud.Specification[entity.RelatedProfile]{}
	spec.Model.RealProfileID = realProfileID
	relations, err := ps.relatedProfileRepo.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}
	return ezutil.MapSlice(relations, func(r entity.RelatedProfile) uuid.UUID { return r.AnonProfileID }), nil
}

func (ps *profileServiceImpl) GetRealProfileID(ctx context.Context, anonProfileID uuid.UUID) (uuid.UUID, error) {
	spec := crud.Specification[entity.RelatedProfile]{}
	spec.Model.AnonProfileID = anonProfileID
	relation, err := ps.relatedProfileRepo.FindFirst(ctx, spec)
	return relation.RealProfileID, err
}

func (ps *profileServiceImpl) getAssociations(ctx context.Context, profile entity.UserProfile) ([]uuid.UUID, uuid.UUID, error) {
	if profile.IsReal() {
		anonProfileIDs, err := ps.getAssociatedProfileIDs(ctx, profile.ID)
		if err != nil {
			return nil, uuid.Nil, err
		}
		return anonProfileIDs, uuid.Nil, nil
	} else {
		profileID, err := ps.GetRealProfileID(ctx, profile.ID)
		if err != nil {
			return nil, uuid.Nil, err
		}
		return nil, profileID, nil
	}
}
