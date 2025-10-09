package dto

import (
	"time"

	"github.com/google/uuid"
)

type FriendshipRequestResponse struct {
	ID        uuid.UUID
	Sender    ProfileResponse
	Recipient ProfileResponse
	Message   string
	CreatedAt time.Time
	BlockedAt time.Time
}
