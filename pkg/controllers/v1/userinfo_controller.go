package v1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/ioc"
	"github.com/psewda/typing/pkg/middlewares"
	"github.com/psewda/typing/pkg/signin/userinfo"
)

// UserinfoController represents all operations on userinfo endpoint.
type UserinfoController struct {
	container ioc.Container
}

// AddRoutes configures all routes of userinfo endpoint
// in the 'echo' server runtime.
func (c *UserinfoController) AddRoutes(e *echo.Echo) {
	if e != nil {
		a := middlewares.Authorization()
		e.GET("/api/v1/signin/userinfo", c.GetUser, a)
	}
}

// GetUser returns the user detail to the client.
func (c *UserinfoController) GetUser(ctx echo.Context) error {
	ui := c.getUserinfo(ctx)
	u, err := ui.Get()
	if err != nil {
		msg := "error occurred while fetching user"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return utils.BuildHTTPError(err, msg)
	}

	return ctx.JSON(http.StatusOK, u)
}

// NewUserinfoController creates a new instance of userinfo controller.
func NewUserinfoController(c ioc.Container) *UserinfoController {
	return &UserinfoController{
		container: c,
	}
}

func (c *UserinfoController) getUserinfo(ctx echo.Context) userinfo.Userinfo {
	accessToken := ctx.Get(middlewares.ContextKeyAccessToken).(string)
	client := utils.ClientWithToken(accessToken)
	instance, _ := c.container.GetInstance(ioc.InstanceTypeUserinfo, client)
	return instance.(userinfo.Userinfo)
}
