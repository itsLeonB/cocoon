package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mapper"
	"github.com/itsLeonB/cocoon/internal/repository"
	"github.com/itsLeonB/ezutil"
)

type profileServiceImpl struct {
	transactor  ezutil.Transactor
	profileRepo repository.UserProfileRepository
}

func NewProfileService(
	transactor ezutil.Transactor,
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
	spec := ezutil.Specification[entity.UserProfile]{}
	spec.Model.ID = id
	profile, err := ps.profileRepo.FindFirst(ctx, spec)
	if err != nil {
		return dto.ProfileResponse{}, err
	}
	if profile.IsZero() {
		return dto.ProfileResponse{}, ezutil.NotFoundError(fmt.Sprintf("profile with ID: %s is not found", id))
	}
	if profile.IsDeleted() {
		return dto.ProfileResponse{}, ezutil.ConflictError(fmt.Sprintf("profile with ID: %s is already deleted", id))
	}

	return mapper.ProfileToResponse(profile), nil
}
