package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
)

type AuthService interface {
	Register(ctx context.Context, request dto.RegisterRequest) (bool, error)
	Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error)
	VerifyToken(ctx context.Context, token string) (dto.AuthData, error)
	VerifyRegistration(ctx context.Context, token string) (dto.LoginResponse, error)
	SendResetPassword(ctx context.Context, resetURL, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) (dto.LoginResponse, error)
}

type OAuthService interface {
	GetOAuthURL(ctx context.Context, provider string) (string, error)
	HandleOAuthCallback(ctx context.Context, data dto.OAuthCallbackData) (dto.LoginResponse, error)
}

type UserService interface {
	CreateNew(ctx context.Context, request dto.NewUserRequest) (entity.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.UserResponse, error)
	FindByEmail(ctx context.Context, email string) (entity.User, error)
	Verify(ctx context.Context, id uuid.UUID, email string, name string, avatar string) (entity.User, error)
	GeneratePasswordResetToken(ctx context.Context, userID uuid.UUID) (string, error)
	ResetPassword(ctx context.Context, userID uuid.UUID, email, resetToken, password string) (entity.User, error)
}

type ProfileService interface {
	Create(ctx context.Context, request dto.NewProfileRequest) (dto.ProfileResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (dto.ProfileResponse, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]dto.ProfileResponse, error)
	Update(ctx context.Context, req dto.UpdateProfileRequest) (dto.ProfileResponse, error)
	GetByEmail(ctx context.Context, email string) (dto.ProfileResponse, error)
	SearchByName(ctx context.Context, query string, limit int) ([]dto.ProfileResponse, error)
}

type FriendshipService interface {
	CreateAnonymous(ctx context.Context, request dto.NewAnonymousFriendshipRequest) (dto.FriendshipResponse, error)
	GetAll(ctx context.Context, profileID uuid.UUID) ([]dto.FriendshipResponse, error)
	GetDetails(ctx context.Context, profileID, friendshipID uuid.UUID) (dto.FriendDetails, error)
	IsFriends(ctx context.Context, profileID1, profileID2 uuid.UUID) (bool, bool, error)
	CreateReal(ctx context.Context, userProfileID, friendProfileID uuid.UUID) (dto.FriendshipResponse, error)
}

type FriendshipRequestService interface {
	Send(ctx context.Context, userProfileID, friendProfileID uuid.UUID) error
	GetAllSent(ctx context.Context, userProfileID uuid.UUID) ([]dto.FriendshipRequestResponse, error)
	Cancel(ctx context.Context, userProfileID, reqID uuid.UUID) error
	GetAllReceived(ctx context.Context, userProfileID uuid.UUID) ([]dto.FriendshipRequestResponse, error)
	Ignore(ctx context.Context, userProfileID, reqID uuid.UUID) error
	Block(ctx context.Context, userProfileID, reqID uuid.UUID) error
	Unblock(ctx context.Context, userProfileID, reqID uuid.UUID) error
	Accept(ctx context.Context, userProfileID, reqID uuid.UUID) (dto.FriendshipResponse, error)
}

type MailService interface {
	Send(ctx context.Context, msg dto.MailMessage) error
}
