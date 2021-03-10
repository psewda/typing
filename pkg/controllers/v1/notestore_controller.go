package v1

import (
	"fmt"
	"net/http"
	"path"

	"github.com/labstack/echo/v4"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/di"
	"github.com/psewda/typing/pkg/middlewares"
	"github.com/psewda/typing/pkg/storage/notestore"
)

// NotestoreController represents all operations on notestore endpoint.
type NotestoreController struct {
	container di.Container
}

// AddRoutes configures all routes of notestore endpoint
// in the 'echo' server runtime.
func (c *NotestoreController) AddRoutes(e *echo.Echo) {
	if e != nil {
		a := middlewares.Authorization()
		d := middlewares.Dependencies(c.container)
		group := e.Group("/api/v1/storage/notes", a, d)
		group.POST(utils.Empty, CreateNote)
		group.GET(utils.Empty, GetNotes)
		group.GET("/:id", GetNote)
		group.PUT("/:id", UpdateNote)
		group.DELETE("/:id", DeleteNote)
	}
}

// NewNotestoreController creates a new instance of notestore controller.
func NewNotestoreController(c di.Container) *NotestoreController {
	return &NotestoreController{
		container: c,
	}
}

// CreateNote builds a new note and returns to the client.
func CreateNote(ctx echo.Context) error {
	ns := getNotestore(ctx)
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
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: msg,
		}
	}

	ctx.Response().Header().Add(echo.HeaderLocation, path.Join(ctx.Path(), note.ID))
	return ctx.JSON(http.StatusCreated, note)
}

// GetNotes fetches all notes from the cloud storage and return to the client.
func GetNotes(ctx echo.Context) error {
	ns := getNotestore(ctx)

	notes, err := ns.GetAll()
	if err != nil {
		msg := "note retrival error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: msg,
		}
	}

	return ctx.JSON(http.StatusOK, notes)
}

// GetNote fetches the single note from the cloud storage and return to the client.
func GetNote(ctx echo.Context) error {
	ns := getNotestore(ctx)
	id := ctx.Param("id")

	note, err := ns.Get(id)
	if err != nil {
		msg := "note retrival error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: msg,
		}
	}

	// note doesn't exist, so return 'NotFound' status
	if note == nil {
		msg := fmt.Sprintf("note with id '%s' not found", id)
		ctx.Logger().Warn(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: msg,
		}
	}
	return ctx.JSON(http.StatusOK, note)
}

// UpdateNote modifies the note and saves back on cloud storage.
func UpdateNote(ctx echo.Context) error {
	ns := getNotestore(ctx)
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
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: msg,
		}
	}

	// note doesn't exist, so return 'NotFound' status
	if note == nil {
		msg := fmt.Sprintf("note with id '%s' not found", id)
		ctx.Logger().Warn(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: msg,
		}
	}

	return ctx.JSON(http.StatusOK, note)
}

// DeleteNote removes the note from cloud storage.
func DeleteNote(ctx echo.Context) error {
	ns := getNotestore(ctx)
	id := ctx.Param("id")

	status, err := ns.Delete(id)
	if err != nil {
		msg := "note deletion error"
		ctx.Logger().Error(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusInternalServerError,
			Message: msg,
		}
	}

	// note doesn't exist, so return 'NotFound' status
	if !status {
		msg := fmt.Sprintf("note with id '%s' not found", id)
		ctx.Logger().Warn(utils.AppendError(msg, err))
		return &echo.HTTPError{
			Code:    http.StatusNotFound,
			Message: msg,
		}
	}

	return ctx.NoContent(http.StatusNoContent)
}

func getNotestore(ctx echo.Context) notestore.Notestore {
	accessToken := ctx.Get(middlewares.KeyAccessToken).(string)
	client := utils.ClientWithToken(accessToken)
	container := ctx.Get(middlewares.KeyContainer).(di.Container)
	instance, _ := container.GetInstance(di.InstanceTypeNotestore, client)
	return instance.(notestore.Notestore)
}
