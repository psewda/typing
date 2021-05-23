package v1

import (
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/di"
	"github.com/psewda/typing/pkg/middlewares"
	"github.com/psewda/typing/pkg/storage/sectionstore"
)

// SectionstoreController represents all operations on sectionstore endpoint.
type SectionstoreController struct {
	container di.Container
}

// AddRoutes configures all routes of sectionstore endpoint
// in the 'echo' server runtime.
func (c *SectionstoreController) AddRoutes(e *echo.Echo) {
	if e != nil {
		a := middlewares.Authorization()
		d := middlewares.Dependencies(c.container)
		group := e.Group("/api/v1/storage/notes/:nid/sections", a, d)
		group.POST(utils.Empty, CreateSection)
		group.GET(utils.Empty, GetSections)
		group.GET("/:id", GetSection)
		group.PUT("/:id", UpdateSection)
		group.DELETE("/:id", DeleteSection)
	}
}

// NewSectionstoreController creates a new instance of sectionstore controller.
func NewSectionstoreController(c di.Container) *SectionstoreController {
	return &SectionstoreController{
		container: c,
	}
}

// CreateSection adds a new section in the note and returns to the client.
func CreateSection(ctx echo.Context) error {
	ss := getSectionstore(ctx)
	nid := ctx.Param("nid")
	s := new(sectionstore.WritableSection)

	if err := ctx.Bind(s); err != nil {
		msg := "spec validation failed"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	// check all rules on section validation
	err := s.Validate()
	if err != nil {
		msg := err.Error()
		ctx.Logger().Warn(msg)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	section, err := ss.Create(nid, s)
	if err != nil {
		msg := "section creation error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return utils.BuildHTTPError(err, msg)
	}

	location := path.Join(ctx.Request().URL.Path, section.ID)
	ctx.Response().Header().Add(echo.HeaderLocation, location)
	return ctx.JSON(http.StatusCreated, section)
}

// GetSections fetches all sections from the note and returns to the client.
func GetSections(ctx echo.Context) error {
	ss := getSectionstore(ctx)
	nid := ctx.Param("nid")

	sections, err := ss.GetAll(nid)
	if err != nil {
		msg := "section retrival error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return utils.BuildHTTPError(err, msg)
	}

	return ctx.JSON(http.StatusOK, sections)
}

// GetSection fetches the single section from the note and returns to the client.
func GetSection(ctx echo.Context) error {
	ss := getSectionstore(ctx)
	nid := ctx.Param("nid")
	id := ctx.Param("id")

	section, err := ss.Get(nid, id)
	if err != nil {
		msg := "section retrival error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return utils.BuildHTTPError(err, msg)
	}

	return ctx.JSON(http.StatusOK, section)
}

// UpdateSection modifies the section and saves it back in the note.
func UpdateSection(ctx echo.Context) error {
	ss := getSectionstore(ctx)
	nid := ctx.Param("nid")
	id := ctx.Param("id")
	s := new(sectionstore.WritableSection)

	if err := ctx.Bind(s); err != nil {
		msg := "spec validation failed"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	// check all rules on section validation
	err := s.Validate()
	if err != nil {
		msg := err.Error()
		ctx.Logger().Warn(msg)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	// try to update the section
	section, err := ss.Update(nid, id, s)
	if err != nil {
		msg := "section updation error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return utils.BuildHTTPError(err, msg)
	}

	return ctx.JSON(http.StatusOK, section)
}

// DeleteSection removes the section from note.
func DeleteSection(ctx echo.Context) error {
	ss := getSectionstore(ctx)
	nid := ctx.Param("nid")
	id := ctx.Param("id")

	err := ss.Delete(nid, id)
	if err != nil {
		msg := "section deletion error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return utils.BuildHTTPError(err, msg)
	}

	return ctx.NoContent(http.StatusNoContent)
}

func getSectionstore(ctx echo.Context) sectionstore.Sectionstore {
	accessToken := ctx.Get(middlewares.KeyAccessToken).(string)
	client := utils.ClientWithToken(accessToken)
	container := ctx.Get(middlewares.KeyContainer).(di.Container)
	instance, _ := container.GetInstance(di.InstanceTypeSectionstore, client)
	return instance.(sectionstore.Sectionstore)
}
