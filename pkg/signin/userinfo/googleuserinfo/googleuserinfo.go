package googleuserinfo

import (
	"context"
	"net/http"

	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/signin/userinfo"
	oauth2v2 "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
)

// GoogleUserinfo is the userinfo implementation
// using google's oauth api.
type GoogleUserinfo struct {
	service *oauth2v2.Service
}

// Get returns the basic detail of user. It
// fetches user detail using google's oauth api.
func (gu *GoogleUserinfo) Get() (*userinfo.User, error) {
	ui, err := gu.service.Userinfo.Get().Do()
	if err != nil {
		return nil, utils.Error("error while getting user info", err)
	}

	return &userinfo.User{
		ID:      ui.Id,
		Name:    ui.Name,
		Email:   ui.Email,
		Picture: ui.Picture,
	}, nil
}

// New creates a new instance of google userinfo.
func New(c *http.Client) (*GoogleUserinfo, error) {
	service, err := oauth2v2.NewService(context.Background(), option.WithHTTPClient(c))
	if err != nil {
		msg := "error while creating new instance of oauth2 service"
		return nil, utils.Error(msg, err)
	}

	return &GoogleUserinfo{
		service: service,
	}, nil
}
