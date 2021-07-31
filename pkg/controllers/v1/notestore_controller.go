package v1

import (
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/ioc"
	"github.com/psewda/typing/pkg/middlewares"
	"github.com/psewda/typing/pkg/storage/notestore"
)

// NotestoreController represents all operations on notestore endpoint.
type NotestoreController struct {
	container ioc.Container
}

// AddRoutes configures all routes of notestore endpoint
// in the 'echo' server runtime.
func (c *NotestoreController) AddRoutes(e *echo.Echo) {
	if e != nil {
		a := middlewares.Authorization()
		group := e.Group("/api/v1/storage/notes", a)
		group.POST(utils.Empty, c.CreateNote)
		group.GET(utils.Empty, c.GetNotes)
		group.GET("/:id", c.GetNote)
		group.PUT("/:id", c.UpdateNote)
		group.DELETE("/:id", c.DeleteNote)
	}
}

// CreateNote builds a new note and returns to the client.
func (c *NotestoreController) CreateNote(ctx echo.Context) error {
	ns := c.getNotestore(ctx)
	n := new(notestore.WritableNote)

	if err := ctx.Bind(n); err != nil {
		msg := "spec validation failed"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	// check all rules on note validation
	err := n.Validate()
	if err != nil {
		msg := err.Error()
		ctx.Logger().Warn(msg)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	note, err := ns.Create(n)
	if err != nil {
		msg := "note creation error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return utils.BuildHTTPError(err, msg)
	}

	ctx.Response().Header().Add(echo.HeaderLocation, path.Join(ctx.Path(), note.ID))
	return ctx.JSON(http.StatusCreated, note)
}

// GetNotes fetches all notes from the cloud storage and return to the client.
func (c *NotestoreController) GetNotes(ctx echo.Context) error {
	ns := c.getNotestore(ctx)

	notes, err := ns.GetAll()
	if err != nil {
		msg := "note retrival error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return utils.BuildHTTPError(err, msg)
	}

	return ctx.JSON(http.StatusOK, notes)
}

// GetNote fetches the single note from the cloud storage and return to the client.
func (c *NotestoreController) GetNote(ctx echo.Context) error {
	ns := c.getNotestore(ctx)
	id := ctx.Param("id")

	note, err := ns.Get(id)
	if err != nil {
		msg := "note retrival error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return utils.BuildHTTPError(err, msg)
	}

	return ctx.JSON(http.StatusOK, note)
}

// UpdateNote modifies the note and saves back on cloud storage.
func (c *NotestoreController) UpdateNote(ctx echo.Context) error {
	ns := c.getNotestore(ctx)
	id := ctx.Param("id")
	n := new(notestore.WritableNote)

	if err := ctx.Bind(n); err != nil {
		msg := "spec validation failed"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	// check all rules on note validation
	err := n.Validate()
	if err != nil {
		msg := err.Error()
		ctx.Logger().Warn(msg)
		return &echo.HTTPError{
			Code:    http.StatusBadRequest,
			Message: msg,
		}
	}

	// try to update the note
	note, err := ns.Update(id, n)
	if err != nil {
		msg := "note updation error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return utils.BuildHTTPError(err, msg)
	}

	return ctx.JSON(http.StatusOK, note)
}

// DeleteNote removes the note from cloud storage.
func (c *NotestoreController) DeleteNote(ctx echo.Context) error {
	ns := c.getNotestore(ctx)
	id := ctx.Param("id")

	err := ns.Delete(id)
	if err != nil {
		msg := "note deletion error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return utils.BuildHTTPError(err, msg)
	}

	return ctx.NoContent(http.StatusNoContent)
}

// NewNotestoreController creates a new instance of notestore controller.
func NewNotestoreController(c ioc.Container) *NotestoreController {
	return &NotestoreController{
		container: c,
	}
}

func (c *NotestoreController) getNotestore(ctx echo.Context) notestore.Notestore {
	accessToken := ctx.Get(middlewares.ContextKeyAccessToken).(string)
	client := utils.ClientWithToken(accessToken)
	instance, _ := c.container.GetInstance(ioc.InstanceTypeNotestore, client)
	return instance.(notestore.Notestore)
}
