package mapper

import (
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/gerpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToFriendshipProto(res dto.FriendshipResponse) *friendship.FriendshipResponse {
	return &friendship.FriendshipResponse{
		Id:          res.ID.String(),
		Type:        ToProtoFriendshipType(res.Type),
		ProfileId:   res.ProfileID.String(),
		ProfileName: res.ProfileName,
		CreatedAt:   timestamppb.New(res.CreatedAt),
		UpdatedAt:   timestamppb.New(res.UpdatedAt),
		DeletedAt:   gerpc.NullableTimeToProto(res.DeletedAt),
	}
}

func ToProtoFriendshipType(ft appconstant.FriendshipType) friendship.FriendshipType {
	switch ft {
	case appconstant.Real:
		return friendship.FriendshipType_FRIENDSHIP_TYPE_REAL
	case appconstant.Anonymous:
		return friendship.FriendshipType_FRIENDSHIP_TYPE_ANON
	default:
		return friendship.FriendshipType_FRIENDSHIP_TYPE_UNSPECIFIED
	}
}
