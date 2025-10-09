package mapper

import (
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/gerpc"
)

func ToFriendshipRequestProto(fr dto.FriendshipRequestResponse) *friendship.Request {
	return &friendship.Request{
		Id:        fr.ID.String(),
		Sender:    ToProfileResponseProto(fr.Sender),
		Recipient: ToProfileResponseProto(fr.Recipient),
		Message:   fr.Message,
		CreatedAt: gerpc.NullableTimeToProto(fr.CreatedAt),
		BlockedAt: gerpc.NullableTimeToProto(fr.BlockedAt),
	}
}
