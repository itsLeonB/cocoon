package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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
}

func NewProfileService(
	transactor crud.Transactor,
	profileRepo repository.UserProfileRepository,
) ProfileService {
	return &profileServiceImpl{
		transactor,
		profileRepo,
	}
}

func (ps *profileServiceImpl) Create(ctx context.Context, request dto.NewProfileRequest) (dto.ProfileResponse, error) {
	var response dto.ProfileResponse

	err := ps.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		newProfile := entity.UserProfile{
			UserID: request.UserID,
			Name:   request.Name,
			Avatar: request.Avatar,
		}

		insertedProfile, err := ps.profileRepo.Insert(ctx, newProfile)
		if err != nil {
			return err
		}

		response = mapper.ProfileToResponse(insertedProfile)

		return nil
	})

	return response, err
}

func (ps *profileServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error) {
	spec := crud.Specification[entity.UserProfile]{}
	spec.Model.ID = id
	profile, err := ps.profileRepo.FindFirst(ctx, spec)
	if err != nil {
		return dto.ProfileResponse{}, err
	}
	if profile.IsZero() {
		return dto.ProfileResponse{}, ungerr.NotFoundError(fmt.Sprintf("profile with ID: %s is not found", id))
	}
	if profile.IsDeleted() {
		return dto.ProfileResponse{}, ungerr.ConflictError(fmt.Sprintf("profile with ID: %s is already deleted", id))
	}

	return mapper.ProfileToResponse(profile), nil
}

func (ps *profileServiceImpl) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]dto.ProfileResponse, error) {
	profiles, err := ps.profileRepo.FindByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(profiles, mapper.ProfileToResponse), nil
}

func (ps *profileServiceImpl) Update(ctx context.Context, req dto.UpdateProfileRequest) (dto.ProfileResponse, error) {
	var response dto.ProfileResponse
	err := ps.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.UserProfile]{}
		spec.Model.ID = req.ID
		if req.UserID != uuid.Nil {
			spec.Model.UserID = req.UserID
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
			profile.UserID = req.UserID
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

		response = mapper.ProfileToResponse(updatedProfile)
		return nil
	})
	return response, err
}
