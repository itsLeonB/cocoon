package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
)

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) error
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
	VerifyToken(ctx context.Context, token string) (dto.AuthData, error)
	GetOAuthURL(ctx context.Context, provider, state string) (string, error)
	HandleOAuthCallback(ctx context.Context, provider, code, state string) (dto.LoginResponse, error)
}

type ProfileService interface {
	Create(ctx context.Context, request dto.NewProfileRequest) (dto.ProfileResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]dto.ProfileResponse, error)
	Update(ctx context.Context, req dto.UpdateProfileRequest) (dto.ProfileResponse, error)
}

type FriendshipService interface {
	CreateAnonymous(ctx context.Context, request dto.NewAnonymousFriendshipRequest) (dto.FriendshipResponse, error)
	GetAll(ctx context.Context, profileID uuid.UUID) ([]dto.FriendshipResponse, error)
	GetDetails(ctx context.Context, profileID, friendshipID uuid.UUID) (dto.FriendDetails, error)
	IsFriends(ctx context.Context, profileID1, profileID2 uuid.UUID) (bool, bool, error)
}
