package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProfileResponse struct {
	UserID    uuid.UUID `json:"userId"`
	ProfileID uuid.UUID `json:"profileId"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	DeletedAt time.Time `json:"deletedAt,omitzero"`
}
