package mapper

import (
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
)

func FriendshipRequestToResponse(fr entity.FriendshipRequest) dto.FriendshipRequestResponse {
	return dto.FriendshipRequestResponse{
		ID:        fr.ID,
		Sender:    ProfileToResponse(fr.SenderProfile, ""),
		Recipient: ProfileToResponse(fr.RecipientProfile, ""),
		Message:   fr.Message.String,
		CreatedAt: fr.CreatedAt,
		BlockedAt: fr.BlockedAt.Time,
	}
}
