package utils_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/errs"
	"google.golang.org/api/googleapi"
)

func TestUtils(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "utils-suite")
}

var _ = Describe("utility functions", func() {
	Context("utility function: GetValueString()", func() {
		It("should cover all cases", func() {
			By("valid value")
			v := utils.GetValueString("value", utils.Empty)
			Expect(v).Should(Equal("value"))

			By("empty value")
			v = utils.GetValueString(utils.Empty, "default")
			Expect(v).Should(Equal("default"))
		})
	})

	Context("utility function: Error()", func() {
		It("should cover all cases", func() {
			By("valid input")
			err := utils.Error("value", errors.New("inner"))
			Expect(err.Error()).Should(Equal("value: [inner]"))

			By("zero input")
			err = utils.Error(utils.Empty, nil)
			Expect(err.Error()).Should(Equal("error"))
		})
	})

	Context("utility function: AppendError()", func() {
		It("should cover all cases", func() {
			By("valid input")
			msg := utils.AppendError("value", errors.New("inner"))
			Expect(msg).Should(Equal("value: [inner]"))

			By("only value")
			msg = utils.AppendError("value", nil)
			Expect(msg).Should(Equal("value"))
		})
	})

	Context("utility function: CheckLocalhostURL()", func() {
		It("should cover all cases", func() {
			By("valid url")
			err := utils.CheckLocalhostURL("http://localhost:5050/redirect")
			Expect(err).ShouldNot(HaveOccurred())

			By("wrong url")
			err = utils.CheckLocalhostURL("invalid-value")
			Expect(err).Should(HaveOccurred())

			By("non localhost url")
			err = utils.CheckLocalhostURL("http://example.com:5050/redirect")
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("utility function: ClientWithToken()", func() {
		It("should cover all cases", func() {
			const accessToken = "access-token"
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				headerValue := r.Header.Get(echo.HeaderAuthorization)
				w.Write([]byte(headerValue))
			}))
			defer ts.Close()

			client := utils.ClientWithToken(accessToken)
			Expect(client).ShouldNot(BeNil())

			res, err := client.Get(ts.URL)
			Expect(err).ShouldNot(HaveOccurred())
			defer res.Body.Close()

			at, _ := ioutil.ReadAll(res.Body)
			t := fmt.Sprintf("Bearer %s", accessToken)
			Expect(string(at)).Should(Equal(t))
		})
	})

	Context("utility function: ClientWithJSON", func() {
		It("should cover all cases", func() {
			j := `{"key": "value"}`
			client := utils.ClientWithJSON(j, http.StatusCreated)
			Expect(client).ShouldNot(BeNil())

			res, err := client.Get("url")
			Expect(err).ShouldNot(HaveOccurred())
			defer res.Body.Close()

			json, _ := ioutil.ReadAll(res.Body)
			Expect(string(json)).Should(Equal(j))
			Expect(res.StatusCode).Should(Equal(http.StatusCreated))
		})
	})

	Context("utility function: ValidateStruct", func() {
		It("should cover all cases", func() {
			type Value struct {
				Name        string `json:"name,omitempty" validate:"required,notblank"`
				Description string `json:"desc,omitempty" validate:"max=250"`
			}

			m := make(map[string]string)
			m["name.required"] = "name is required field"
			m["name.notblank"] = "name can't be empty value"
			m["description.max"] = "desc must be less than 250 chars"

			By("validation successful")
			{
				v := &Value{
					Name:        "name",
					Description: "desc",
				}
				err := utils.ValidateStruct(v, m)
				Expect(err).ShouldNot(HaveOccurred())
			}

			By("validation error")
			{
				v := &Value{
					Name:        utils.Empty,
					Description: "desc",
				}
				err := utils.ValidateStruct(v, m)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("name is required field"))
			}
		})
	})

	Context("utility function: Sanitize", func() {
		It("should cover all cases", func() {
			m := map[string]string{
				"k1":   "v1",
				" k2 ": "v2",
				" ":    "v3",
			}
			sanitized := utils.Sanitize(m)

			Expect(sanitized).Should(HaveLen(2))
			Expect(sanitized).Should(HaveKeyWithValue("k1", "v1"))
			Expect(sanitized).Should(HaveKeyWithValue("k2", "v2"))
		})
	})

	Context("utility function: GetStatusCode", func() {
		It("should cover all cases", func() {
			By("status code 401")
			{
				err := googleapi.Error{
					Code:    http.StatusUnauthorized,
					Message: "unauthorized",
				}
				code := utils.GetStatusCode(&err)
				Expect(code).Should(Equal(http.StatusUnauthorized))
			}

			By("status code -1")
			{
				code := utils.GetStatusCode(errors.New("error"))
				Expect(code).Should(Equal(-1))
			}
		})
	})

	Context("utility function: BuildHTTPError", func() {
		It("should cover all cases", func() {
			By("unauthorized")
			{
				err := errs.NewUnauthorizedError()
				httpError := utils.BuildHTTPError(err, utils.Empty)
				Expect(httpError.Code).Should(Equal(http.StatusUnauthorized))
			}

			By("not found")
			{
				err := errs.NewNotFoundError("msg")
				httpError := utils.BuildHTTPError(err, utils.Empty)
				Expect(httpError.Code).Should(Equal(http.StatusNotFound))
			}

			By("internal server error")
			{
				err := errors.New("msg")
				httpError := utils.BuildHTTPError(err, utils.Empty)
				Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
			}
		})
	})
})
