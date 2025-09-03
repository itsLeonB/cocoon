package mapper

import (
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/gerpc"
	"github.com/rotisserie/eris"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToFriendshipProto(res dto.FriendshipResponse) (*friendship.FriendshipResponse, error) {
	friendshipType, err := ToProtoFriendshipType(res.Type)
	if err != nil {
		return nil, err
	}

	return &friendship.FriendshipResponse{
		Id:          res.ID.String(),
		Type:        friendshipType,
		ProfileId:   res.ProfileID.String(),
		ProfileName: res.ProfileName,
		CreatedAt:   timestamppb.New(res.CreatedAt),
		UpdatedAt:   timestamppb.New(res.UpdatedAt),
		DeletedAt:   gerpc.NullableTimeToProto(res.DeletedAt),
	}, nil
}

func ToProtoFriendshipType(ft appconstant.FriendshipType) (friendship.FriendshipType, error) {
	switch ft {
	case appconstant.Real:
		return friendship.FriendshipType_FRIENDSHIP_TYPE_REAL, nil
	case appconstant.Anonymous:
		return friendship.FriendshipType_FRIENDSHIP_TYPE_ANON, nil
	default:
		return friendship.FriendshipType_FRIENDSHIP_TYPE_UNSPECIFIED, eris.Errorf("undefined FriendshipType constant: %s", ft)
	}
}
