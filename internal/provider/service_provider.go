package provider

import (
	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/service"
	"github.com/itsLeonB/cocoon/internal/store"
	"github.com/itsLeonB/ezutil/v2"
)

type Services struct {
	Auth       service.AuthService
	OAuth      service.OAuthService
	Profile    service.ProfileService
	Friendship service.FriendshipService
}

func ProvideServices(
	configs config.Config,
	repos *Repositories,
	logger ezutil.Logger,
	store store.StateStore,
) (*Services, error) {
	profileService := service.NewProfileService(
		repos.Transactor,
		repos.UserProfile,
	)

	userSvc := service.NewUserService(
		repos.Transactor,
		repos.User,
		profileService,
	)

	mailSvc := service.NewMailService(configs.Mail)

	authService := service.NewAuthService(
		repos.Transactor,
		configs.Auth,
		userSvc,
		mailSvc,
	)

	oAuthSvc := service.NewOAuthService(
		repos.Transactor,
		repos.OAuthAccount,
		logger,
		configs,
		store,
		userSvc,
	)

	friendshipService := service.NewFriendshipService(
		repos.Transactor,
		repos.Friendship,
		profileService,
	)

	return &Services{
		Auth:       authService,
		OAuth:      oAuthSvc,
		Profile:    profileService,
		Friendship: friendshipService,
	}, nil
}
