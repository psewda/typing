package drvnotestore

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/storage/notestore"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

const (
	appdir         = "appDataFolder"
	fileFields     = "id, name, description, properties, createdTime, modifiedTime"
	fileListFields = "files(id, name, description, properties, createdTime, modifiedTime)"
)

// DrvNotestore is the notestore implementation
// using google drive api.
type DrvNotestore struct {
	service *drive.Service
}

// Create builds a new note and saves it on google drive.
func (ns *DrvNotestore) Create(n *notestore.WritableNote) (*notestore.Note, error) {
	err := checkNote(n)
	if err != nil {
		return nil, utils.Error("note validation failed", err)
	}

	note := sanitize(n)
	f := drive.File{
		Name:        fmt.Sprintf("%s.json", note.Name),
		Description: note.Description,
		MimeType:    "application/json",
		Parents:     []string{appdir},
		Properties:  fillProps(note),
	}

	file, err := ns.service.Files.Create(&f).Fields(fileFields).Do()
	if err != nil {
		return nil, utils.Error("file creation error", err)
	}

	return toNote(file), nil
}

// GetAll returns a list of all notes from google drive
func (ns *DrvNotestore) GetAll() ([]*notestore.Note, error) {
	list, err := ns.service.Files.List().Spaces(appdir).Fields(fileListFields).Do()
	if err != nil {
		return nil, utils.Error("file listing error", err)
	}

	var notes []*notestore.Note
	for _, f := range list.Files {
		notes = append(notes, toNote(f))
	}
	return notes, nil
}

// Get returns the single note from google drive.
func (ns *DrvNotestore) Get(id string) (*notestore.Note, error) {
	if len(id) == 0 {
		return nil, errors.New("note id is nil")
	}

	file, err := getFile(id, ns.service)
	if err != nil {
		return nil, err
	}

	// file found, so convert it into note
	if file != nil {
		return toNote(file), nil
	}

	// file doesn't exist, so return nil
	return nil, nil
}

// Update modifies the note and saves back on google drive.
func (ns *DrvNotestore) Update(id string, n *notestore.WritableNote) (*notestore.Note, error) {
	if len(id) == 0 {
		return nil, errors.New("note id is nil")
	}

	file, err := getFile(id, ns.service)
	if err != nil {
		return nil, err
	}

	// file found, so try to update it
	if file != nil {
		note := sanitize(n)

		// update api uses patch sementics, so to clear a field, it must be sent as
		// null value. If any field has an empty value, then it will be ignored
		// in json serialization and will NOT be updated. We have to use NullFields
		// and ForceSendFields to tell which fields are required to be sent as
		// null value. The below link has more detail on the issue -
		// https://github.com/googleapis/google-api-go-client/issues/201
		f := drive.File{
			Name:            fmt.Sprintf("%s.json", note.Name),
			Description:     note.Description,
			Properties:      fillProps(note),
			NullFields:      getNullFields(note, file),
			ForceSendFields: []string{"Description", "Properties"},
		}

		updated, err := ns.service.Files.Update(file.Id, &f).Fields(fileFields).Do()
		if err != nil {
			return nil, utils.Error("file updation error", err)
		}
		return toNote(updated), nil
	}

	// file doesn't exist, so return nil
	return nil, nil
}

// Delete removes the note from google drive.
func (ns *DrvNotestore) Delete(id string) (bool, error) {
	if len(id) == 0 {
		return false, errors.New("note id is nil")
	}

	err := ns.service.Files.Delete(id).Do()
	if err != nil {
		// file doesn't exist, so return false
		if checkNotFound(err) {
			return false, nil
		}
		return false, utils.Error("file deletion error", err)
	}

	// file deleted, so return true
	return true, nil
}

// New creates a new instance of google drive notestore.
func New(c *http.Client) (*DrvNotestore, error) {
	service, err := drive.New(c)
	if err != nil {
		return nil, utils.Error("drive service creation error", err)
	}

	return &DrvNotestore{
		service: service,
	}, nil
}

func parseTime(value string) time.Time {
	if len(value) > 0 {
		t, err := time.Parse(time.RFC3339, value)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}

func sanitize(n *notestore.WritableNote) *notestore.WritableNote {
	note := notestore.WritableNote{
		Name:        strings.TrimSpace(n.Name),
		Description: strings.TrimSpace(n.Description),
	}

	for _, l := range n.Labels {
		cleanLabel := strings.TrimSpace(l)
		if len(cleanLabel) > 0 {
			note.Labels = append(note.Labels, cleanLabel)
		}
	}
	note.Metadata = utils.Sanitize(n.Metadata)
	return &note
}

func fillProps(n *notestore.WritableNote) map[string]string {
	props := make(map[string]string)

	if len(n.Labels) > 0 {
		labels := strings.Join(n.Labels, ",")
		props["labels"] = labels
	}

	if len(n.Metadata) > 0 {
		for k, v := range n.Metadata {
			props[fmt.Sprintf("meta!%s", k)] = v
		}
	}

	return props
}

func toNote(f *drive.File) *notestore.Note {
	n := notestore.Note{
		ID:          f.Id,
		Name:        strings.TrimRight(f.Name, ".json"),
		Description: f.Description,
		DateCreated: parseTime(f.CreatedTime),
		DateUpdated: parseTime(f.ModifiedTime),
	}

	if len(f.Properties["labels"]) > 0 {
		n.Labels = strings.Split(f.Properties["labels"], ",")
	}
	for k, v := range f.Properties {
		if strings.HasPrefix(k, "meta!") {
			if n.Metadata == nil {
				n.Metadata = make(map[string]string)
			}
			n.Metadata[k[5:]] = v
		}
	}

	return &n
}

func checkNote(n *notestore.WritableNote) error {
	if n == nil {
		return errors.New("note is nil")
	}
	return n.Validate()
}

func getFile(id string, service *drive.Service) (*drive.File, error) {
	file, err := service.Files.Get(id).Fields(fileFields).Do()
	if err != nil {
		if checkNotFound(err) {
			return nil, nil
		}
		return nil, utils.Error("file retrival error", err)
	}

	return file, nil
}

func checkNotFound(err error) bool {
	e, ok := err.(*googleapi.Error)
	if ok {
		if e.Code == http.StatusNotFound {
			return true
		}
	}
	return false
}

func getNullFields(n *notestore.WritableNote, f *drive.File) []string {
	nullFields := make([]string, 0)
	if len(n.Labels) == 0 {
		nullFields = append(nullFields, "Properties.labels")
	}

	for k := range f.Properties {
		if strings.HasPrefix(k, "meta!") {
			_, ok := n.Metadata[k[5:]]
			if !ok {
				nullFields = append(nullFields, fmt.Sprintf("Properties.%s", k))
			}
		}
	}
	return nullFields
}
