package mapper

import (
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToFriendshipProto(res dto.FriendshipResponse) *friendship.FriendshipResponse {
	return &friendship.FriendshipResponse{
		Id:          res.ID.String(),
		Type:        ToFriendshipTypeEnum(res.Type),
		ProfileId:   res.ProfileID.String(),
		ProfileName: res.ProfileName,
		CreatedAt:   timestamppb.New(res.CreatedAt),
		UpdatedAt:   timestamppb.New(res.UpdatedAt),
		DeletedAt:   timestamppb.New(res.DeletedAt),
	}
}

func ToFriendshipTypeEnum(friendshipType appconstant.FriendshipType) friendship.FriendshipType {
	return friendship.FriendshipType(friendship.FriendshipType_value[string(friendshipType)])
}
