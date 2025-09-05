package mapper

import (
	"github.com/itsLeonB/cocoon-protos/gen/go/domain/v1"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/gerpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProfileProto(res dto.ProfileResponse) *domain.Profile {
	return &domain.Profile{
		UserId:      res.UserID.String(),
		Name:        res.Name,
		IsAnonymous: res.IsAnonymous,
		AuditMetadata: &domain.AuditMetadata{
			Id:        res.ID.String(),
			CreatedAt: timestamppb.New(res.CreatedAt),
			UpdatedAt: timestamppb.New(res.UpdatedAt),
			DeletedAt: gerpc.NullableTimeToProto(res.DeletedAt),
		},
	}
}
