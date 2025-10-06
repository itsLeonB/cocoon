package provider

import (
	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/cocoon/internal/store"
	"github.com/itsLeonB/ezutil/v2"
)

type Services struct {
	Auth       service.AuthService
	Profile    service.ProfileService
	Friendship service.FriendshipService
}

func ProvideServices(
	configs config.Config,
	repos *Repositories,
	logger ezutil.Logger,
	store store.StateStore,
) *Services {
	profileService := service.NewProfileService(
		repos.Transactor,
		repos.UserProfile,
	)

	authService := service.NewAuthService(
		repos.User,
		repos.Transactor,
		profileService,
		repos.OAuthAccount,
		logger,
		configs,
		store,
	)

	friendshipService := service.NewFriendshipService(
		repos.Transactor,
		repos.Friendship,
		profileService,
	)

	return &Services{
		Auth:       authService,
		Profile:    profileService,
		Friendship: friendshipService,
	}
}
