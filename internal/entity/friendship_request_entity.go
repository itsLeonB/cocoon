package entity

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/itsLeonB/go-crud"
)

type FriendshipRequest struct {
	crud.BaseEntity
	SenderProfileID    uuid.UUID
	RecipientProfileID uuid.UUID
	Message            sql.NullString
	BlockedAt          sql.NullTime
	SenderProfile      UserProfile
	RecipientProfile   UserProfile
}
