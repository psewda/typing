package drvsectionstore_test

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"mime"
	"mime/multipart"
	"net/http"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/errs"
	"github.com/psewda/typing/pkg/storage/sectionstore"
	"github.com/psewda/typing/pkg/storage/sectionstore/drvsectionstore"
)

func TestDrvSectionstore(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "drvsectionstore-suite")
}

var _ = Describe("googledrive sectionstore", func() {
	Context("create new section", func() {
		It("should succeed when correct input - w/ existing content", func() {
			verifyReq := func(req *http.Request) {
				var sections []*sectionstore.Section
				content := readContent(req)
				json.Unmarshal(content, &sections)

				Expect(sections).Should(HaveLen(2))
				Expect(sections[1].ID).ShouldNot(BeEmpty())
				Expect(sections[1].Name).Should(Equal("new-section"))
				Expect(sections[1].Labels).Should(HaveLen(2))
				Expect(sections[1].Labels).Should(ContainElements("label1", "label2"))
				Expect(sections[1].Metadata).Should(HaveLen(1))
				Expect(sections[1].Metadata).Should(HaveKeyWithValue("meta1", "value1"))
				Expect(sections[1].Data).Should(HaveLen(1))
				Expect(sections[1].Data).Should(HaveKeyWithValue("item1", "value1"))
			}

			client := http.DefaultClient
			client.Transport = utils.TransportFunc(func(req *http.Request) (*http.Response, error) {
				j := `[
						{
							"id": "gdtg45w9mjh10ds",
							"name": "section",
							"data": {
								"item1": "value1"
							}
						}
					]`
				if req.Method == "PATCH" {
					verifyReq(req)
					j = `{ "id": "nid" }`
				}
				return buildResponse(http.StatusOK, j), nil
			})

			dss, _ := drvsectionstore.New(client)
			section, err := dss.Create("nid", &sectionstore.WritableSection{
				Name:   "new-section",
				Labels: []string{"label1", "label2"},
				Metadata: map[string]string{
					"meta1": "value1",
				},
				Data: map[string]string{
					"item1": "value1",
				},
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(section).ShouldNot(BeNil())
			Expect(section.ID).ShouldNot(BeEmpty())
			Expect(section.Name).Should(Equal("new-section"))
		})

		It("should succeed when correct input - w/o existing content", func() {
			verifyReq := func(req *http.Request) {
				var sections []*sectionstore.Section
				content := readContent(req)
				json.Unmarshal(content, &sections)

				Expect(sections).Should(HaveLen(1))
				Expect(sections[0].ID).ShouldNot(BeEmpty())
				Expect(sections[0].Name).Should(Equal("new-section"))
				Expect(sections[0].Data).Should(HaveLen(1))
				Expect(sections[0].Data).Should(HaveKeyWithValue("item1", "value1"))
			}

			client := http.DefaultClient
			client.Transport = utils.TransportFunc(func(req *http.Request) (*http.Response, error) {
				j := ``
				if req.Method == "PATCH" {
					verifyReq(req)
					j = `{ "id": "nid" }`
				}
				return buildResponse(http.StatusOK, j), nil
			})

			dss, _ := drvsectionstore.New(client)
			section, err := dss.Create("nid", &sectionstore.WritableSection{
				Name: "new-section",
				Data: map[string]string{
					"item1": "value1",
				},
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(section).ShouldNot(BeNil())
			Expect(section.ID).ShouldNot(BeEmpty())
			Expect(section.Name).Should(Equal("new-section"))
		})

		It("should succeed when unsanitized input", func() {
			verifyReq := func(req *http.Request) {
				var sections []*sectionstore.Section
				content := readContent(req)
				json.Unmarshal(content, &sections)

				Expect(sections).Should(HaveLen(1))
				Expect(sections[0].ID).ShouldNot(BeEmpty())
				Expect(sections[0].Name).Should(Equal("new-section"))
				Expect(sections[0].Labels).Should(HaveLen(2))
				Expect(sections[0].Labels).Should(ContainElements("label1", "label2"))
				Expect(sections[0].Metadata).Should(HaveLen(1))
				Expect(sections[0].Metadata).Should(HaveKeyWithValue("meta1", "value1"))
				Expect(sections[0].Data).Should(HaveLen(1))
				Expect(sections[0].Data).Should(HaveKeyWithValue("item1", "value1"))
			}

			client := http.DefaultClient
			client.Transport = utils.TransportFunc(func(req *http.Request) (*http.Response, error) {
				j := ``
				if req.Method == "PATCH" {
					verifyReq(req)
					j = `{ "id": "nid" }`
				}
				return buildResponse(http.StatusOK, j), nil
			})

			dss, _ := drvsectionstore.New(client)
			section, err := dss.Create("nid", &sectionstore.WritableSection{
				Name:   " new-section ",
				Labels: []string{"label1", " ", "label2  "},
				Metadata: map[string]string{
					"meta1  ": "value1   ",
					" ":       "  value2",
				},
				Data: map[string]string{
					"item1": "value1  ",
				},
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(section).ShouldNot(BeNil())
			Expect(section.ID).ShouldNot(BeEmpty())
			Expect(section.Name).Should(Equal("new-section"))
		})

		It("should return error when wrong input", func() {
			dss, _ := drvsectionstore.New(http.DefaultClient)
			section, err := dss.Create("nid", &sectionstore.WritableSection{
				Labels: []string{"label1", "label2"},
			})

			Expect(err).Should(HaveOccurred())
			Expect(section).Should(BeNil())
		})

		It("should return error when wrong note id", func() {
			code := http.StatusNotFound
			client := utils.ClientWithJSON("{}", code)
			dss, _ := drvsectionstore.New(client)
			_, err := dss.Create("nid", &sectionstore.WritableSection{
				Name: "section",
			})

			Expect(err).Should(HaveOccurred())
			Expect(err).Should(BeAssignableToTypeOf(errs.NewNotFoundError("msg")))
		})

		It("should return error when authorization failure", func() {
			code := http.StatusUnauthorized
			client := utils.ClientWithJSON("{}", code)
			dss, _ := drvsectionstore.New(client)
			_, err := dss.Create("nid", &sectionstore.WritableSection{
				Name: "section",
			})

			Expect(err).Should(HaveOccurred())
			Expect(err).Should(BeAssignableToTypeOf(errs.NewUnauthorizedError()))
		})
	})

	Context("get all sections", func() {
		It("should return all sections when correct setup", func() {
			j := `[
						{
							"id": "gdtg45gdre6d3s",
							"name": "section1",
							"data": {
								"item1": "value1"
							}
						},
						{
							"id": "nfgd4w6sj7sh9s",
							"name": "section2",
							"data": {
								"item1": "value1"
							}
						}
					]`
			client := utils.ClientWithJSON(j, http.StatusOK)
			dss, _ := drvsectionstore.New(client)
			sections, err := dss.GetAll("nid")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(sections).ShouldNot(BeNil())
			Expect(sections).Should(HaveLen(2))
			Expect(sections[0].ID).ShouldNot(BeEmpty())
			Expect(sections[0].Name).Should(Equal("section1"))
			Expect(sections[0].Data).Should(HaveLen(1))
			Expect(sections[0].Data).Should(HaveKeyWithValue("item1", "value1"))
			Expect(sections[1].ID).ShouldNot(BeEmpty())
			Expect(sections[1].Name).Should(Equal("section2"))
			Expect(sections[1].Data).Should(HaveLen(1))
			Expect(sections[1].Data).Should(HaveKeyWithValue("item1", "value1"))
		})

		It("should return nil when no note content", func() {
			client := utils.ClientWithJSON(``, http.StatusOK)
			dss, _ := drvsectionstore.New(client)
			sections, err := dss.GetAll("nid")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(sections).Should(BeNil())
		})

		It("should return error when download failure", func() {
			code := getRndHTTPErrorCode()
			client := utils.ClientWithJSON("{}", code)
			dss, _ := drvsectionstore.New(client)
			_, err := dss.GetAll("nid")

			assertDownloadError(err, code)
		})
	})

	Context("get section by id", func() {
		It("should return the section when valid section id", func() {
			j := `[
					{
						"id": "secid1",
						"name": "section1",
						"data": {
							"item1": "value1"
						}
					},
					{
						"id": "secid2",
						"name": "section2",
						"data": {
							"item1": "value1"
						}
					}
				]`
			client := utils.ClientWithJSON(j, http.StatusOK)
			dss, _ := drvsectionstore.New(client)
			section, err := dss.Get("nid", "secid1")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(section).ShouldNot(BeNil())

			Expect(section.ID).Should(Equal("secid1"))
			Expect(section.Name).Should(Equal("section1"))
			Expect(section.Data).Should(HaveLen(1))
			Expect(section.Data).Should(HaveKeyWithValue("item1", "value1"))
		})

		It("should return error when wrong section id", func() {
			j := `[
					{
						"id": "secid1",
						"name": "section1",
						"data": {
							"item1": "value1"
						}
					}
				]`
			client := utils.ClientWithJSON(j, http.StatusOK)
			dss, _ := drvsectionstore.New(client)
			_, err := dss.Get("nid", "wrong")

			Expect(err).Should(HaveOccurred())
			Expect(err).Should(BeAssignableToTypeOf(errs.NewNotFoundError("msg")))
		})

		It("should return error when download failure", func() {
			code := getRndHTTPErrorCode()
			client := utils.ClientWithJSON("{}", code)
			dss, _ := drvsectionstore.New(client)
			_, err := dss.Get("nid", "secid")

			assertDownloadError(err, code)
		})
	})

	Context("update section", func() {
		It("should succeed when correct input", func() {
			verifyReq := func(req *http.Request) {
				var sections []*sectionstore.Section
				content := readContent(req)
				json.Unmarshal(content, &sections)

				Expect(sections).Should(HaveLen(1))
				Expect(sections[0].ID).ShouldNot(BeEmpty())
				Expect(sections[0].Name).Should(Equal("section-updated"))
				Expect(sections[0].Labels).Should(HaveLen(2))
				Expect(sections[0].Labels).Should(ContainElements("label1", "label2"))
				Expect(sections[0].Metadata).Should(HaveLen(1))
				Expect(sections[0].Metadata).Should(HaveKeyWithValue("meta1", "value1"))
				Expect(sections[0].Data).Should(HaveLen(1))
				Expect(sections[0].Data).Should(HaveKeyWithValue("item1", "value1-updated"))
			}

			client := http.DefaultClient
			client.Transport = utils.TransportFunc(func(req *http.Request) (*http.Response, error) {
				j := `[
						{
							"id": "secid",
							"name": "section",
							"data": {
								"item1": "value1"
							}
						}
					]`
				if req.Method == "PATCH" {
					verifyReq(req)
					j = `{ "id": "nid" }`
				}
				return buildResponse(http.StatusOK, j), nil
			})

			dss, _ := drvsectionstore.New(client)
			section, err := dss.Update("nid", "secid", &sectionstore.WritableSection{
				Name:   "section-updated",
				Labels: []string{"label1", "label2"},
				Metadata: map[string]string{
					"meta1": "value1",
				},
				Data: map[string]string{
					"item1": "value1-updated",
				},
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(section).ShouldNot(BeNil())
			Expect(section.ID).ShouldNot(BeEmpty())
			Expect(section.Name).Should(Equal("section-updated"))
		})

		It("should succeed when unsanitized input", func() {
			verifyReq := func(req *http.Request) {
				var sections []*sectionstore.Section
				content := readContent(req)
				json.Unmarshal(content, &sections)

				Expect(sections).Should(HaveLen(1))
				Expect(sections[0].ID).ShouldNot(BeEmpty())
				Expect(sections[0].Name).Should(Equal("section-updated"))
				Expect(sections[0].Labels).Should(HaveLen(2))
				Expect(sections[0].Labels).Should(ContainElements("label1", "label2"))
				Expect(sections[0].Metadata).Should(HaveLen(1))
				Expect(sections[0].Metadata).Should(HaveKeyWithValue("meta1", "value1"))
				Expect(sections[0].Data).Should(HaveLen(1))
				Expect(sections[0].Data).Should(HaveKeyWithValue("item1", "value1-updated"))
			}

			client := http.DefaultClient
			client.Transport = utils.TransportFunc(func(req *http.Request) (*http.Response, error) {
				j := `[
						{
							"id": "secid",
							"name": "section",
							"data": {
								"item1": "value1"
							}
						}
					]`
				if req.Method == "PATCH" {
					verifyReq(req)
					j = `{ "id": "nid" }`
				}
				return buildResponse(http.StatusOK, j), nil
			})

			dss, _ := drvsectionstore.New(client)
			section, err := dss.Update("nid", "secid", &sectionstore.WritableSection{
				Name:   " section-updated ",
				Labels: []string{"label1", " ", "label2  "},
				Metadata: map[string]string{
					"meta1  ": "value1   ",
					" ":       "  value2",
				},
				Data: map[string]string{
					"item1": "value1-updated  ",
				},
			})

			Expect(err).ShouldNot(HaveOccurred())
			Expect(section).ShouldNot(BeNil())
			Expect(section.ID).ShouldNot(BeEmpty())
			Expect(section.Name).Should(Equal("section-updated"))
		})

		It("should return error when wrong input", func() {
			dss, _ := drvsectionstore.New(http.DefaultClient)
			section, err := dss.Update("nid", "sid", &sectionstore.WritableSection{
				Labels: []string{"label1", "label2"},
			})

			Expect(err).Should(HaveOccurred())
			Expect(section).Should(BeNil())
		})

		It("should return error when wrong section id", func() {
			j := `[
					{
						"id": "secid1",
						"name": "section1",
						"data": {
							"item1": "value1"
						}
					}
				]`
			client := utils.ClientWithJSON(j, http.StatusOK)
			dss, _ := drvsectionstore.New(client)
			_, err := dss.Update("nid", "wrong", &sectionstore.WritableSection{
				Name: "section-updated",
				Data: map[string]string{
					"item1": "value1-updated",
				},
			})

			Expect(err).Should(HaveOccurred())
			Expect(err).Should(BeAssignableToTypeOf(errs.NewNotFoundError("msg")))
		})

		It("should return error when download failure", func() {
			code := getRndHTTPErrorCode()
			client := utils.ClientWithJSON("{}", code)
			dss, _ := drvsectionstore.New(client)
			_, err := dss.Update("nid", "sid", &sectionstore.WritableSection{
				Name: "section-updated",
			})

			assertDownloadError(err, code)
		})
	})

	Context("delete section", func() {
		It("should succeed when correct section id", func() {
			verifyReq := func(req *http.Request) {
				var sections []*sectionstore.Section
				content := readContent(req)
				json.Unmarshal(content, &sections)

				Expect(sections).Should(HaveLen(1))
				Expect(sections[0].ID).ShouldNot(Equal("secid1"))
				Expect(sections[0].Name).ShouldNot(Equal("section1"))
			}

			client := http.DefaultClient
			client.Transport = utils.TransportFunc(func(req *http.Request) (*http.Response, error) {
				j := `[
						{
							"id": "secid1",
							"name": "section1",
							"data": {
								"item1": "value1"
							}
						},
						{
							"id": "secid2",
							"name": "section2",
							"data": {
								"item1": "value1"
							}
						}
					]`
				if req.Method == "PATCH" {
					verifyReq(req)
					j = `{ "id": "nid" }`
				}
				return buildResponse(http.StatusOK, j), nil
			})

			dss, _ := drvsectionstore.New(client)
			err := dss.Delete("nid", "secid1")

			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return error when wrong section id", func() {
			j := `[
					{
						"id": "secid1",
						"name": "section1",
						"data": {
							"item1": "value1"
						}
					}
				]`
			client := utils.ClientWithJSON(j, http.StatusOK)
			dss, _ := drvsectionstore.New(client)
			err := dss.Delete("nid", "wrong")

			Expect(err).Should(HaveOccurred())
			Expect(err).Should(BeAssignableToTypeOf(errs.NewNotFoundError("msg")))
		})

		It("should return error when download failure", func() {
			code := getRndHTTPErrorCode()
			client := utils.ClientWithJSON("{}", code)
			dss, _ := drvsectionstore.New(client)
			err := dss.Delete("nid", "sid")

			assertDownloadError(err, code)
		})
	})
})

func readContent(req *http.Request) []byte {
	_, params, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
	reader := multipart.NewReader(req.Body, params["boundary"])

	// move to second part to read uploaded data.
	reader.NextPart()
	p, _ := reader.NextPart()

	content, _ := ioutil.ReadAll(p)
	return content
}

func buildResponse(code int, j string) *http.Response {
	return &http.Response{
		StatusCode: code,
		Body:       ioutil.NopCloser(bytes.NewBufferString(j)),
		Header:     map[string][]string{"Content-Type": {"application/json"}},
	}
}

func getRndHTTPErrorCode() int {
	httpErrorCodes := []int{
		http.StatusUnauthorized,
		http.StatusNotFound,
		http.StatusInternalServerError,
	}

	src := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(src)
	idx := rnd.Intn(len(httpErrorCodes) - 1)
	return httpErrorCodes[idx]
}

func assertDownloadError(err error, code int) {
	switch code {
	case http.StatusNotFound:
		Expect(err).Should(BeAssignableToTypeOf(errs.NewNotFoundError("msg")))
	case http.StatusUnauthorized:
		Expect(err).Should(BeAssignableToTypeOf(errs.NewUnauthorizedError()))

	default:
		Expect(err).Should(HaveOccurred())
	}
}
