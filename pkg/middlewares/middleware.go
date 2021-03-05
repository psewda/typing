package middlewares

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/di"
)

const (
	// KeyContainer is the key for container object.
	KeyContainer = "#container"

	// KeyAccessToken is the key for access token value.
	KeyAccessToken = "#accesstoken"
)

// Authorization middleware authorizes http request by validating
// bearer token in authorization header. If validation is
// successful, it inserts the access token in the echo context.
func Authorization() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			value := ctx.Request().Header.Get(echo.HeaderAuthorization)
			if len(value) == 0 {
				msg := "authorization header is empty, set valid authorization token"
				ctx.Logger().Warn(msg)
				return &echo.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: msg,
				}
			}
			t := fetchToken(value)
			if t == utils.Empty {
				msg := "authorization token is in invalid format"
				ctx.Logger().Warn(msg)
				return &echo.HTTPError{
					Code:    http.StatusUnauthorized,
					Message: msg,
				}
			}

			ctx.Set(KeyAccessToken, t)
			return next(ctx)
		}
	}
}

// Dependencies middleware inserts the container in the echo ontext.
func Dependencies(c di.Container) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set(KeyContainer, c)
			return next(ctx)
		}
	}
}

func fetchToken(value string) string {
	const scheme = "Bearer"
	if strings.HasPrefix(value, scheme) {
		return strings.TrimSpace(value[len(scheme):])
	}
	return utils.Empty
}
