package service

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mapper"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/ungerr"
)

type friendshipRequestServiceImpl struct {
	transactor     crud.Transactor
	friendshipSvc  FriendshipService
	profileService ProfileService
	requestRepo    crud.Repository[entity.FriendshipRequest]
}

func NewFriendshipRequestService(
	transactor crud.Transactor,
	friendshipSvc FriendshipService,
	profileService ProfileService,
	requestRepo crud.Repository[entity.FriendshipRequest],
) FriendshipRequestService {
	return &friendshipRequestServiceImpl{
		transactor,
		friendshipSvc,
		profileService,
		requestRepo,
	}
}

func (fs *friendshipRequestServiceImpl) Send(ctx context.Context, userProfileID, friendProfileID uuid.UUID) error {
	return fs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.FriendshipRequest]{}
		spec.Model.SenderProfileID = userProfileID
		spec.Model.RecipientProfileID = friendProfileID
		existingRequest, err := fs.requestRepo.FindFirst(ctx, spec)
		if err != nil {
			return err
		}
		if !existingRequest.IsZero() {
			if existingRequest.BlockedAt.Valid {
				return ungerr.UnprocessableEntityError("user is blocked by recipient")
			}
			return ungerr.UnprocessableEntityError("user still has existing request")
		}

		isFriends, _, err := fs.friendshipSvc.IsFriends(ctx, userProfileID, friendProfileID)
		if err != nil {
			return err
		}
		if isFriends {
			return ungerr.UnprocessableEntityError("already friends")
		}

		friendProfile, err := fs.profileService.GetByID(ctx, friendProfileID)
		if err != nil {
			return err
		}
		if friendProfile.UserID == uuid.Nil {
			return ungerr.UnprocessableEntityError("cannot request friendship with anonymous profile")
		}

		newRequest := entity.FriendshipRequest{
			SenderProfileID:    userProfileID,
			RecipientProfileID: friendProfileID,
		}

		_, err = fs.requestRepo.Insert(ctx, newRequest)
		return err
	})
}

func (fs *friendshipRequestServiceImpl) GetAllSent(ctx context.Context, userProfileID uuid.UUID) ([]dto.FriendshipRequestResponse, error) {
	spec := crud.Specification[entity.FriendshipRequest]{}
	spec.Model.SenderProfileID = userProfileID
	spec.PreloadRelations = []string{"SenderProfile", "RecipientProfile"}
	requests, err := fs.requestRepo.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(requests, mapper.FriendshipRequestToResponse), nil
}

func (fs *friendshipRequestServiceImpl) Cancel(ctx context.Context, userProfileID, reqID uuid.UUID) error {
	return fs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.FriendshipRequest]{}
		spec.Model.ID = reqID
		spec.Model.SenderProfileID = userProfileID
		spec.ForUpdate = true
		request, err := fs.getPendingRequest(ctx, spec)
		if err != nil {
			return err
		}
		return fs.requestRepo.Delete(ctx, request)
	})
}

func (fs *friendshipRequestServiceImpl) GetAllReceived(ctx context.Context, userProfileID uuid.UUID) ([]dto.FriendshipRequestResponse, error) {
	spec := crud.Specification[entity.FriendshipRequest]{}
	spec.Model.RecipientProfileID = userProfileID
	spec.PreloadRelations = []string{"SenderProfile", "RecipientProfile"}
	requests, err := fs.requestRepo.FindAll(ctx, spec)
	if err != nil {
		return nil, err
	}

	return ezutil.MapSlice(requests, mapper.FriendshipRequestToResponse), nil
}

func (fs *friendshipRequestServiceImpl) Ignore(ctx context.Context, userProfileID, reqID uuid.UUID) error {
	return fs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.FriendshipRequest]{}
		spec.Model.ID = reqID
		spec.Model.RecipientProfileID = userProfileID
		spec.ForUpdate = true
		request, err := fs.getPendingRequest(ctx, spec)
		if err != nil {
			return err
		}
		return fs.requestRepo.Delete(ctx, request)
	})
}

func (fs *friendshipRequestServiceImpl) Block(ctx context.Context, userProfileID, reqID uuid.UUID) error {
	return fs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.FriendshipRequest]{}
		spec.Model.ID = reqID
		spec.Model.RecipientProfileID = userProfileID
		spec.ForUpdate = true
		request, err := fs.getRequest(ctx, spec)
		if err != nil {
			return err
		}
		if request.BlockedAt.Valid {
			return nil
		}

		request.BlockedAt = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		_, err = fs.requestRepo.Update(ctx, request)
		return err
	})
}

func (fs *friendshipRequestServiceImpl) Unblock(ctx context.Context, userProfileID, reqID uuid.UUID) error {
	return fs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.FriendshipRequest]{}
		spec.Model.ID = reqID
		spec.Model.RecipientProfileID = userProfileID
		spec.ForUpdate = true
		request, err := fs.getRequest(ctx, spec)
		if err != nil {
			return err
		}
		if !request.BlockedAt.Valid {
			return nil
		}

		request.BlockedAt = sql.NullTime{}
		_, err = fs.requestRepo.Update(ctx, request)
		return err
	})
}

func (fs *friendshipRequestServiceImpl) Accept(ctx context.Context, userProfileID, reqID uuid.UUID) (dto.FriendshipResponse, error) {
	var response dto.FriendshipResponse
	err := fs.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		spec := crud.Specification[entity.FriendshipRequest]{}
		spec.Model.ID = reqID
		spec.Model.RecipientProfileID = userProfileID
		spec.ForUpdate = true
		request, err := fs.getPendingRequest(ctx, spec)
		if err != nil {
			return err
		}

		response, err = fs.friendshipSvc.CreateReal(ctx, userProfileID, request.SenderProfileID)
		if err != nil {
			return err
		}

		return fs.requestRepo.Delete(ctx, request)
	})
	return response, err
}

func (fs *friendshipRequestServiceImpl) getPendingRequest(ctx context.Context, spec crud.Specification[entity.FriendshipRequest]) (entity.FriendshipRequest, error) {
	request, err := fs.getRequest(ctx, spec)
	if err != nil {
		return entity.FriendshipRequest{}, err
	}
	if request.BlockedAt.Valid {
		return entity.FriendshipRequest{}, ungerr.UnprocessableEntityError("sender is blocked")
	}
	return request, nil
}

func (fs *friendshipRequestServiceImpl) getRequest(ctx context.Context, spec crud.Specification[entity.FriendshipRequest]) (entity.FriendshipRequest, error) {
	request, err := fs.requestRepo.FindFirst(ctx, spec)
	if err != nil {
		return entity.FriendshipRequest{}, err
	}
	if request.IsZero() {
		return entity.FriendshipRequest{}, ungerr.NotFoundError("request not found")
	}
	return request, nil
}
