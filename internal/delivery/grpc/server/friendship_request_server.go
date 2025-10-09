package server

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon-protos/gen/go/friendship/v1"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/delivery/grpc/mapper"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
	"google.golang.org/protobuf/types/known/emptypb"
)

type friendshipRequestServer struct {
	friendship.UnimplementedRequestServiceServer
	svc service.FriendshipRequestService
}

func newFriendshipRequestServer(svc service.FriendshipRequestService) friendship.RequestServiceServer {
	return &friendshipRequestServer{
		svc: svc,
	}
}

func (frs *friendshipRequestServer) Send(ctx context.Context, req *friendship.SendRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	userProfileID, err := ezutil.Parse[uuid.UUID](req.GetUserProfileId())
	if err != nil {
		return nil, err
	}
	friendProfileID, err := ezutil.Parse[uuid.UUID](req.GetFriendProfileId())
	if err != nil {
		return nil, err
	}
	return nil, frs.svc.Send(ctx, userProfileID, friendProfileID, req.GetMessage())
}

func (frs *friendshipRequestServer) GetAllSent(ctx context.Context, req *friendship.GetAllSentRequest) (*friendship.GetAllSentResponse, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	userProfileID, err := ezutil.Parse[uuid.UUID](req.GetUserProfileId())
	if err != nil {
		return nil, err
	}

	response, err := frs.svc.GetAllSent(ctx, userProfileID)
	if err != nil {
		return nil, err
	}

	return &friendship.GetAllSentResponse{
		Requests: ezutil.MapSlice(response, mapper.ToFriendshipRequestProto),
	}, nil
}

func (frs *friendshipRequestServer) Cancel(ctx context.Context, req *friendship.CancelRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	userProfileID, err := ezutil.Parse[uuid.UUID](req.GetUserProfileId())
	if err != nil {
		return nil, err
	}
	reqID, err := ezutil.Parse[uuid.UUID](req.GetRequestId())
	if err != nil {
		return nil, err
	}
	return nil, frs.svc.Cancel(ctx, userProfileID, reqID)
}

func (frs *friendshipRequestServer) GetAllReceived(ctx context.Context, req *friendship.GetAllReceivedRequest) (*friendship.GetAllReceivedResponse, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	userProfileID, err := ezutil.Parse[uuid.UUID](req.GetUserProfileId())
	if err != nil {
		return nil, err
	}

	response, err := frs.svc.GetAllReceived(ctx, userProfileID)
	if err != nil {
		return nil, err
	}

	return &friendship.GetAllReceivedResponse{
		Requests: ezutil.MapSlice(response, mapper.ToFriendshipRequestProto),
	}, nil
}

func (frs *friendshipRequestServer) Ignore(ctx context.Context, req *friendship.IgnoreRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	userProfileID, err := ezutil.Parse[uuid.UUID](req.GetUserProfileId())
	if err != nil {
		return nil, err
	}
	reqID, err := ezutil.Parse[uuid.UUID](req.GetRequestId())
	if err != nil {
		return nil, err
	}
	return nil, frs.svc.Ignore(ctx, userProfileID, reqID)
}

func (frs *friendshipRequestServer) Block(ctx context.Context, req *friendship.BlockRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	userProfileID, err := ezutil.Parse[uuid.UUID](req.GetUserProfileId())
	if err != nil {
		return nil, err
	}
	reqID, err := ezutil.Parse[uuid.UUID](req.GetRequestId())
	if err != nil {
		return nil, err
	}
	return nil, frs.svc.Block(ctx, userProfileID, reqID)
}

func (frs *friendshipRequestServer) Unblock(ctx context.Context, req *friendship.UnblockRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	userProfileID, err := ezutil.Parse[uuid.UUID](req.GetUserProfileId())
	if err != nil {
		return nil, err
	}
	reqID, err := ezutil.Parse[uuid.UUID](req.GetRequestId())
	if err != nil {
		return nil, err
	}
	return nil, frs.svc.Unblock(ctx, userProfileID, reqID)
}

func (frs *friendshipRequestServer) Accept(ctx context.Context, req *friendship.AcceptRequest) (*friendship.AcceptResponse, error) {
	if req == nil {
		return nil, eris.New(appconstant.ErrNilRequest)
	}
	userProfileID, err := ezutil.Parse[uuid.UUID](req.GetUserProfileId())
	if err != nil {
		return nil, err
	}
	reqID, err := ezutil.Parse[uuid.UUID](req.GetRequestId())
	if err != nil {
		return nil, err
	}

	friendshipResp, err := frs.svc.Accept(ctx, userProfileID, reqID)
	if err != nil {
		return nil, err
	}

	resp, err := mapper.ToFriendshipProto(friendshipResp)
	if err != nil {
		return nil, err
	}

	return &friendship.AcceptResponse{
		Friendship: resp,
	}, nil
}
