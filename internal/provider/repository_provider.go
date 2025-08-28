package provider

import (
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/repository"
	crud "github.com/itsLeonB/go-crud"
	"gorm.io/gorm"
)

type Repositories struct {
	Transactor  crud.Transactor
	User        repository.UserRepository
	UserProfile repository.UserProfileRepository
	Friendship  repository.FriendshipRepository
}

func ProvideRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Transactor:  crud.NewTransactor(db),
		User:        crud.NewCRUDRepository[entity.User](db),
		UserProfile: repository.NewProfileRepository(db),
		Friendship:  repository.NewFriendshipRepository(db),
	}
}
