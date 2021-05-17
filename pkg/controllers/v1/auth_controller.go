package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/di"
	"github.com/psewda/typing/pkg/middlewares"
	"github.com/psewda/typing/pkg/signin/auth"
	"github.com/psewda/typing/pkg/types"
)

// AuthController represents all operations on auth endpoint.
type AuthController struct {
	container di.Container
}

// AddRoutes configures all routes of auth endpoint
// in the 'echo' server runtime.
func (c *AuthController) AddRoutes(e *echo.Echo) {
	if e != nil {
		d := middlewares.Dependencies(c.container)
		group := e.Group("/api/v1/signin/auth", d)
		group.GET("/url", GetURL)
		group.POST("/token", Exchange)
		group.POST("/refresh", Refresh)
		group.POST("/revoke", Revoke)
	}
}

// NewAuthController creates a new instance of auth controller.
func NewAuthController(c di.Container) *AuthController {
	return &AuthController{
		container: c,
	}
}

// GetURL returns the authorization workflow url to the client.
func GetURL(ctx echo.Context) error {
	auth := getAuth(ctx)
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
func Exchange(ctx echo.Context) error {
	auth := getAuth(ctx)
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
func Refresh(ctx echo.Context) error {
	auth := getAuth(ctx)
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
func Revoke(ctx echo.Context) error {
	auth := getAuth(ctx)
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

func getAuth(ctx echo.Context) auth.Auth {
	container := ctx.Get(middlewares.KeyContainer).(di.Container)
	instance, _ := container.GetInstance(di.InstanceTypeAuth)
	return instance.(auth.Auth)

}
