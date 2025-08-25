package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/appconstant"
)

type NewAnonymousFriendshipRequest struct {
	ProfileID uuid.UUID
	Name      string `json:"name" validate:"required,min=3"`
}

type FriendshipResponse struct {
	ID          uuid.UUID                  `json:"id"`
	Type        appconstant.FriendshipType `json:"type"`
	ProfileID   uuid.UUID                  `json:"profileId"`
	ProfileName string                     `json:"profileName"`
	CreatedAt   time.Time                  `json:"createdAt"`
	UpdatedAt   time.Time                  `json:"updatedAt"`
	DeletedAt   time.Time                  `json:"deletedAt,omitzero"`
}

type FriendshipWithProfile struct {
	Friendship    FriendshipResponse
	UserProfile   ProfileResponse
	FriendProfile ProfileResponse
}

type FriendDetails struct {
	ID         uuid.UUID                  `json:"id"`
	ProfileID  uuid.UUID                  `json:"profileId"`
	Name       string                     `json:"name"`
	Type       appconstant.FriendshipType `json:"type"`
	Email      string                     `json:"email,omitempty"`
	Phone      string                     `json:"phone,omitempty"`
	Avatar     string                     `json:"avatar,omitempty"`
	CreatedAt  time.Time                  `json:"createdAt"`
	UpdatedAt  time.Time                  `json:"updatedAt"`
	DeletedAt  time.Time                  `json:"deletedAt,omitzero"`
	ProfileID1 uuid.UUID
	ProfileID2 uuid.UUID
}
