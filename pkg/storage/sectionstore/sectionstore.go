package sectionstore

import (
	"github.com/psewda/typing/internal/utils"
)

// Sectionstore is the base interface having all operations on section.
type Sectionstore interface {
	// Create adds a new section in the note.
	Create(nid string, s *WritableSection) (*Section, error)

	// GetAll fetches all sections from the note.
	GetAll(nid string) ([]*Section, error)

	// Get returns a single section from the note.
	Get(nid, sid string) (*Section, error)

	// Update modifies the section and saves it back in the note.
	Update(nid, sid string, s *WritableSection) (*Section, error)

	// Delete removes the section from note.
	Delete(nid, sid string) error
}

// WritableSection is used for creating and updating section.
type WritableSection struct {
	Name     string            `json:"name,omitempty" validate:"required,notblank,max=100"`
	Labels   []string          `json:"labels,omitempty" validate:"max=5,dive,max=20"`
	Metadata map[string]string `json:"metadata,omitempty" validate:"max=20,dive,keys,max=20,endkeys,max=100"`
	Data     map[string]string `json:"data,omitempty" validate:"max=50,dive,keys,max=50,endkeys,max=2000"`
}

// Validate checks all validation rules on writable section fields. It returns
// error on any validation failure.
func (s *WritableSection) Validate() error {
	return utils.ValidateStruct(s, messages)
}

// Section represents full detail about section.
type Section struct {
	ID       string            `json:"id,omitempty"`
	Name     string            `json:"name,omitempty"`
	Labels   []string          `json:"labels,omitempty"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Data     map[string]string `json:"data,omitempty"`
}

var messages map[string]string

func init() {
	messages = make(map[string]string)
	messages["name.required"] = "name is required field"
	messages["name.notblank"] = "name can't be empty value"
	messages["name.max"] = "name must be less than 100 chars"
	messages["labels.max"] = "label count can't be more than 5"
	messages["labels.item.max"] = "label must be less than 20 chars"
	messages["metadata.max"] = "metadata count can't be more than 20"
	messages["metadata.item.max"] = "metadata key and value must be less than 20 and 100 chars respectively"
	messages["data.max"] = "data count can't be more than 50"
	messages["data.item.max"] = "data key and value must be less than 50 and 2000 chars respectively"
}
