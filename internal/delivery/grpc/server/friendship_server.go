package server

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc/mapper"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/ezutil"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FriendshipServer struct {
	friendship.UnimplementedFriendshipServiceServer
	validate          *validator.Validate
	friendshipService service.FriendshipService
}

func NewFriendshipServer(
	validate *validator.Validate,
	friendshipService service.FriendshipService,
) friendship.FriendshipServiceServer {
	return &FriendshipServer{
		validate:          validate,
		friendshipService: friendshipService,
	}
}

func (fs *FriendshipServer) CreateAnonymous(ctx context.Context, req *friendship.CreateAnonymousRequest) (*friendship.CreateAnonymousResponse, error) {
	profileID, err := ezutil.Parse[uuid.UUID](req.GetProfileId())
	if err != nil {
		return nil, err
	}

	request := dto.NewAnonymousFriendshipRequest{
		ProfileID: profileID,
		Name:      req.GetName(),
	}

	if err := fs.validate.Struct(request); err != nil {
		return nil, err
	}

	response, err := fs.friendshipService.CreateAnonymous(ctx, request)
	if err != nil {
		return nil, err
	}

	return &friendship.CreateAnonymousResponse{
		Friendship: mapper.ToFriendshipProto(response),
	}, nil
}

func (fs *FriendshipServer) GetAll(ctx context.Context, req *friendship.GetAllRequest) (*friendship.GetAllResponse, error) {
	profileID, err := ezutil.Parse[uuid.UUID](req.GetProfileId())
	if err != nil {
		return nil, err
	}

	response, err := fs.friendshipService.GetAll(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return &friendship.GetAllResponse{Friendships: ezutil.MapSlice(response, mapper.ToFriendshipProto)}, nil
}

func (fs *FriendshipServer) GetDetails(ctx context.Context, req *friendship.GetDetailsRequest) (*friendship.GetDetailsResponse, error) {
	profileID, err := ezutil.Parse[uuid.UUID](req.GetProfileId())
	if err != nil {
		return nil, err
	}

	friendshipID, err := ezutil.Parse[uuid.UUID](req.GetFriendshipId())
	if err != nil {
		return nil, err
	}

	response, err := fs.friendshipService.GetDetails(ctx, profileID, friendshipID)
	if err != nil {
		return nil, err
	}

	return &friendship.GetDetailsResponse{
		Id:         response.ID.String(),
		ProfileId:  response.ProfileID.String(),
		Name:       response.Name,
		Type:       mapper.ToProtoFriendshipType(response.Type),
		Email:      response.Email,
		Phone:      response.Phone,
		Avatar:     response.Avatar,
		CreatedAt:  timestamppb.New(response.CreatedAt),
		UpdatedAt:  timestamppb.New(response.UpdatedAt),
		DeletedAt:  timestamppb.New(response.DeletedAt),
		ProfileId1: response.ProfileID1.String(),
		ProfileId2: response.ProfileID2.String(),
	}, nil
}

func (fs *FriendshipServer) IsFriends(ctx context.Context, req *friendship.IsFriendsRequest) (*friendship.IsFriendsResponse, error) {
	profileID1, err := ezutil.Parse[uuid.UUID](req.GetProfileId_1())
	if err != nil {
		return nil, err
	}

	profileID2, err := ezutil.Parse[uuid.UUID](req.GetProfileId_2())
	if err != nil {
		return nil, err
	}

	isFriends, isAnonymous, err := fs.friendshipService.IsFriends(ctx, profileID1, profileID2)
	if err != nil {
		return nil, err
	}

	return &friendship.IsFriendsResponse{
		IsFriends:   isFriends,
		IsAnonymous: isAnonymous,
	}, nil
}
