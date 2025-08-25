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
}

type ProfileService interface {
	Create(ctx context.Context, request dto.NewProfileRequest) (dto.ProfileResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error)
	GetNames(ctx context.Context, ids []uuid.UUID) (map[uuid.UUID]string, error)
}

type FriendshipService interface {
	CreateAnonymous(ctx context.Context, request dto.NewAnonymousFriendshipRequest) (dto.FriendshipResponse, error)
	GetAll(ctx context.Context, profileID uuid.UUID) ([]dto.FriendshipResponse, error)
	GetDetails(ctx context.Context, profileID, friendshipID uuid.UUID) (dto.FriendDetails, error)
	IsFriends(ctx context.Context, profileID1, profileID2 uuid.UUID) (bool, error)
}
