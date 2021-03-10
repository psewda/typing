package v1_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/psewda/typing/internal/utils"
	"github.com/psewda/typing/mocks"
	ctrlv1 "github.com/psewda/typing/pkg/controllers/v1"
	"github.com/psewda/typing/pkg/signin/auth"
	"github.com/psewda/typing/pkg/types"
)

var _ = Describe("auth controller", func() {
	var (
		mockContainer *mocks.MockContainer
		mockAuth      *mocks.MockAuth
		rec           *httptest.ResponseRecorder
	)

	BeforeEach(func() {
		mockContainer = mocks.NewMockContainer(mockCtrl)
		mockAuth = mocks.NewMockAuth(mockCtrl)
		rec = httptest.NewRecorder()
	})

	Context("get auth url", func() {
		It("should return valid url or appropriate error", func() {
			const returnURL = "https://accounts.google.com/o/oauth2/auth"

			By("no param")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				mockAuth.EXPECT().GetURL(gomock.Any(), gomock.Any()).Return(returnURL, nil)
				req := httptest.NewRequest(http.MethodGet, urlRoute, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer))

				ctrlv1.GetURL(ctx)
				Expect(rec.Code).Should(Equal(http.StatusOK))

				var urlValue types.URLValue
				json.NewDecoder(rec.Body).Decode(&urlValue)
				Expect(urlValue).ShouldNot(BeZero())
				Expect(urlValue.URL).Should(Equal(returnURL))
			}

			By("valid redirect url")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				mockAuth.EXPECT().GetURL(gomock.Any(), gomock.Any()).Return(returnURL, nil)
				u := fmt.Sprintf("%s?redirect=http://localhost:7070/redirect", urlRoute)
				req := httptest.NewRequest(http.MethodGet, u, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer))

				ctrlv1.GetURL(ctx)
				Expect(rec.Code).Should(Equal(http.StatusOK))

				var urlValue types.URLValue
				json.NewDecoder(rec.Body).Decode(&urlValue)
				Expect(urlValue).ShouldNot(BeZero())
				Expect(urlValue.URL).Should(Equal(returnURL))
			}

			By("invalid redirect url")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				u := fmt.Sprintf("%s?redirect=invalid-url", urlRoute)
				req := httptest.NewRequest(http.MethodGet, u, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer))

				err := ctrlv1.GetURL(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusBadRequest))
			}

			By("inner error")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				mockAuth.EXPECT().GetURL(gomock.Any(), gomock.Any()).Return(utils.Empty, errors.New("error"))
				req := httptest.NewRequest(http.MethodGet, urlRoute, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer))

				err := ctrlv1.GetURL(ctx)
				httpError := err.(*echo.HTTPError)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
			}
		})
	})

	Context("token exchange", func() {
		It("should return valid token or appropriate error", func() {

			By("valid authcode")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				token := &auth.Token{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					Expiry:       time.Now().Add(time.Second * 30),
				}
				mockAuth.EXPECT().Exchange(gomock.Any()).Return(token, nil)
				req := httptest.NewRequest(http.MethodPost, tokenRoute, nil)
				form := url.Values{}
				form.Add("auth_code", "valid-authcode")
				req.PostForm = form
				ctx := newCtx(req, rec, withContainer(mockContainer))

				ctrlv1.Exchange(ctx)
				Expect(rec.Code).Should(Equal(http.StatusOK))

				var t auth.Token
				json.NewDecoder(rec.Body).Decode(&t)
				Expect(t).ShouldNot(BeNil())
				Expect(t.AccessToken).Should(Equal("access-token"))
				Expect(t.RefreshToken).Should(Equal("refresh-token"))
				Expect(t.Expiry).Should(BeTemporally(">", time.Now()))
			}

			By("empty authcode")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				req := httptest.NewRequest(http.MethodPost, tokenRoute, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer))

				err := ctrlv1.Exchange(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusBadRequest))
			}

			By("inner error")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				mockAuth.EXPECT().Exchange(gomock.Any()).Return(nil, errors.New("error"))
				req := httptest.NewRequest(http.MethodPost, tokenRoute, nil)
				form := url.Values{}
				form.Add("auth_code", "valid-authcode")
				req.PostForm = form
				ctx := newCtx(req, rec, withContainer(mockContainer))

				err := ctrlv1.Exchange(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
			}
		})
	})

	Context("token refresh", func() {
		It("should return valid token or appropriate error", func() {

			By("valid token")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				token := &auth.Token{
					AccessToken:  "access-token",
					RefreshToken: "refresh-token",
					Expiry:       time.Now().Add(time.Second * 30),
				}
				mockAuth.EXPECT().Refresh(gomock.Any()).Return(token, nil)
				req := httptest.NewRequest(http.MethodPost, refreshRoute, nil)
				form := url.Values{}
				form.Add("refresh_token", "valid-token")
				req.PostForm = form
				ctx := newCtx(req, rec, withContainer(mockContainer))

				ctrlv1.Refresh(ctx)
				Expect(rec.Code).Should(Equal(http.StatusOK))

				var t auth.Token
				json.NewDecoder(rec.Body).Decode(&t)
				Expect(t).ShouldNot(BeNil())
				Expect(t.AccessToken).Should(Equal("access-token"))
				Expect(t.RefreshToken).Should(Equal("refresh-token"))
				Expect(t.Expiry).Should(BeTemporally(">", time.Now()))
			}

			By("empty token")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				req := httptest.NewRequest(http.MethodPost, refreshRoute, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer))

				err := ctrlv1.Refresh(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusBadRequest))
			}

			By("inner error")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				mockAuth.EXPECT().Refresh(gomock.Any()).Return(nil, errors.New("error"))
				req := httptest.NewRequest(http.MethodPost, tokenRoute, nil)
				form := url.Values{}
				form.Add("refresh_token", "valid-token")
				req.PostForm = form
				ctx := newCtx(req, rec, withContainer(mockContainer))

				err := ctrlv1.Refresh(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
			}
		})
	})

	Context("revoke token", func() {
		It("should return 'NoContent' or appropriate error", func() {

			By("valid token")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				mockAuth.EXPECT().Revoke(gomock.Any()).Return(nil)
				req := httptest.NewRequest(http.MethodPost, revokeRoute, nil)
				form := url.Values{}
				form.Add("access_token", "valid-token")
				req.PostForm = form
				ctx := newCtx(req, rec, withContainer(mockContainer))

				ctrlv1.Revoke(ctx)
				Expect(rec.Code).Should(Equal(http.StatusNoContent))
			}

			By("empty token")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				req := httptest.NewRequest(http.MethodPost, revokeRoute, nil)
				ctx := newCtx(req, rec, withContainer(mockContainer))

				err := ctrlv1.Revoke(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusBadRequest))
			}

			By("inner error")
			{
				mockContainer.EXPECT().GetInstance(gomock.Any()).Return(mockAuth, nil)
				mockAuth.EXPECT().Revoke(gomock.Any()).Return(errors.New("error"))
				req := httptest.NewRequest(http.MethodGet, revokeRoute, nil)
				form := url.Values{}
				form.Add("access_token", "valid-token")
				req.PostForm = form
				ctx := newCtx(req, rec, withContainer(mockContainer))

				err := ctrlv1.Revoke(ctx)
				httpError := toHTTPError(err)
				Expect(httpError).Should(HaveOccurred())
				Expect(httpError.Code).Should(Equal(http.StatusInternalServerError))
			}
		})
	})
})