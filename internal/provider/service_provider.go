package provider

import (
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/ezutil"
)

type Services struct {
	Auth       service.AuthService
	User       service.UserService
	Friendship service.FriendshipService
}

func ProvideServices(configs *ezutil.Config, repos *Repositories) *Services {
	hashService := ezutil.NewHashService(0)

	jwtService := ezutil.NewJwtService(configs.Auth)

	authService := service.NewAuthService(
		hashService,
		jwtService,
		repos.User,
		repos.Transactor,
		repos.UserProfile,
	)

	userService := service.NewUserService(
		repos.Transactor,
		repos.User,
		repos.UserProfile,
	)

	friendshipService := service.NewFriendshipService(
		repos.Transactor,
		repos.UserProfile,
		repos.Friendship,
		userService,
	)

	return &Services{
		Auth:       authService,
		User:       userService,
		Friendship: friendshipService,
	}
}
