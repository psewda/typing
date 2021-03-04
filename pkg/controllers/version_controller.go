package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/psewda/typing"
	"github.com/psewda/typing/pkg/types"
)

// VersionController represents all operations on version endpoint.
type VersionController struct{}

// AddRoutes configures all routes of version endpoint
// in the 'echo' server runtime.
func (c *VersionController) AddRoutes(e *echo.Echo) {
	if e != nil {
		e.GET("/api/version", GetVersion)
	}
}

// NewVersionController creates a new instance of version controller.
func NewVersionController() *VersionController {
	return &VersionController{}
}

// GetVersion returns the version string to the client.
func GetVersion(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, types.VersionValue{
		Version: typing.GetVersionString(),
	})
}
