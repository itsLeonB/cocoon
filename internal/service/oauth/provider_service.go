package oauth

import (
	"context"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/ezutil/v2"
)

type ProviderService interface {
	IsTrusted() bool
	GetAuthCodeURL(ctx context.Context, state string) (string, error)
	HandleCallback(ctx context.Context, code, state string) (UserInfo, error)
}

func NewOAuthProviderServices(logger ezutil.Logger, cfgs config.OAuthProviders) map[string]ProviderService {
	return map[string]ProviderService{
		"google": newGoogleProviderService(logger, cfgs.Google),
	}
}
