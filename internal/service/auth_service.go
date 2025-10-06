package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/itsLeonB/cocoon/internal/appconstant"
	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/cocoon/internal/dto"
	"github.com/itsLeonB/cocoon/internal/entity"
	"github.com/itsLeonB/cocoon/internal/mapper"
	"github.com/itsLeonB/cocoon/internal/repository"
	"github.com/itsLeonB/cocoon/internal/service/oauth"
	"github.com/itsLeonB/cocoon/internal/store"
	"github.com/itsLeonB/cocoon/internal/util"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/itsLeonB/go-crud"
	"github.com/itsLeonB/sekure"
	"github.com/itsLeonB/ungerr"
	"github.com/rotisserie/eris"
)

type authServiceImpl struct {
	hashService      sekure.HashService
	jwtService       sekure.JWTService
	userRepository   repository.UserRepository
	transactor       crud.Transactor
	profileService   ProfileService
	oauthProviders   map[string]oauth.ProviderService
	oauthAccountRepo crud.Repository[entity.OAuthAccount]
	stateStore       store.StateStore
}

func NewAuthService(
	userRepository repository.UserRepository,
	transactor crud.Transactor,
	profileService ProfileService,
	oauthAccountRepo crud.Repository[entity.OAuthAccount],
	logger ezutil.Logger,
	configs config.Config,
	stateStore store.StateStore,
) AuthService {
	return &authServiceImpl{
		sekure.NewHashService(configs.HashCost),
		sekure.NewJwtService(configs.Issuer, configs.SecretKey, configs.TokenDuration),
		userRepository,
		transactor,
		profileService,
		oauth.NewOAuthProviderServices(logger, configs.OAuthProviders),
		oauthAccountRepo,
		stateStore,
	}
}

func (as *authServiceImpl) Register(ctx context.Context, request dto.RegisterRequest) error {
	return as.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		existingUser, err := as.findUserByEmail(ctx, request.Email)
		if err != nil {
			return err
		}
		if !existingUser.IsZero() {
			return ungerr.ConflictError(fmt.Sprintf(appconstant.ErrAuthDuplicateUser, request.Email))
		}

		hash, err := as.hashService.Hash(request.Password)
		if err != nil {
			return err
		}

		newUser := entity.User{
			Email:    request.Email,
			Password: hash,
		}

		user, err := as.userRepository.Insert(ctx, newUser)
		if err != nil {
			return err
		}

		profile := dto.NewProfileRequest{
			UserID: user.ID,
			Name:   util.GetNameFromEmail(request.Email),
		}

		if _, err = as.profileService.Create(ctx, profile); err != nil {
			return err
		}

		return nil
	})
}

func (as *authServiceImpl) Login(ctx context.Context, request dto.LoginRequest) (dto.LoginResponse, error) {
	user, err := as.findUserByEmail(ctx, request.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if user.IsZero() {
		return dto.LoginResponse{}, ungerr.NotFoundError(appconstant.ErrAuthUnknownCredentials)
	}

	ok, err := as.hashService.CheckHash(user.Password, request.Password)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	if !ok {
		return dto.LoginResponse{}, ungerr.NotFoundError(appconstant.ErrAuthUnknownCredentials)
	}

	return as.createLoginResponse(user)
}

func (as *authServiceImpl) VerifyToken(ctx context.Context, token string) (dto.AuthData, error) {
	claims, err := as.jwtService.VerifyToken(token)
	if err != nil {
		return dto.AuthData{}, err
	}

	tokenUserId, exists := claims.Data[appconstant.ContextUserID]
	if !exists {
		return dto.AuthData{}, eris.New("missing user ID from token")
	}
	stringUserID, ok := tokenUserId.(string)
	if !ok {
		return dto.AuthData{}, eris.New("error asserting userID, is not a string")
	}
	userID, err := ezutil.Parse[uuid.UUID](stringUserID)
	if err != nil {
		return dto.AuthData{}, err
	}

	spec := crud.Specification[entity.User]{}
	spec.Model.ID = userID
	spec.PreloadRelations = []string{"Profile"}
	user, err := as.userRepository.FindFirst(ctx, spec)
	if err != nil {
		return dto.AuthData{}, err
	}
	if user.IsZero() {
		return dto.AuthData{}, eris.New("user ID is not found")
	}
	if user.IsDeleted() {
		return dto.AuthData{}, eris.New("user is deleted")
	}

	return dto.AuthData{
		ProfileID: user.Profile.ID,
	}, nil
}

func (as *authServiceImpl) GetOAuthURL(ctx context.Context, provider string) (string, error) {
	oauthProvider, ok := as.oauthProviders[provider]
	if !ok {
		return "", eris.Errorf("unsupported oauth provider: %s", provider)
	}

	state, err := as.generateState()
	if err != nil {
		return "", err
	}

	url, err := oauthProvider.GetAuthCodeURL(ctx, state)
	if err != nil {
		return "", err
	}

	if err = as.stateStore.Store(ctx, state, 5*time.Minute); err != nil {
		return "", err
	}

	return url, nil
}

func (as *authServiceImpl) HandleOAuthCallback(ctx context.Context, data dto.OAuthCallbackData) (dto.LoginResponse, error) {
	var response dto.LoginResponse
	err := as.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
		oauthProvider, ok := as.oauthProviders[data.Provider]
		if !ok {
			return eris.Errorf("unsupported oauth provider: %s", data.Provider)
		}

		stateIsValid, err := as.stateStore.VerifyAndDelete(ctx, data.State)
		if err != nil {
			return err
		}
		if !stateIsValid {
			return ungerr.BadRequestError("state is not valid")
		}

		userInfo, err := oauthProvider.HandleCallback(ctx, data.Code)
		if err != nil {
			return err
		}

		existingOAuth, err := as.findOAuthAccount(ctx, userInfo.Provider, userInfo.ProviderID)
		if err != nil {
			return err
		}
		if !existingOAuth.IsZero() && !existingOAuth.IsDeleted() {
			if existingOAuth.User.IsDeleted() {
				return ungerr.NotFoundError(appconstant.ErrAuthUnknownCredentials)
			}
			response, err = as.createLoginResponse(existingOAuth.User)
			return err
		}

		user, err := as.createNewUserOAuth(ctx, userInfo)
		if err != nil {
			return err
		}
		response, err = as.createLoginResponse(user)
		return err
	})

	return response, err
}

func (as *authServiceImpl) createNewUserOAuth(ctx context.Context, userInfo oauth.UserInfo) (entity.User, error) {
	user, err := as.findUserByEmail(ctx, userInfo.Email)
	if err != nil {
		return entity.User{}, err
	}
	if user.IsZero() {
		// New user
		newUser := entity.User{Email: userInfo.Email}
		user, err = as.userRepository.Insert(ctx, newUser)
		if err != nil {
			return entity.User{}, err
		}
		newProfile := dto.NewProfileRequest{
			UserID: user.ID,
			Name:   userInfo.Name,
			Avatar: userInfo.Avatar,
		}
		if _, err = as.profileService.Create(ctx, newProfile); err != nil {
			return entity.User{}, err
		}
	} else if user.IsDeleted() {
		return entity.User{}, ungerr.NotFoundError(appconstant.ErrAuthUnknownCredentials)
	}

	if !as.oauthProviders[userInfo.Provider].IsTrusted() {
		return entity.User{}, eris.New("provider temporarily disabled")
	}

	// New oauth method
	newOAuthAccount := entity.OAuthAccount{
		UserID:     user.ID,
		Provider:   userInfo.Provider,
		ProviderID: userInfo.ProviderID,
		Email:      userInfo.Email,
	}

	if _, err = as.oauthAccountRepo.Insert(ctx, newOAuthAccount); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (as *authServiceImpl) createLoginResponse(user entity.User) (dto.LoginResponse, error) {
	authData := mapper.UserToAuthData(user)

	token, err := as.jwtService.CreateToken(authData)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	return dto.NewBearerTokenResp(token), nil
}

func (as *authServiceImpl) findOAuthAccount(ctx context.Context, provider, providerID string) (entity.OAuthAccount, error) {
	oauthSpec := crud.Specification[entity.OAuthAccount]{}
	oauthSpec.Model.Provider = provider
	oauthSpec.Model.ProviderID = providerID
	oauthSpec.DeletedFilter = crud.ExcludeDeleted
	oauthSpec.PreloadRelations = []string{"User"}
	return as.oauthAccountRepo.FindFirst(ctx, oauthSpec)
}

func (as *authServiceImpl) findUserByEmail(ctx context.Context, email string) (entity.User, error) {
	userSpec := crud.Specification[entity.User]{}
	userSpec.Model.Email = email
	userSpec.DeletedFilter = crud.ExcludeDeleted
	return as.userRepository.FindFirst(ctx, userSpec)
}

func (as *authServiceImpl) generateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", eris.Wrap(err, "error generating random string")
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
