package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
)

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) error
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
}

type UserService interface {
	GetProfile(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error)
	GetEntityByID(ctx context.Context, id uuid.UUID) (entity.User, error)
}

type FriendshipService interface {
	CreateAnonymous(ctx context.Context, request dto.NewAnonymousFriendshipRequest) (dto.FriendshipResponse, error)
	GetAll(ctx context.Context, userID uuid.UUID) ([]dto.FriendshipResponse, error)
	GetDetails(ctx context.Context, userID, friendshipID uuid.UUID) (dto.FriendDetails, error)
	IsFriends(ctx context.Context, profileID1, profileID2 uuid.UUID) (bool, error)
}
