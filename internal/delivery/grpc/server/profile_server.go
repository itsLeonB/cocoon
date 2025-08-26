package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile/v1"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc/mapper"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/ezutil"
)

type ProfileServer struct {
	profile.UnimplementedProfileServiceServer
	validate       *validator.Validate
	profileService service.ProfileService
}

func NewProfileServer(
	validate *validator.Validate,
	profileService service.ProfileService,
) profile.ProfileServiceServer {
	return &ProfileServer{
		validate:       validate,
		profileService: profileService,
	}
}

func (ps *ProfileServer) Get(ctx context.Context, req *profile.GetRequest) (*profile.GetResponse, error) {
	id, err := uuid.Parse(req.GetProfileId())
	if err != nil {
		return nil, ezutil.ValidationError("profile_id is not a valid uuid")
	}

	prof, err := ps.profileService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &profile.GetResponse{
		Profile: mapper.ToProfileProto(prof),
	}, nil
}

func (ps *ProfileServer) Create(ctx context.Context, req *profile.CreateRequest) (*profile.CreateResponse, error) {
	var userID uuid.UUID
	if req.GetUserId() != "" {
		parsedID, err := ezutil.Parse[uuid.UUID](req.GetUserId())
		if err != nil {
			return nil, err
		}
		userID = parsedID
	}

	request := dto.NewProfileRequest{
		UserID: userID,
		Name:   req.GetName(),
	}

	if err := ps.validate.Struct(request); err != nil {
		return nil, err
	}

	createdProfile, err := ps.profileService.Create(ctx, request)
	if err != nil {
		return nil, err
	}

	return &profile.CreateResponse{
		Profile: mapper.ToProfileProto(createdProfile),
	}, nil
}

func (ps *ProfileServer) GetNames(ctx context.Context, req *profile.GetNamesRequest) (*profile.GetNamesResponse, error) {
	profileIDs, err := ezutil.MapSliceWithError(req.GetProfileIds(), ezutil.Parse[uuid.UUID])
	if err != nil {
		return nil, err
	}

	namesMap, err := ps.profileService.GetNames(ctx, profileIDs)
	if err != nil {
		return nil, err
	}

	namesByProfileID := make(map[string]string, len(namesMap))
	for id, name := range namesMap {
		namesByProfileID[id.String()] = name
	}

	return &profile.GetNamesResponse{
		NamesByProfileId: namesByProfileID,
	}, nil
}

func (ps *ProfileServer) GetByIDs(ctx context.Context, req *profile.GetByIDsRequest) (*profile.GetByIDsResponse, error) {
	profileIDs, err := ezutil.MapSliceWithError(req.GetProfileIds(), ezutil.Parse[uuid.UUID])
	if err != nil {
		return nil, err
	}

	profiles, err := ps.profileService.GetByIDs(ctx, profileIDs)
	if err != nil {
		return nil, err
	}

	responses := ezutil.MapSlice(profiles, mapper.ToProfileProto)

	return &profile.GetByIDsResponse{Profiles: responses}, nil
}
