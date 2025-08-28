package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/entity"
	crud "github.com/itsLeonB/go-crud"
)

type UserRepository interface {
	crud.CRUDRepository[entity.User]
}

type UserProfileRepository interface {
	crud.CRUDRepository[entity.UserProfile]
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]entity.UserProfile, error)
}

type FriendshipRepository interface {
	crud.CRUDRepository[entity.Friendship]
	Insert(ctx context.Context, friendship entity.Friendship) (entity.Friendship, error)
	FindAllBySpec(ctx context.Context, spec entity.FriendshipSpecification) ([]entity.Friendship, error)
	FindFirstBySpec(ctx context.Context, spec entity.FriendshipSpecification) (entity.Friendship, error)
	FindByProfileIDs(ctx context.Context, profileID1, profileID2 uuid.UUID) (entity.Friendship, error)
}
