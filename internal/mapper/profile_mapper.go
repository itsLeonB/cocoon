package mapper

import (
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
)

func ProfileToResponse(profile entity.UserProfile) dto.ProfileResponse {
	return dto.ProfileResponse{
		UserID:    profile.UserID,
		ID:        profile.ID,
		Name:      profile.Name,
		Avatar:    profile.Avatar,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
		DeletedAt: profile.DeletedAt.Time,
	}
}
