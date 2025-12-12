package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon-protos/gen/go/profile/v1"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc/mapper"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/ungerr"
	"google.golang.org/protobuf/types/known/emptypb"
)

type profileServer struct {
	profile.UnimplementedProfileServiceServer
	validate       *validator.Validate
	profileService service.ProfileService
}

func newProfileServer(
	validate *validator.Validate,
	profileService service.ProfileService,
) profile.ProfileServiceServer {
	return &profileServer{
		validate:       validate,
		profileService: profileService,
	}
}

func (ps *profileServer) Get(ctx context.Context, req *profile.GetRequest) (*profile.GetResponse, error) {
	id, err := uuid.Parse(req.GetProfileId())
	if err != nil {
		return nil, ungerr.ValidationError("profile_id is not a valid uuid")
	}

	prof, err := ps.profileService.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &profile.GetResponse{
		Profile: mapper.ToProfileResponseProto(prof),
	}, nil
}

func (ps *profileServer) Create(ctx context.Context, req *profile.CreateRequest) (*profile.CreateResponse, error) {
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
		Profile: mapper.ToProfileResponseProto(createdProfile),
	}, nil
}

func (ps *profileServer) GetByIDs(ctx context.Context, req *profile.GetByIDsRequest) (*profile.GetByIDsResponse, error) {
	profileIDs, err := ezutil.MapSliceWithError(req.GetProfileIds(), ezutil.Parse[uuid.UUID])
	if err != nil {
		return nil, err
	}

	profiles, err := ps.profileService.GetByIDs(ctx, profileIDs)
	if err != nil {
		return nil, err
	}

	responses := ezutil.MapSlice(profiles, mapper.ToProfileResponseProto)

	return &profile.GetByIDsResponse{Profiles: responses}, nil
}

func (ps *profileServer) Update(ctx context.Context, req *profile.UpdateRequest) (*profile.UpdateResponse, error) {
	request, err := mapper.FromUpdateProfileRequestProto(req)
	if err != nil {
		return nil, err
	}

	response, err := ps.profileService.Update(ctx, request)
	if err != nil {
		return nil, err
	}

	return &profile.UpdateResponse{
		Profile: mapper.ToProfileResponseProto(response),
	}, nil
}

func (ps *profileServer) GetByEmail(ctx context.Context, req *profile.GetByEmailRequest) (*profile.GetByEmailResponse, error) {
	if req == nil {
		return nil, ungerr.BadRequestError("request is nil")
	}
	if req.GetEmail() == "" {
		return nil, ungerr.BadRequestError("email is empty")
	}

	response, err := ps.profileService.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		return nil, err
	}

	return &profile.GetByEmailResponse{
		Profile: mapper.ToProfileResponseProto(response),
	}, nil
}

func (ps *profileServer) SearchByName(ctx context.Context, req *profile.SearchByNameRequest) (*profile.SearchByNameResponse, error) {
	if req == nil {
		return nil, ungerr.BadRequestError("request is nil")
	}
	if req.GetQuery() == "" {
		return nil, ungerr.BadRequestError("query is empty")
	}
	if req.GetLimit() < 1 {
		return nil, ungerr.BadRequestError("limit must be greater than 0")
	}
	if req.GetLimit() > 100 {
		return nil, ungerr.BadRequestError("limit must be less than or equal to 100")
	}

	responses, err := ps.profileService.SearchByName(ctx, req.GetQuery(), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}

	return &profile.SearchByNameResponse{
		Profiles: ezutil.MapSlice(responses, mapper.ToProfileResponseProto),
	}, nil
}

func (ps *profileServer) Associate(ctx context.Context, req *profile.AssociateRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, ungerr.BadRequestError("request is nil")
	}
	userProfileID, err := ezutil.Parse[uuid.UUID](req.GetUserProfileId())
	if err != nil {
		return nil, err
	}
	realProfileID, err := ezutil.Parse[uuid.UUID](req.GetRealProfileId())
	if err != nil {
		return nil, err
	}
	anonProfileID, err := ezutil.Parse[uuid.UUID](req.GetAnonProfileId())
	if err != nil {
		return nil, err
	}

	request := dto.AssociateProfileRequest{
		UserProfileID: userProfileID,
		RealProfileID: realProfileID,
		AnonProfileID: anonProfileID,
	}

	return nil, ps.profileService.Associate(ctx, request)
}
