package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/di"
	"github.com/psewda/typing/pkg/middlewares"
	"github.com/psewda/typing/pkg/signin/userinfo"
)

// UserinfoController represents all operations on userinfo endpoint.
type UserinfoController struct {
	container di.Container
}

// AddRoutes configures all routes of userinfo endpoint
// in the 'echo' server runtime.
func (c *UserinfoController) AddRoutes(e *echo.Echo) {
	if e != nil {
		a := middlewares.Authorization()
		d := middlewares.Dependencies(c.container)
		e.GET("/api/v1/signin/userinfo", GetUser, a, d)
	}
}

// NewUserinfoController creates a new instance of userinfo controller.
func NewUserinfoController(c di.Container) *UserinfoController {
	return &UserinfoController{
		container: c,
	}
}

// GetUser returns the user detail to the client.
func GetUser(ctx echo.Context) error {
	ui := getUserinfo(ctx)
	u, err := ui.Get()
	if err != nil {
		msg := "error occurred while fetching user"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: msg,
		}
	}

	return ctx.JSON(http.StatusOK, u)
}

func getUserinfo(ctx echo.Context) userinfo.Userinfo {
	accessToken := ctx.Get(middlewares.KeyAccessToken).(string)
	client := utils.ClientWithToken(accessToken)
	container := ctx.Get(middlewares.KeyContainer).(di.Container)
	instance, _ := container.GetInstance(di.InstanceTypeUserinfo, client)
	return instance.(userinfo.Userinfo)
}
