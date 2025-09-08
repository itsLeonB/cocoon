package mapper

import (
	"github.com/itsLeonB/audit/gen/go/audit/v1"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile/v1"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/gerpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProfileProto(res dto.ProfileResponse) *profile.Profile {
	return &profile.Profile{
		UserId:      res.UserID.String(),
		Name:        res.Name,
		IsAnonymous: res.IsAnonymous,
		AuditMetadata: &audit.Metadata{
			Id:        res.ID.String(),
			CreatedAt: timestamppb.New(res.CreatedAt),
			UpdatedAt: timestamppb.New(res.UpdatedAt),
			DeletedAt: gerpc.NullableTimeToProto(res.DeletedAt),
		},
	}
}
