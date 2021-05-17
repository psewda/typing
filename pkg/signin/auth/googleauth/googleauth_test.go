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
		It("should return valid url", func() {
			cred := strings.Replace(credTemplate, "{token_url}", tokenURL, 1)
			ga, _ := googleauth.New([]byte(cred))

			By("zero param")
			{
				u := ga.GetURL(utils.Empty, utils.Empty)
				Expect(u).ShouldNot(BeZero())
			}

			By("redirect and state")
			{
				redirect := "http://localhost:5050/redirect"
				u := ga.GetURL(redirect, "state")
				redirectQuery := fmt.Sprintf("&redirect_uri=%s", redirect)
				stateQuery := "&state=state"
				Expect(url.QueryUnescape(u)).Should(ContainSubstring(redirectQuery))
				Expect(u).Should(ContainSubstring(stateQuery))
			}
		})
	})

	Context("exchange authcode", func() {
		It("should return valid token", func() {
			var writeContent writeFunc
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writeContent(w)
			}))
			defer ts.Close()
			cred := strings.Replace(credTemplate, "{token_url}", ts.URL, 1)
			ga, _ := googleauth.New([]byte(cred))

			By("valid authcode")
			{
				accessToken := "90d6446as32safy868asa0d14870c"
				refreshToken := "302305ma3s837uags57ffvfs9jksq"
				expiresIn := 3600
				writeContent = func(w http.ResponseWriter) {
					w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
					format := "access_token=%s&refresh_token=%s&expires_in=%d"
					data := fmt.Sprintf(format, accessToken, refreshToken, expiresIn)
					w.Write([]byte(data))
				}

				token, err := ga.Exchange("valid--authcode")
				Expect(token).ShouldNot(BeNil())
				Expect(token.AccessToken).Should(Equal(accessToken))
				Expect(token.RefreshToken).Should(Equal(refreshToken))
				Expect(token.Expiry).Should(BeTemporally(">", time.Now()))
				Expect(err).ShouldNot(HaveOccurred())
			}

			By("wrong authcode")
			{
				writeContent = func(w http.ResponseWriter) {
					w.WriteHeader(http.StatusInternalServerError)
				}
				_, err := ga.Exchange("wrong-authcode")
				Expect(err).Should(HaveOccurred())
			}
		})
	})

	Context("refresh access token", func() {
		It("should return new access token", func() {
			var writeContent writeFunc
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writeContent(w)
			}))
			defer ts.Close()
			cred := strings.Replace(credTemplate, "{token_url}", ts.URL, 1)
			ga, _ := googleauth.New([]byte(cred))

			By("valid token")
			{
				writeContent = func(w http.ResponseWriter) {
					w.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
					j := `{
						"access_token": "access-token",
						"expires_in": 3600
					}`
					w.Write([]byte(j))
				}
				token, err := ga.Refresh("refresh-token")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(token.AccessToken).Should(Equal("access-token"))
				Expect(token.RefreshToken).Should(Equal("refresh-token"))
				Expect(token.Expiry).Should(BeTemporally(">", time.Now()))
			}

			By("non-ok status")
			{
				writeContent = func(w http.ResponseWriter) {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("{}"))
				}
				_, err := ga.Refresh("refresh-token")
				Expect(err).Should(HaveOccurred())
			}
		})
	})

	Context("revoke access token", func() {
		It("should revoke successfully w/o error", func() {
			var writeContent writeFunc
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				writeContent(w)
			}))
			defer ts.Close()
			cred := strings.Replace(credTemplate, "{token_url}", ts.URL, 1)
			ga, _ := googleauth.New([]byte(cred))

			By("valid token")
			{
				writeContent = func(w http.ResponseWriter) {
					w.WriteHeader(http.StatusOK)
				}
				err := ga.Revoke("access-token")
				Expect(err).ShouldNot(HaveOccurred())
			}

			By("wrong token")
			{
				writeContent = func(w http.ResponseWriter) {
					w.WriteHeader(http.StatusInternalServerError)
				}
				err := ga.Revoke("wrong-token")
				Expect(err).Should(HaveOccurred())
			}
		})
	})
})
