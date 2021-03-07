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
})
