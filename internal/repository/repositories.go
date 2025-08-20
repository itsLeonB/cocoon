package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/ezutil"
)

type UserRepository interface {
	ezutil.CRUDRepository[entity.User]
}

type UserProfileRepository interface {
	ezutil.CRUDRepository[entity.UserProfile]
}

type FriendshipRepository interface {
	ezutil.CRUDRepository[entity.Friendship]
	Insert(ctx context.Context, friendship entity.Friendship) (entity.Friendship, error)
	FindAllBySpec(ctx context.Context, spec entity.FriendshipSpecification) ([]entity.Friendship, error)
	FindFirstBySpec(ctx context.Context, spec entity.FriendshipSpecification) (entity.Friendship, error)
	FindByProfileIDs(ctx context.Context, profileID1, profileID2 uuid.UUID) (entity.Friendship, error)
}
