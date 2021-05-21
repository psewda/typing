package googleauth_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/pkg/signin/auth/googleauth"
)

func TestGoogleAuth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "googleauth-suite")
}

var _ = Describe("googleauth", func() {
	type writeFunc func(w http.ResponseWriter)
	const tokenURL = "https://oauth2.googleapis.com/token"
	const credTemplate = `{
		"installed": {
		  "client_id": "1122-8sahye34j.apps.googleusercontent.com",
		  "project_id": "typing-888",
		  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
		  "token_uri": "{token_url}",
		  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		  "client_secret": "vdfal523nn233",
		  "redirect_uris": [
			"urn:ietf:wg:oauth:2.0:oob",
			"http://localhost"
		  ]
		}
	  }`

	Context("get auth url", func() {
		It("should return auth url when having zero param", func() {
			cred := strings.Replace(credTemplate, "{token_url}", tokenURL, 1)
			ga, _ := googleauth.New([]byte(cred))
			u := ga.GetURL(utils.Empty, utils.Empty)
			Expect(u).ShouldNot(BeZero())
		})

		It("should return auth url when having redirect and state values", func() {
			cred := strings.Replace(credTemplate, "{token_url}", tokenURL, 1)
			ga, _ := googleauth.New([]byte(cred))
			redirect := "http://localhost:5050/redirect"
			u := ga.GetURL(redirect, "state")
			redirectQuery := fmt.Sprintf("&redirect_uri=%s", redirect)
			stateQuery := "&state=state"
			Expect(url.QueryUnescape(u)).Should(ContainSubstring(redirectQuery))
			Expect(u).Should(ContainSubstring(stateQuery))
		})
	})

	Context("exchange authcode", func() {
		var writeContent writeFunc
		var ts *httptest.Server

		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writeContent(w)
			}))
		})

		AfterEach(func() {
			ts.Close()
		})

		It("should succeed when valid authcode", func() {
			accessToken := "90d6446as32safy868asa0d14870c"
			refreshToken := "302305ma3s837uags57ffvfs9jksq"
			writeContent = func(w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
				format := "access_token=%s&refresh_token=%s&expires_in=%d"
				data := fmt.Sprintf(format, accessToken, refreshToken, 3600)
				w.Write([]byte(data))
			}

			cred := strings.Replace(credTemplate, "{token_url}", ts.URL, 1)
			ga, _ := googleauth.New([]byte(cred))
			token, err := ga.Exchange("valid-authcode")

			Expect(token).ShouldNot(BeNil())
			Expect(token.AccessToken).Should(Equal(accessToken))
			Expect(token.RefreshToken).Should(Equal(refreshToken))
			Expect(token.Expiry).Should(BeTemporally(">", time.Now()))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return error when wrong authcode", func() {
			writeContent = func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusInternalServerError)
			}

			cred := strings.Replace(credTemplate, "{token_url}", ts.URL, 1)
			ga, _ := googleauth.New([]byte(cred))
			_, err := ga.Exchange("wrong-authcode")
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("refresh token", func() {
		var writeContent writeFunc
		var ts *httptest.Server

		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writeContent(w)
			}))
		})

		AfterEach(func() {
			ts.Close()
		})

		It("should succeed when valid refresh token", func() {
			writeContent = func(w http.ResponseWriter) {
				w.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
				j := `{
						"access_token": "access-token",
						"expires_in": 3600
					}`
				w.Write([]byte(j))
			}
			cred := strings.Replace(credTemplate, "{token_url}", ts.URL, 1)
			ga, _ := googleauth.New([]byte(cred))
			token, err := ga.Refresh("refresh-token")

			Expect(err).ShouldNot(HaveOccurred())
			Expect(token.AccessToken).Should(Equal("access-token"))
			Expect(token.RefreshToken).Should(Equal("refresh-token"))
			Expect(token.Expiry).Should(BeTemporally(">", time.Now()))
		})

		It("should return error when http status code != OK", func() {
			writeContent = func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("{}"))
			}
			cred := strings.Replace(credTemplate, "{token_url}", ts.URL, 1)
			ga, _ := googleauth.New([]byte(cred))
			_, err := ga.Refresh("refresh-token")
			Expect(err).Should(HaveOccurred())
		})
	})

	Context("revoke token", func() {
		var writeContent writeFunc
		var ts *httptest.Server

		BeforeEach(func() {
			ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writeContent(w)
			}))
		})

		AfterEach(func() {
			ts.Close()
		})

		It("should succeed when valid token", func() {
			writeContent = func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
			}
			cred := strings.Replace(credTemplate, "{token_url}", ts.URL, 1)
			ga, _ := googleauth.New([]byte(cred))
			err := ga.Revoke("valid-token")
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return error when wrong token", func() {
			writeContent = func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusInternalServerError)
			}
			cred := strings.Replace(credTemplate, "{token_url}", ts.URL, 1)
			ga, _ := googleauth.New([]byte(cred))
			err := ga.Revoke("wrong-token")
			Expect(err).Should(HaveOccurred())
		})
	})
})
