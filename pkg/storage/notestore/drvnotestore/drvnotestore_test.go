package drvnotestore_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/errs"
	"github.com/psewda/typing/pkg/storage/notestore"
	"github.com/psewda/typing/pkg/storage/notestore/drvnotestore"
	"google.golang.org/api/drive/v3"
)

func TestDrvNotestore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "drvnotestore-suite")
}

var _ = Describe("googledrive notestore", func() {
	Context("create new note", func() {
		It("should create new note successfully", func() {
			By("correct input")
			{
				j := `{
					"id": "gdtg45w9mjh10ds",
					"name": "note",
					"description": "desc",
					"createdTime": "2021-02-12T07:20:50.52Z"
				}`
				client := utils.ClientWithJSON(j, http.StatusCreated)
				drvns, _ := drvnotestore.New(client)
				note, err := drvns.Create(&notestore.WritableNote{
					Name:        "note",
					Description: "desc",
					Labels:      []string{"label1", "label2"},
					Metadata:    map[string]string{"key": "value"},
				})

				Expect(err).ShouldNot(HaveOccurred())
				Expect(note).ShouldNot(BeNil())
				Expect(note.Name).Should(Equal("note"))
				Expect(note.DateCreated).ShouldNot(BeZero())
			}

			By("unsanitized input")
			{
				verifyReq := func(req *http.Request) {
					var f drive.File
					json.NewDecoder(req.Body).Decode(&f)
					labels := f.Properties["labels"]
					meta1 := f.Properties["meta!key1"]
					meta2 := f.Properties["meta!key2"]

					Expect(f.Name).Should(Equal("note.json"))
					Expect(f.Description).Should(Equal("desc"))
					Expect(labels).Should(Equal("label1,label2"))
					Expect(meta1).Should(Equal("value1"))
					Expect(meta2).Should(Equal("value2"))

					metaCount := count(f.Properties, func(k, v string) bool {
						return strings.HasPrefix(k, "meta!")
					})
					Expect(metaCount).Should(Equal(2))
				}

				client := http.DefaultClient
				client.Transport = utils.TransportFunc(func(req *http.Request) (*http.Response, error) {
					verifyReq(req)
					j := `{
						"id": "gdtg45w9mjh10ds",
						"name": "note"
					}`
					return &http.Response{
						StatusCode: http.StatusCreated,
						Body:       ioutil.NopCloser(bytes.NewBufferString(j)),
						Header:     map[string][]string{"Content-Type": {"application/json"}},
					}, nil
				})

				drvns, _ := drvnotestore.New(client)
				note, err := drvns.Create(&notestore.WritableNote{
					Name:        "note",
					Description: " desc  ",
					Labels:      []string{"label1", " ", "label2  "},
					Metadata: map[string]string{
						"key1":   "value1",
						"key2  ": "value2   ",
						" ":      "  value",
					},
				})

				Expect(err).ShouldNot(HaveOccurred())
				Expect(note).ShouldNot(BeNil())
				Expect(note.Name).Should(Equal("note"))
			}

			By("wrong input")
			{
				client := utils.ClientWithJSON("{}", http.StatusCreated)
				drvns, _ := drvnotestore.New(client)
				note, err := drvns.Create(&notestore.WritableNote{
					Description: "desc",
				})

				Expect(err).Should(HaveOccurred())
				Expect(note).Should(BeNil())
			}

			By("authorization failure")
			{
				code := http.StatusUnauthorized
				client := utils.ClientWithJSON("{}", code)
				drvns, _ := drvnotestore.New(client)
				_, err := drvns.Create(&notestore.WritableNote{
					Name: "note",
				})

				Expect(err).Should(HaveOccurred())
				Expect(err).Should(BeAssignableToTypeOf(errs.NewUnauthorizedError()))
			}

			By("inner error")
			{
				code := http.StatusInternalServerError
				client := utils.ClientWithJSON("error", code)
				drvns, _ := drvnotestore.New(client)
				note, err := drvns.Create(&notestore.WritableNote{
					Name: "note",
				})

				Expect(err).Should(HaveOccurred())
				Expect(note).Should(BeNil())
			}
		})
	})

	Context("get all notes", func() {
		It("should return a list of all notes", func() {
			By("correct setup")
			{
				j := `{
				"files": [
				  {
					"id": "gdtg45w9mjh10ds",
					"name": "note1",
					"properties": {
					  "labels": "label1,label2",
					  "meta!key1": "value1"
					},
					"createdTime": "2021-02-12T07:20:50.52Z"
				  },
				  {
					"id": "gdtg11w9mjh10ds",
					"name": "note2",
					"properties": {
					  "labels": "label1,label2",
					  "meta!key1": "value1"
					},
					"createdTime": "2021-02-12T08:20:50.52Z"
				  }
				]
			  }`
				client := utils.ClientWithJSON(j, http.StatusOK)
				drvns, _ := drvnotestore.New(client)
				notes, err := drvns.GetAll()

				Expect(err).ShouldNot(HaveOccurred())
				Expect(notes).ShouldNot(BeNil())
				Expect(len(notes)).Should(Equal(2))
				Expect(notes[0].ID).ShouldNot(BeZero())
				Expect(notes[0].Name).Should(Equal("note1"))
				Expect(len(notes[0].Labels)).Should(Equal(2))
				Expect(len(notes[0].Metadata)).Should(Equal(1))
				Expect(notes[1].ID).ShouldNot(BeZero())
				Expect(notes[1].Name).Should(Equal("note2"))
				Expect(len(notes[1].Labels)).Should(Equal(2))
				Expect(len(notes[1].Metadata)).Should(Equal(1))
			}

			By("authorization failure")
			{
				code := http.StatusUnauthorized
				client := utils.ClientWithJSON("{}", code)
				drvns, _ := drvnotestore.New(client)
				_, err := drvns.GetAll()
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(BeAssignableToTypeOf(errs.NewUnauthorizedError()))
			}

			By("inner error")
			{
				code := http.StatusInternalServerError
				client := utils.ClientWithJSON("error", code)
				drvns, _ := drvnotestore.New(client)
				notes, err := drvns.GetAll()
				Expect(err).Should(HaveOccurred())
				Expect(notes).Should(BeNil())
			}
		})
	})

	Context("get note by id", func() {
		It("should return the note successfully", func() {
			By("valid note id")
			{
				j := `{
					"id": "gdtg45w9mjh10ds",
					"name": "note",
					"properties": {
					  "labels": "label1,label2",
					  "meta!key1": "value1"
					},
					"createdTime": "2021-02-12T07:20:50.52Z"
				  }`
				client := utils.ClientWithJSON(j, http.StatusOK)
				drvns, _ := drvnotestore.New(client)
				note, err := drvns.Get("id")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(note).ShouldNot(BeNil())
				Expect(note.ID).ShouldNot(BeZero())
				Expect(note.Name).Should(Equal("note"))
				Expect(len(note.Labels)).Should(Equal(2))
				Expect(len(note.Metadata)).Should(Equal(1))
			}

			By("wrong note id")
			{
				client := utils.ClientWithJSON("{}", http.StatusNotFound)
				drvns, _ := drvnotestore.New(client)
				_, err := drvns.Get("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(BeAssignableToTypeOf(errs.NewNotFoundError("msg")))
			}

			By("authorization failure")
			{
				client := utils.ClientWithJSON("{}", http.StatusUnauthorized)
				drvns, _ := drvnotestore.New(client)
				_, err := drvns.Get("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(BeAssignableToTypeOf(errs.NewUnauthorizedError()))
			}

			By("inner error")
			{
				code := http.StatusInternalServerError
				client := utils.ClientWithJSON("error", code)
				drvns, _ := drvnotestore.New(client)
				note, err := drvns.Get("id")
				Expect(err).Should(HaveOccurred())
				Expect(note).Should(BeNil())
			}
		})
	})

	Context("update note", func() {
		It("should update the note successfully", func() {
			By("valid note id")
			{
				j := `{
					"id": "gdtg45w9mjh10ds",
					"name": "note",
					"description": "desc"
				}`
				client := utils.ClientWithJSON(j, http.StatusOK)
				drvns, _ := drvnotestore.New(client)
				note, err := drvns.Update("id", &notestore.WritableNote{
					Name:        "note",
					Description: "desc",
				})
				Expect(err).ShouldNot(HaveOccurred())
				Expect(note).ShouldNot(BeNil())
				Expect(note.Name).Should(Equal("note"))
			}

			By("wrong note id")
			{
				client := utils.ClientWithJSON("{}", http.StatusNotFound)
				drvns, _ := drvnotestore.New(client)
				_, err := drvns.Update("id", &notestore.WritableNote{
					Name: "note",
				})
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(BeAssignableToTypeOf(errs.NewNotFoundError("msg")))
			}

			By("authorization failure")
			{
				client := utils.ClientWithJSON("{}", http.StatusUnauthorized)
				drvns, _ := drvnotestore.New(client)
				_, err := drvns.Update("id", &notestore.WritableNote{
					Name: "note",
				})
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(BeAssignableToTypeOf(errs.NewUnauthorizedError()))
			}

			By("inner error")
			{
				code := http.StatusInternalServerError
				client := utils.ClientWithJSON("error", code)
				drvns, _ := drvnotestore.New(client)
				note, err := drvns.Update("id", &notestore.WritableNote{
					Name: "note",
				})
				Expect(err).Should(HaveOccurred())
				Expect(note).Should(BeNil())
			}
		})
	})

	Context("delete note", func() {
		It("should delete the note successfully", func() {
			By("valid note error")
			{
				client := utils.ClientWithJSON("{}", http.StatusOK)
				drvns, _ := drvnotestore.New(client)
				err := drvns.Delete("id")
				Expect(err).ShouldNot(HaveOccurred())
			}

			By("wrong note id")
			{
				client := utils.ClientWithJSON("{}", http.StatusNotFound)
				drvns, _ := drvnotestore.New(client)
				err := drvns.Delete("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(BeAssignableToTypeOf(errs.NewNotFoundError("msg")))
			}

			By("authorization failure")
			{
				client := utils.ClientWithJSON("{}", http.StatusUnauthorized)
				drvns, _ := drvnotestore.New(client)
				err := drvns.Delete("id")
				Expect(err).Should(HaveOccurred())
				Expect(err).Should(BeAssignableToTypeOf(errs.NewUnauthorizedError()))
			}

			By("inner error")
			{
				code := http.StatusInternalServerError
				client := utils.ClientWithJSON("error", code)
				drvns, _ := drvnotestore.New(client)
				err := drvns.Delete("id")
				Expect(err).Should(HaveOccurred())
			}
		})
	})
})

func count(m map[string]string, h func(k, v string) bool) int {
	count := 0
	for k, v := range m {
		if h(k, v) {
			count++
		}
	}
	return count
}
