package provider

import (
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/ezutil"
)

type Services struct {
	Auth       service.AuthService
	Profile    service.ProfileService
	Friendship service.FriendshipService
}

func ProvideServices(configs *ezutil.Config, repos *Repositories) *Services {
	hashService := ezutil.NewHashService(0)

	jwtService := ezutil.NewJwtService(configs.Auth)

	profileService := service.NewProfileService(
		repos.Transactor,
		repos.UserProfile,
	)

	authService := service.NewAuthService(
		hashService,
		jwtService,
		repos.User,
		repos.Transactor,
		profileService,
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
