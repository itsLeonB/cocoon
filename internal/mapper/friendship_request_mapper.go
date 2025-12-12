package mapper

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
)

func FriendshipRequestToResponse(fr entity.FriendshipRequest) dto.FriendshipRequestResponse {
	return dto.FriendshipRequestResponse{
		ID:        fr.ID,
		Sender:    ProfileToResponse(fr.SenderProfile, "", nil, uuid.Nil),
		Recipient: ProfileToResponse(fr.RecipientProfile, "", nil, uuid.Nil),
		CreatedAt: fr.CreatedAt,
		BlockedAt: fr.BlockedAt.Time,
	}
}
