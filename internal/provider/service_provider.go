package provider

import (
	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/sekure"
)

type Services struct {
	Auth       service.AuthService
	Profile    service.ProfileService
	Friendship service.FriendshipService
}

func ProvideServices(configs config.Auth, repos *Repositories) *Services {
	hashService := sekure.NewHashService(configs.HashCost)

	jwtService := sekure.NewJwtService(configs.Issuer, configs.SecretKey, configs.TokenDuration)

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
