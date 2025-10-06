package oauth

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/itsLeonB/cocoon/internal/config"
	"github.com/itsLeonB/ezutil/v2"
	"github.com/rotisserie/eris"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type googleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type googleProviderService struct {
	userInfoURL string
	cfg         *oauth2.Config
	logger      ezutil.Logger
}

func newGoogleProviderService(logger ezutil.Logger, cfg config.OAuthProvider) ProviderService {
	return &googleProviderService{
		userInfoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
		cfg: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			RedirectURL:  cfg.RedirectURL,
			Endpoint:     google.Endpoint,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
		},
		logger: logger,
	}
}

func (*googleProviderService) IsTrusted() bool {
	return true
}

func (gps *googleProviderService) GetAuthCodeURL(ctx context.Context, state string) (string, error) {
	url := gps.cfg.AuthCodeURL(state, oauth2.AccessTypeOffline)
	if url == "" {
		return "", eris.Errorf("OAuth2 google provider returns empty string for auth code URL")
	}
	return url, nil
}

func (gps *googleProviderService) HandleCallback(ctx context.Context, code, state string) (UserInfo, error) {
	token, err := gps.cfg.Exchange(ctx, code)
	if err != nil {
		return UserInfo{}, eris.Wrap(err, "error exchange OAuth2 token at callback")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, gps.userInfoURL, nil)
	if err != nil {
		return UserInfo{}, eris.Wrap(err, "error creating new HTTP request")
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return UserInfo{}, eris.Wrap(err, "error making HTTP request")
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			err = eris.Wrap(err, "error closing HTTP response body")
			gps.logger.Error(eris.ToString(err, true))
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return UserInfo{}, eris.Wrap(err, "error reading response body")
	}
	if resp.StatusCode != http.StatusOK {
		return UserInfo{}, eris.Errorf("error getting user info: %s", string(body))
	}

	var userInfo googleUserInfo
	if err = json.Unmarshal(body, &userInfo); err != nil {
		return UserInfo{}, eris.Wrap(err, "error unmarshaling google user info")
	}

	return UserInfo{
		Provider:    "google",
		ProviderID:  userInfo.ID,
		Email:       userInfo.Email,
		Name:        userInfo.Name,
		Avatar:      userInfo.Picture,
		AccessToken: token.AccessToken,
	}, nil
}
