package provider

import (
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/repository"
	"github.com/itsLeonB/go-crud"
	"gorm.io/gorm"
)

type Repositories struct {
	Transactor   crud.Transactor
	User         crud.Repository[entity.User]
	UserProfile  repository.UserProfileRepository
	Friendship   repository.FriendshipRepository
	OAuthAccount crud.Repository[entity.OAuthAccount]
}

func ProvideRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		Transactor:   crud.NewTransactor(db),
		User:         crud.NewRepository[entity.User](db),
		UserProfile:  repository.NewProfileRepository(db),
		Friendship:   repository.NewFriendshipRepository(db),
		OAuthAccount: crud.NewRepository[entity.OAuthAccount](db),
	}
}
