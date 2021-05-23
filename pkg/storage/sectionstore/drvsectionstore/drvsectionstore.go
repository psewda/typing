package drvsectionstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/errs"
	secstore "github.com/psewda/typing/pkg/storage/sectionstore"
	"github.com/rs/xid"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
)

// DrvSectionstore is the sectionstore implementation
// using google drive api.
type DrvSectionstore struct {
	service *drive.Service
}

// Create adds a new section in the note and stores the data on google drive.
func (ss *DrvSectionstore) Create(nid string, s *secstore.WritableSection) (*secstore.Section, error) {
	err := checkSection(s)
	if err != nil {
		return nil, utils.Error("section validation failed", err)
	}

	// download existing note content
	content, err := download(ss.service, nid)
	if err != nil {
		return nil, err
	}

	// build new section instance
	sanitized := sanitize(s)
	section := &secstore.Section{
		ID:       xid.New().String(),
		Name:     sanitized.Name,
		Labels:   sanitized.Labels,
		Metadata: sanitized.Metadata,
		Data:     sanitized.Data,
	}

	// append the new section
	var sections []*secstore.Section
	if len(content) > 0 {
		sections, err = unmarshal(content)
		if err != nil {
			return nil, err
		}
	}
	sections = append(sections, section)

	// upload the note content
	j, _ := json.Marshal(sections)
	if err := upload(ss.service, nid, j); err != nil {
		return nil, err
	}
	return section, nil
}

// GetAll fetches all sections from the note.
func (ss *DrvSectionstore) GetAll(nid string) ([]*secstore.Section, error) {
	content, err := download(ss.service, nid)
	if err != nil {
		return nil, err
	}

	var sections []*secstore.Section
	if len(content) > 0 {
		sections, err = unmarshal(content)
		if err != nil {
			return nil, err
		}
	}
	return sections, nil
}

// Get returns a single section from the note.
func (ss *DrvSectionstore) Get(nid, sid string) (*secstore.Section, error) {
	content, err := download(ss.service, nid)
	if err != nil {
		return nil, err
	}

	var sections []*secstore.Section
	if len(content) > 0 {
		sections, err = unmarshal(content)
		if err != nil {
			return nil, err
		}
	}

	// find the section in the array
	idx := indexOf(sections, sid)
	if idx == -1 {
		msg := fmt.Sprintf("section with id '%s' not found", sid)
		return nil, errs.NewNotFoundError(msg)
	}

	// section found, so return the section
	return sections[idx], nil
}

// Update modifies the section and saves it back in the note.
func (ss *DrvSectionstore) Update(nid, sid string, s *secstore.WritableSection) (*secstore.Section, error) {
	err := checkSection(s)
	if err != nil {
		return nil, utils.Error("section validation failed", err)
	}

	// download existing note content
	content, err := download(ss.service, nid)
	if err != nil {
		return nil, err
	}

	var sections []*secstore.Section
	if len(content) > 0 {
		sections, err = unmarshal(content)
		if err != nil {
			return nil, err
		}
	}

	// find the section in the array
	idx := indexOf(sections, sid)
	if idx == -1 {
		msg := fmt.Sprintf("section with id '%s' not found", sid)
		return nil, errs.NewNotFoundError(msg)
	}

	// update section fields
	sanitized := sanitize(s)
	sections[idx].Name = sanitized.Name
	sections[idx].Labels = sanitized.Labels
	sections[idx].Metadata = sanitized.Metadata
	sections[idx].Data = sanitized.Data

	// upload the note content
	j, _ := json.Marshal(sections)
	if err := upload(ss.service, nid, j); err != nil {
		return nil, err
	}
	return sections[idx], nil
}

// Delete removes the section from note.
func (ss *DrvSectionstore) Delete(nid, sid string) error {
	content, err := download(ss.service, nid)
	if err != nil {
		return err
	}

	var sections []*secstore.Section
	if len(content) > 0 {
		sections, err = unmarshal(content)
		if err != nil {
			return err
		}
	}

	// find the section in the array
	idx := indexOf(sections, sid)
	if idx == -1 {
		msg := fmt.Sprintf("section with id '%s' not found", sid)
		return errs.NewNotFoundError(msg)
	}

	// delete the section
	sections[idx] = sections[len(sections)-1]
	sections[len(sections)-1] = nil
	sections = sections[:len(sections)-1]

	// upload the note content
	j, _ := json.Marshal(sections)
	if err := upload(ss.service, nid, j); err != nil {
		return err
	}
	return nil
}

// New creates a new instance of google drive sectionstore.
func New(c *http.Client) (*DrvSectionstore, error) {
	service, err := drive.New(c)
	if err != nil {
		return nil, utils.Error("drive service creation error", err)
	}

	return &DrvSectionstore{
		service: service,
	}, nil
}

func checkSection(s *secstore.WritableSection) error {
	if s == nil {
		return errors.New("section is nil")
	}
	return s.Validate()
}

func sanitize(s *secstore.WritableSection) *secstore.WritableSection {
	section := secstore.WritableSection{
		Name: strings.TrimSpace(s.Name),
	}

	for _, l := range s.Labels {
		cleanLabel := strings.TrimSpace(l)
		if len(cleanLabel) > 0 {
			section.Labels = append(section.Labels, cleanLabel)
		}
	}
	section.Metadata = utils.Sanitize(s.Metadata)
	section.Data = utils.Sanitize(s.Data)

	return &section
}

func download(ds *drive.Service, nid string) ([]byte, error) {
	res, err := ds.Files.Get(nid).Download()
	if err != nil {
		if utils.GetStatusCode(err) == http.StatusUnauthorized {
			return nil, errs.NewUnauthorizedError()
		}
		if utils.GetStatusCode(err) == http.StatusNotFound {
			msg := fmt.Sprintf("note with id '%s' not found", nid)
			return nil, errs.NewNotFoundError(msg)
		}
		return nil, utils.Error("note download error", err)
	}

	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func upload(ds *drive.Service, nid string, content []byte) error {
	reader := bytes.NewReader(content)
	_, err := ds.Files.Update(nid, nil).Media(reader, googleapi.ContentType("application/json")).Do()

	if err != nil {
		if utils.GetStatusCode(err) == http.StatusUnauthorized {
			return errs.NewUnauthorizedError()
		}
		return utils.Error("note upload error", err)
	}
	return nil
}

func unmarshal(content []byte) ([]*secstore.Section, error) {
	var sections []*secstore.Section
	err := json.Unmarshal(content, &sections)
	if err != nil {
		return nil, utils.Error("error on unmarshalling sections", err)
	}
	return sections, nil
}

func indexOf(sections []*secstore.Section, sid string) int {
	for i, s := range sections {
		if s.ID == sid {
			return i
		}
	}
	return -1
}
