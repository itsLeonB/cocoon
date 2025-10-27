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
	transactor  crud.Transactor
	profileRepo repository.UserProfileRepository
	userRepo    crud.Repository[entity.User]
}

func NewProfileService(
	transactor crud.Transactor,
	profileRepo repository.UserProfileRepository,
	userRepo crud.Repository[entity.User],
) ProfileService {
	return &profileServiceImpl{
		transactor,
		profileRepo,
		userRepo,
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

		response = mapper.ProfileToResponse(insertedProfile, "")

		return nil
	})

	return response, err
}

func (ps *profileServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error) {
	profile, err := ps.getByID(ctx, id)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	userSpec := crud.Specification[entity.User]{}
	userSpec.Model.ID = profile.UserID.UUID
	user, err := ps.userRepo.FindFirst(ctx, userSpec)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	return mapper.ProfileToResponse(profile, user.Email), nil
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

		response = mapper.ProfileToResponse(updatedProfile, "")
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
	return mapper.ProfileToResponse(user.Profile, user.Email), nil
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
