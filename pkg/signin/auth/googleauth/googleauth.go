package googleauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/signin/auth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v2"
	oauth2v2 "google.golang.org/api/oauth2/v2"
)

var (
	scopes = []string{
		drive.DriveAppdataScope,
		oauth2v2.UserinfoEmailScope,
		oauth2v2.UserinfoProfileScope,
	}
)

// GoogleAuth is the authorization workflow implementation
// for google oauth api.
type GoogleAuth struct {
	clientCred []byte
	config     *oauth2.Config
}

// GetURL build the authorization workflow url for google oauth api.
func (ga *GoogleAuth) GetURL(redirect, state string) (string, error) {
	if len(redirect) > 0 {
		ga.config.RedirectURL = redirect
	}

	s := utils.GetValueString(state, "0")
	return ga.config.AuthCodeURL(s, oauth2.AccessTypeOffline), nil
}

// Exchange converts the authorization code to access token.
func (ga *GoogleAuth) Exchange(code string) (*auth.Token, error) {
	token, err := ga.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, utils.Error("error on converting auth code into token", err)
	}

	return &auth.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		Expiry:       token.Expiry,
	}, nil
}

// Refresh renews access token using refresh token.
func (ga *GoogleAuth) Refresh(refreshToken string) (*auth.Token, error) {
	at, expiresIn, err := ga.doRefresh(refreshToken)
	if err != nil {
		msg := "error on doing token refresh"
		return nil, utils.Error(msg, err)
	}

	return &auth.Token{
		AccessToken:  at,
		RefreshToken: refreshToken,
		Expiry:       time.Now().Add(time.Duration(expiresIn) * time.Second),
	}, nil
}

// Revoke cancels the access token and resets the authorization workflow.
func (ga *GoogleAuth) Revoke(accessToken string) error {
	client := http.DefaultClient
	values := url.Values{}
	values.Set("token", accessToken)
	revokeURL := toRevokeURL(ga.config.Endpoint.TokenURL)

	res, err := client.PostForm(revokeURL, values)
	if err != nil {
		return utils.Error("access token revocation failed", err)
	}

	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("access token revocation failed with '%d' status code", res.StatusCode)
		return errors.New(msg)
	}
	return nil
}

// New creates a new instance of google auth struct.
func New(cred []byte) (*GoogleAuth, error) {
	config, err := google.ConfigFromJSON(cred, scopes...)
	if err != nil {
		msg := "error occurred while unmarshalling google client cred"
		return nil, utils.Error(msg, err)
	}

	return &GoogleAuth{
		clientCred: cred,
		config:     config,
	}, nil
}

func (ga *GoogleAuth) doRefresh(refreshToken string) (string, int, error) {
	client := http.DefaultClient
	values := url.Values{}
	values.Set("client_id", ga.config.ClientID)
	values.Set("client_secret", ga.config.ClientSecret)
	values.Set("refresh_token", refreshToken)
	values.Set("grant_type", "refresh_token")

	res, err := client.PostForm(ga.config.Endpoint.TokenURL, values)
	if err != nil {
		return utils.Empty, 0, utils.Error("token refresh failed", err)
	}

	if res.StatusCode != http.StatusOK {
		msg := fmt.Sprintf("token refresh failed with '%d' status code", res.StatusCode)
		return utils.Empty, 0, errors.New(msg)
	}

	type body struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	var b body
	if err := json.NewDecoder(res.Body).Decode(&b); err != nil {
		msg := "error on response unmarshalling"
		return utils.Empty, 0, utils.Error(msg, err)
	}

	return b.AccessToken, b.ExpiresIn, nil
}

func toRevokeURL(tokenURL string) string {
	u, _ := url.Parse(tokenURL)
	return fmt.Sprintf("%s://%s/revoke", u.Scheme, u.Host)
}
