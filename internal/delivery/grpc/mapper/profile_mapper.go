package mapper

import (
	"github.com/google/uuid"
	"github.com/itsLeonB/audit/gen/go/audit/v1"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile/v1"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/gerpc"
	"github.com/rotisserie/eris"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToProfileResponseProto(res dto.ProfileResponse) *profile.ProfileResponse {
	return &profile.ProfileResponse{
		Profile: &profile.Profile{
			UserId: res.UserID.String(),
			Name:   res.Name,
			Avatar: res.Avatar,
		},
		AuditMetadata: &audit.Metadata{
			Id:        res.ID.String(),
			CreatedAt: timestamppb.New(res.CreatedAt),
			UpdatedAt: timestamppb.New(res.UpdatedAt),
			DeletedAt: gerpc.NullableTimeToProto(res.DeletedAt),
		},
	}
}

func FromUpdateProfileRequestProto(req *profile.UpdateRequest) (dto.UpdateProfileRequest, error) {
	if req == nil {
		return dto.UpdateProfileRequest{}, eris.New("request is nil")
	}
	if req.GetId() == "" {
		return dto.UpdateProfileRequest{}, eris.New("id is nil")
	}
	id, err := ezutil.Parse[uuid.UUID](req.GetId())
	if err != nil {
		return dto.UpdateProfileRequest{}, err
	}
	profile := req.GetProfile()
	if profile == nil {
		return dto.UpdateProfileRequest{}, eris.New("profile is nil")
	}

	userID := uuid.Nil
	if profile.GetUserId() != "" {
		parsedID, err := ezutil.Parse[uuid.UUID](profile.GetUserId())
		if err != nil {
			return dto.UpdateProfileRequest{}, err
		}
		userID = parsedID
	}

	return dto.UpdateProfileRequest{
		ID:     id,
		UserID: userID,
		Name:   profile.GetName(),
		Avatar: profile.GetAvatar(),
	}, nil
}
