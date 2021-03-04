package controllers

import (
	"github.com/labstack/echo/v4"
)

// Controller is the base type for all specific controllers.
type Controller interface {
	AddRoutes(e *echo.Echo)
}
