package mapper

import (
	"github.com/itsLeonB/cocoon-protos/gen/go/profile/v1"
	"github.com/itsLeonB/cocoon/internal/dto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProfileProto(res dto.ProfileResponse) *profile.ProfileResponse {
	return &profile.ProfileResponse{
		Id:          res.ID.String(),
		UserId:      res.UserID.String(),
		Name:        res.Name,
		CreatedAt:   timestamppb.New(res.CreatedAt),
		UpdatedAt:   timestamppb.New(res.UpdatedAt),
		DeletedAt:   NullableTimeToProto(res.DeletedAt),
		IsAnonymous: res.IsAnonymous,
	}
}
