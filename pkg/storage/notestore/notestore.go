package notestore

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/psewda/typing/internal/utils"
)

// Notestore is the base interface having all operations on note.
type Notestore interface {
	// Create builds a new note and saves it on cloud storage.
	Create(n *WritableNote) (*Note, error)

	// GetAll fetches all notes from cloud storage.
	GetAll() ([]*Note, error)

	// Get returns the single note from cloud storage.
	Get(id string) (*Note, error)

	// Update modifies the note and saves it on from cloud storage.
	Update(id string, n *WritableNote) (*Note, error)

	// Delete removes the note from cloud storage.
	Delete(id string) (bool, error)
}

// WritableNote is used for creating and updating note.
type WritableNote struct {
	Name        string            `json:"name,omitempty" validate:"required,notblank,max=100"`
	Description string            `json:"desc,omitempty" validate:"max=250"`
	Labels      []string          `json:"labels,omitempty" validate:"max=5,dive,max=20"`
	Metadata    map[string]string `json:"metadata,omitempty" validate:"max=20,dive,keys,max=20,endkeys,max=100"`
}

// Validate checks all validation rules on writable note fields. It returns
// error on any validation failure.
func (n *WritableNote) Validate() error {
	val := validator.New()
	_ = val.RegisterValidation("notblank", validators.NotBlank)
	err := val.Struct(n)
	if err != nil {
		errs := err.(validator.ValidationErrors)
		return errors.New(translate(errs))
	}
	return nil
}

// Note represent full detail about note.
type Note struct {
	ID          string            `json:"id,omitempty"`
	Name        string            `json:"name,omitempty"`
	Description string            `json:"desc,omitempty"`
	Labels      []string          `json:"labels,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	DateCreated time.Time         `json:"dateCreated,omitempty"`
	DateUpdated time.Time         `json:"dateUpdated,omitempty"`
}

var messages map[string]string

func init() {
	messages = make(map[string]string)
	messages["name.required"] = "name is required field"
	messages["name.notblank"] = "name can't be empty value"
	messages["name.max"] = "name must be less than 100 chars"
	messages["description.max"] = "desc must be less than 250 chars"
	messages["labels.max"] = "label count can't be more than 5"
	messages["labels.item.max"] = "label must be less than 20 chars"
	messages["metadata.max"] = "metadata count can't be more than 20"
	messages["metadata.item.max"] = "metadata key and value must be less than 20 and 100 chars respectively"
}

func translate(errs validator.ValidationErrors) string {
	getKey := func(e validator.FieldError) string {
		field := strings.ToLower(e.Field())
		tag := strings.ToLower(e.Tag())
		if strings.Contains(field, "[") {
			index := strings.Index(field, "[")
			return fmt.Sprintf("%s.item.%s", field[0:index], tag)
		}
		return fmt.Sprintf("%s.%s", field, tag)
	}

	var msgs []string
	for _, e := range errs {
		key := getKey(e)
		defmsg := fmt.Sprintf("validation for '%s' failed", key)
		msg := utils.GetValueString(messages[key], defmsg)
		msgs = append(msgs, msg)

	}
	return strings.Join(msgs, ", ")
}
