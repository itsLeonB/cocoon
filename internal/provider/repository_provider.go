package provider

import (
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/repository"
	"github.com/itsLeonB/ezutil"
	"gorm.io/gorm"
)

type Repositories struct {
	Transactor  ezutil.Transactor
	User        repository.UserRepository
	UserProfile repository.UserProfileRepository
	Friendship  repository.FriendshipRepository
}

func ProvideRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Transactor:  ezutil.NewTransactor(db),
		User:        ezutil.NewCRUDRepository[entity.User](db),
		UserProfile: repository.NewProfileRepository(db),
		Friendship:  repository.NewFriendshipRepository(db),
	}
}
