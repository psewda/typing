package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/ioc"
	"github.com/psewda/typing/pkg/signin/auth"
	"github.com/psewda/typing/pkg/types"
)

// AuthController represents all operations on auth endpoint.
type AuthController struct {
	container ioc.Container
}

// AddRoutes configures all routes of auth endpoint
// in the 'echo' server runtime.
func (c *AuthController) AddRoutes(e *echo.Echo) {
	if e != nil {
		group := e.Group("/api/v1/signin/auth")
		group.GET("/url", c.GetURL)
		group.POST("/token", c.Exchange)
		group.POST("/refresh", c.Refresh)
		group.POST("/revoke", c.Revoke)
	}
}

// GetURL returns the authorization workflow url to the client.
func (c *AuthController) GetURL(ctx echo.Context) error {
	auth := c.getAuth()
	redirect := ctx.QueryParam("redirect")
	state := ctx.QueryParam("state")

	if len(redirect) > 0 {
		if err := utils.CheckLocalhostURL(redirect); err != nil {
			msg := "redirect url is invalid or is not a localhost url"
			ctx.Logger().Warn(utils.AppendError(msg, err))
			return &echo.HTTPError{
				Code:    http.StatusBadRequest,
				Message: msg,
			}
		}
	}

	url := auth.GetURL(redirect, state)
	return ctx.JSON(http.StatusOK, types.URLValue{
		URL: url,
	})
}

// Exchange converts the authorization code into token and
// returns it to the client.
func (c *AuthController) Exchange(ctx echo.Context) error {
	auth := c.getAuth()
	code := ctx.FormValue("auth_code")
	if len(code) == 0 {
		msg := "authorization code is empty"
		ctx.Logger().Warn(msg)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	token, err := auth.Exchange(code)
	if err != nil {
		msg := "token exchange failed, check the authorization code"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: msg,
		}
	}

	return ctx.JSON(http.StatusOK, token)
}

// Refresh renews access token using refresh token.
func (c *AuthController) Refresh(ctx echo.Context) error {
	auth := c.getAuth()
	refreshToken := ctx.FormValue("refresh_token")
	if len(refreshToken) == 0 {
		msg := "refresh token is empty"
		ctx.Logger().Warn(msg)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	token, err := auth.Refresh(refreshToken)
	if err != nil {
		msg := "access token refresh failed, check the token"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: msg,
		}
	}

	return ctx.JSON(http.StatusOK, token)
}

// Revoke resets the authorization workflow. After calling revoke, the
// user needs to start authorization workflow again.
func (c *AuthController) Revoke(ctx echo.Context) error {
	auth := c.getAuth()
	token := ctx.FormValue("token")
	if len(token) == 0 {
		msg := "token value is empty"
		ctx.Logger().Warn(msg)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	err := auth.Revoke(token)
	if err != nil {
		msg := "token revocation failed, check the token value"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: msg,
		}
	}

	return ctx.NoContent(http.StatusNoContent)
}

// NewAuthController creates a new instance of auth controller.
func NewAuthController(c ioc.Container) *AuthController {
	return &AuthController{
		container: c,
	}
}

func (c *AuthController) getAuth() auth.Auth {
	instance, _ := c.container.GetInstance(ioc.InstanceTypeAuth)
	return instance.(auth.Auth)
}
