package mapper

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
)

func ProfileToResponse(profile entity.UserProfile, email string, anonProfileIDs []uuid.UUID, realProfileID uuid.UUID) dto.ProfileResponse {
	return dto.ProfileResponse{
		UserID:                   profile.UserID.UUID,
		ID:                       profile.ID,
		Name:                     profile.Name,
		Avatar:                   profile.Avatar,
		Email:                    email,
		CreatedAt:                profile.CreatedAt,
		UpdatedAt:                profile.UpdatedAt,
		DeletedAt:                profile.DeletedAt.Time,
		AssociatedAnonProfileIDs: anonProfileIDs,
		RealProfileID:            realProfileID,
	}
}

func SimpleProfileToResponse(profile entity.UserProfile) dto.ProfileResponse {
	return ProfileToResponse(profile, "", nil, uuid.Nil)
}
